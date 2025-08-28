package usecase

import (
	"auth-service/constants"
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"math/rand"
	"time"

	"github.com/anhvanhoa/service-core/common"
	"github.com/anhvanhoa/service-core/domain/cache"
	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"

	"github.com/anhvanhoa/service-core/domain/goid"
	"github.com/anhvanhoa/service-core/domain/queue"
	"github.com/anhvanhoa/service-core/domain/saga"
	"github.com/anhvanhoa/service-core/domain/token"
)

type ResRegister struct {
	UserInfor entity.UserInfor
	Token     string
}

type RegisterReq struct {
	Email           string
	FullName        string
	Password        string
	ConfirmPassword string
	Code            string
}

type RegisterUsecase interface {
	CheckUserExist(email string) (entity.User, error)
	hashPassword(password string) (string, error)
	Register(user RegisterReq, os string, exp time.Time) (ResRegister, error)
	RegisterWithSaga(sagaID string, execute common.ExecuteSaga) error
	GengerateCode(length int8) string
	createOrUpdateUser(user RegisterReq, ctx context.Context) (entity.UserInfor, error)
	saveToken(token string, id string, os string) error
	SendMail(payload queue.PayloadI) (string, error)
	CompensateRegister(ctx context.Context, userId string, token string) error
	CompensateSendMail(ctx context.Context, taskID string) error
}

type registerUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwt         token.TokenAuthI
	tx          repository.ManagerTransaction
	saga        saga.SagaManager
	goid        goid.GoUUID
	hashPass    hashpass.HashPassI
	cache       cache.CacheI
	qc          queue.QueueClient
}

func NewRegisterUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwt token.TokenAuthI,
	tx repository.ManagerTransaction,
	goid goid.GoUUID,
	hashPass hashpass.HashPassI,
	cache cache.CacheI,
	queue queue.QueueClient,
	saga saga.SagaManager,
) RegisterUsecase {
	return &registerUsecaseImpl{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tx:          tx,
		jwt:         jwt,
		goid:        goid,
		hashPass:    hashPass,
		cache:       cache,
		qc:          queue,
		saga:        saga,
	}
}

func (uc *registerUsecaseImpl) CheckUserExist(email string) (entity.User, error) {
	return uc.userRepo.GetUserByEmail(email)
}

func (uc *registerUsecaseImpl) GengerateCode(length int8) string {
	const digits = "0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[r.Intn(len(digits))]
	}
	return string(result)
}

func (uc *registerUsecaseImpl) hashPassword(password string) (string, error) {
	return uc.hashPass.HashPassword(password)
}

func (uc *registerUsecaseImpl) Register(user RegisterReq, os string, exp time.Time) (ResRegister, error) {
	res := ResRegister{}
	err := uc.tx.RunInTransaction(func(ctx context.Context) error {
		var err error
		if res.UserInfor, err = uc.createOrUpdateUser(user, ctx); err != nil {
			return err
		}
		if err = uc.sessionRepo.DeleteSessionVerifyByUserID(ctx, res.UserInfor.ID); err != nil {
			return err
		}
		if res.Token, err = uc.jwt.GenAuthToken(res.UserInfor.ID, user.Code, exp); err != nil {
			return err
		}
		if err = uc.saveToken(res.Token, res.UserInfor.ID, os); err != nil {
			return err
		}
		return nil
	})
	return res, err
}

func (uc *registerUsecaseImpl) createOrUpdateUser(user RegisterReq, ctx context.Context) (entity.UserInfor, error) {
	var userInfo entity.UserInfor
	var err error
	id := uc.goid.Gen()
	newUser := entity.User{
		ID:         id,
		Email:      user.Email,
		Password:   user.Password,
		FullName:   user.FullName,
		CodeVerify: user.Code,
	}
	if newUser.Password, err = uc.hashPassword(newUser.Password); err != nil {
		return userInfo, err
	}

	if isExist, err := uc.userRepo.CheckUserExist(newUser.Email); err != nil {
		return userInfo, err
	} else if !isExist {
		if userInfo, err = uc.userRepo.Tx(ctx).CreateUser(newUser); err != nil {
			return userInfo, err
		}
	} else {
		newUser.ID = "" // Nó sẽ không cập nhật ID bởi ID là khóa chính | thêm cho dễ hiểu
		if ok, err := uc.userRepo.Tx(ctx).UpdateUserByEmail(newUser.Email, newUser); err != nil {
			return userInfo, err
		} else if u, err := uc.userRepo.GetUserByEmail(newUser.Email); ok && err == nil {
			userInfo = u.GetInfor()
		} else {
			return userInfo, err
		}
	}
	return userInfo, nil
}

func (uc *registerUsecaseImpl) saveToken(token string, userId string, os string) error {
	session := entity.Session{
		Token:     token,
		UserID:    userId,
		Os:        os,
		Type:      entity.SessionTypeVerify,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(constants.VerifyExpiredAt * time.Second),
	}
	if err := uc.cache.Set(token, []byte(token), constants.VerifyExpiredAt*time.Second); err != nil {
		if err := uc.sessionRepo.CreateSession(session); err != nil {
			return err
		}
	} else {
		go uc.sessionRepo.CreateSession(session)
	}
	return nil
}

func (uc *registerUsecaseImpl) SendMail(payload queue.PayloadI) (string, error) {
	return uc.qc.EnqueueAnyTask(string(constants.QUEUE_MAIL), payload)
}

func (uc *registerUsecaseImpl) RegisterWithSaga(sagaID string, execute common.ExecuteSaga) error {
	ctx := context.Background()
	sagaTx := uc.saga.NewTransaction(sagaID, ctx)
	if err := execute(sagaTx.GetContext(), sagaTx); err != nil {
		return err
	}
	return sagaTx.Execute(sagaTx.GetContext(), sagaID)
}

func (uc *registerUsecaseImpl) CompensateRegister(ctx context.Context, userID string, token string) error {
	if err := uc.userRepo.DeleteByID(ctx, userID); err != nil {
		return err
	}
	go uc.sessionRepo.DeleteSessionVerifyByUserID(ctx, userID)
	return uc.cache.Delete(token)
}

func (uc *registerUsecaseImpl) CompensateSendMail(ctx context.Context, taskID string) error {
	return uc.qc.CancelTask(queue.TypeMail, taskID)
}
