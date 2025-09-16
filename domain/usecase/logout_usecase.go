package usecase

import (
	"auth-service/domain/repository"
	"context"
	"errors"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/token"
)

var (
	ErrNotFoundSession = errors.New("không tìm thấy phiên làm việc")
)

type LogoutUsecase interface {
	VerifyToken(token string) error
	Logout(token string) error
}

type logoutUsecaseImpl struct {
	sessionRepo repository.SessionRepository
	token       token.TokenAuthorizeI
	cache       cache.CacheI
}

func NewLogoutUsecase(
	sessionRepo repository.SessionRepository,
	token token.TokenAuthorizeI,
	cache cache.CacheI,
) LogoutUsecase {
	return &logoutUsecaseImpl{
		sessionRepo,
		token,
		cache,
	}
}

func (l *logoutUsecaseImpl) VerifyToken(token string) error {
	_, err := l.cache.Get(token)
	if err != nil {
		return ErrNotFoundSession
	}
	_, err = l.token.VerifyAuthorizeToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (l *logoutUsecaseImpl) Logout(token string) error {
	if err := l.cache.Delete(token); err != nil {
		return err
	}
	if err := l.sessionRepo.DeleteSessionAuthByToken(context.Background(), token); err != nil {
		return err
	}
	return nil
}
