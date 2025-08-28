package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/token"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type VerifyAccountUsecase interface {
	VerifyRegister(t string) (*token.AuthClaims, error)
	GetUserById(id string) (entity.User, error)
	VerifyAccount(id string) error
}

type verifyAccountUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	token       token.TokenAuthI
	cache       cache.CacheI
}

func NewVerifyAccountUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	token token.TokenAuthI,
	cache cache.CacheI,
) VerifyAccountUsecase {
	return &verifyAccountUsecaseImpl{
		userRepo,
		sessionRepo,
		token,
		cache,
	}
}

func (u *verifyAccountUsecaseImpl) VerifyRegister(t string) (*token.AuthClaims, error) {
	if _, err := u.cache.Get(t); err != nil {
		if isExist := u.sessionRepo.TokenExists(t); !isExist {
			return nil, ErrTokenNotFound
		}
	} else {
		go func() {
			u.sessionRepo.DeleteSessionAuthByToken(context.Background(), t)
			u.cache.Delete(t)
		}()
	}

	data, err := u.token.VerifyAuthToken(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u *verifyAccountUsecaseImpl) GetUserById(id string) (entity.User, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *verifyAccountUsecaseImpl) VerifyAccount(id string) error {
	t := time.Now()
	user := entity.User{
		ID:         id,
		CodeVerify: "",
		Veryfied:   &t,
	}
	_, err := u.userRepo.UpdateUser(id, user)
	return err
}
