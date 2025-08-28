package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"time"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/token"
)

type RefreshUsecase interface {
	GetSessionByToken(token string) (entity.Session, error)
	ClearSessionExpired() error
	VerifyToken(token string) (*token.AuthorizeClaims, error)
	GengerateAccessToken(id, fullName, email string, exp time.Time) (string, error)
	GengerateRefreshToken(id, fullName, email string, exp time.Time, os string) (string, error)
}

type refreshUsecaseImpl struct {
	sessionRepo repository.SessionRepository
	access      token.TokenAuthorizeI
	refresh     token.TokenAuthorizeI
	cache       cache.CacheI
}

func NewRefreshUsecase(
	sessionRepo repository.SessionRepository,
	access token.TokenAuthorizeI,
	refresh token.TokenAuthorizeI,
	cache cache.CacheI,
) RefreshUsecase {
	return &refreshUsecaseImpl{
		sessionRepo: sessionRepo,
		access:      access,
		refresh:     refresh,
		cache:       cache,
	}
}

func (uc *refreshUsecaseImpl) GetSessionByToken(token string) (entity.Session, error) {
	session, err := uc.sessionRepo.GetSessionAliveByToken(entity.SessionTypeAuth, token)
	if err != nil {
		return entity.Session{}, err
	}
	if err := uc.sessionRepo.DeleteSessionAuthByToken(context.Background(), token); err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (uc *refreshUsecaseImpl) ClearSessionExpired() error {
	if err := uc.sessionRepo.DeleteAllSessionsExpired(context.Background()); err != nil {
		return err
	}
	return nil
}

func (uc *refreshUsecaseImpl) VerifyToken(token string) (*token.AuthorizeClaims, error) {
	claims, err := uc.refresh.VerifyAuthorizeToken(token)
	if err != nil {
		return claims, err
	}
	return claims, nil
}

func (uc *refreshUsecaseImpl) GengerateAccessToken(id, fullName, email string, exp time.Time) (string, error) {
	return uc.access.GenAuthorizeToken(id, fullName, email, exp)
}

func (uc *refreshUsecaseImpl) GengerateRefreshToken(id, fullName, email string, exp time.Time, os string) (string, error) {
	token, err := uc.refresh.GenAuthorizeToken(id, fullName, email, exp)
	if err != nil {
		return "", err
	}

	session := entity.Session{
		Token:     token,
		UserID:    id,
		Os:        os,
		Type:      entity.SessionTypeAuth,
		ExpiredAt: exp,
		CreatedAt: time.Now(),
	}

	if err := uc.cache.Set(token, []byte(id), time.Until(exp)); err != nil {
		if err := uc.sessionRepo.CreateSession(session); err != nil {
			return "", err
		}
	} else {
		go uc.sessionRepo.CreateSession(session)
	}
	return token, nil
}
