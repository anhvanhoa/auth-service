package usecase

import (
	"auth-service/constants"
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/anhvanhoa/service-core/common"
	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/oops"
	"github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/domain/saga"
	"github.com/anhvanhoa/service-core/domain/token"
)

type ForgotPasswordType string

const (
	ForgotByCode  ForgotPasswordType = "ForgotByCode"
	ForgotByToken ForgotPasswordType = "ForgotByToken"
)

var (
	ErrValidateForgotPassword = oops.New("Phương thức xác thực không hợp lệ, vui lòng chọn code hoặc token")
	ErrCreateSession          = oops.New("Không thể tạo phiên làm việc")
	ErrInvalidCompensateType  = oops.New("Loại bù không hợp lệ")
)

type ForgotPasswordRes struct {
	User  entity.UserInfor
	Token string
	Code  string
}

type CompensateForgotPassword struct {
	UserID string
	Token  string
	Code   string
	Type   ForgotPasswordType
}

type ForgotPasswordUsecase interface {
	ForgotPassword(email, os string, method ForgotPasswordType) (ForgotPasswordRes, error)
	saveCodeOrToken(typeForgot ForgotPasswordType, userID, codeOrToken, os string, exp time.Time) error
	SendEmailForgotPassword(payload queue.PayloadI) (string, error)
	CompensateSendEmail(ctx context.Context, taskID string) error
	generateRandomCode(length int) string
	ForgotPasswordWithSaga(sagaID string, execute common.ExecuteSaga) error
	CompensateForgotPassword(ctx context.Context, data CompensateForgotPassword) error
}

type forgotPasswordUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tx          repository.ManagerTransaction
	token       token.TokenForgotPasswordI
	cache       cache.CacheI
	qc          queue.QueueClient
	saga        saga.SagaManager
	log         *log.LogGRPCImpl
}

func NewForgotPasswordUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	tx repository.ManagerTransaction,
	token token.TokenForgotPasswordI,
	cache cache.CacheI,
	qc queue.QueueClient,
	saga saga.SagaManager,
	log *log.LogGRPCImpl,
) ForgotPasswordUsecase {
	return &forgotPasswordUsecaseImpl{
		userRepo,
		sessionRepo,
		tx,
		token,
		cache,
		qc,
		saga,
		log,
	}
}

func (uc *forgotPasswordUsecaseImpl) saveCodeOrToken(typeForgot ForgotPasswordType, userID, codeOrToken, os string, exp time.Time) error {
	session := entity.Session{
		Token:     codeOrToken,
		UserID:    userID,
		Type:      entity.SessionTypeForgot,
		Os:        os,
		ExpiredAt: exp,
		CreatedAt: time.Now(),
	}
	key := codeOrToken
	if typeForgot == ForgotByCode && len(codeOrToken) == 6 {
		key = fmt.Sprintf("%s:%s", codeOrToken, userID)
	}
	if err := uc.cache.Set(key, []byte(codeOrToken), constants.ForgotExpiredAt*time.Minute); err != nil {
		if err := uc.sessionRepo.CreateSession(session); err != nil {
			return ErrCreateSession
		}
	} else {
		go uc.sessionRepo.CreateSession(session)
	}
	return nil
}

func (uc *forgotPasswordUsecaseImpl) SendEmailForgotPassword(payload queue.PayloadI) (string, error) {
	go uc.sessionRepo.DeleteAllSessionsForgot(context.Background())
	return uc.qc.EnqueueAnyTask(payload)
}

func (uc *forgotPasswordUsecaseImpl) ForgotPassword(email, os string, method ForgotPasswordType) (ForgotPasswordRes, error) {
	var resForgotPassword ForgotPasswordRes
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return resForgotPassword, ErrUserNotFound
	}
	resForgotPassword.User = user.GetInfor()
	exp := time.Now().Add(constants.ForgotExpiredAt * time.Minute)
	switch method {
	case ForgotByCode:
		resForgotPassword.Code = uc.generateRandomCode(6)
		if err := uc.saveCodeOrToken(ForgotByCode, user.ID, resForgotPassword.Code, os, exp); err != nil {
			return resForgotPassword, err
		}
		return resForgotPassword, nil
	case ForgotByToken:
		code := uc.generateRandomCode(6)
		resForgotPassword.Token, err = uc.token.GenForgotPasswordToken(user.ID, code, exp)
		if err != nil {
			return resForgotPassword, err
		}
		if err := uc.saveCodeOrToken(ForgotByToken, user.ID, resForgotPassword.Token, os, exp); err != nil {
			return resForgotPassword, err
		}
		return resForgotPassword, nil
	}
	return resForgotPassword, ErrValidateForgotPassword
}

func (uc *forgotPasswordUsecaseImpl) generateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := int64(1)
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1
	num := r.Int63n(max-min+1) + min
	return strconv.FormatInt(num, 10)
}

func (uc *forgotPasswordUsecaseImpl) ForgotPasswordWithSaga(sagaID string, execute common.ExecuteSaga) error {
	ctx := context.Background()
	sagaTx := uc.saga.NewTransaction(sagaID, ctx)
	if err := execute(sagaTx.GetContext(), sagaTx); err != nil {
		return err
	}
	return sagaTx.Execute(sagaTx.GetContext(), sagaID)
}

func (uc *forgotPasswordUsecaseImpl) CompensateSendEmail(ctx context.Context, taskID string) error {
	return uc.qc.CancelTask(taskID)
}

func (uc *forgotPasswordUsecaseImpl) CompensateForgotPassword(ctx context.Context, data CompensateForgotPassword) error {
	switch data.Type {
	case ForgotByCode:
		key := fmt.Sprintf("%s:%s", data.Code, data.UserID)
		if err := uc.cache.Delete(key); err != nil {
			if err := uc.sessionRepo.DeleteSessionForgotByTokenAndIdUser(ctx, data.Code, data.UserID); err != nil {
				return err
			}
		}
		go func() {
			if err := uc.sessionRepo.DeleteSessionForgotByTokenAndIdUser(ctx, data.Code, data.UserID); err != nil {
				uc.log.Error("async delete failed: " + err.Error())
			}
		}()
		return nil

	case ForgotByToken:
		if err := uc.cache.Delete(data.Token); err != nil {
			if err := uc.sessionRepo.DeleteSessionForgotByToken(ctx, data.Token); err != nil {
				return err
			}
		}
		go func() {
			if err := uc.sessionRepo.DeleteSessionForgotByToken(ctx, data.Token); err != nil {
				uc.log.Error("async delete failed: " + err.Error())
			}
		}()
		return nil

	default:
		return ErrInvalidCompensateType
	}
}
