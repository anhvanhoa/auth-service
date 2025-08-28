package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"time"

	"github.com/anhvanhoa/service-core/domain/cache"
	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
	"github.com/anhvanhoa/service-core/domain/token"
)

type LoginUsecase interface {
	GetUserByEmailOrPhone(val string) (entity.User, error)
	CheckHashPassword(password, hash string) bool
	GengerateAccessToken(id, fullName, email string, exp time.Time) (string, error)
	GengerateRefreshToken(id, fullName, email string, exp time.Time, os string) (string, error)
}

type loginUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtAccess   token.TokenAuthorizeI
	jwtRefresh  token.TokenAuthorizeI
	hassPass    hashpass.HashPassI
	cache       cache.CacheI
}

func NewLoginUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	jwtAccess token.TokenAuthorizeI,
	jwtRefresh token.TokenAuthorizeI,
	hassPass hashpass.HashPassI,
	cache cache.CacheI,
) LoginUsecase {
	return &loginUsecaseImpl{
		userRepo,
		sessionRepo,
		jwtAccess,
		jwtRefresh,
		hassPass,
		cache,
	}
}

func (uc *loginUsecaseImpl) GetUserByEmailOrPhone(val string) (entity.User, error) {
	return uc.userRepo.GetUserByEmailOrPhone(val)
}

func (uc *loginUsecaseImpl) CheckHashPassword(password, hash string) bool {
	mach, err := uc.hassPass.VerifyPassword(hash, password)
	if err != nil {
		return false
	}
	return mach
}

func (uc *loginUsecaseImpl) GengerateAccessToken(id, fullName, email string, exp time.Time) (string, error) {
	return uc.jwtAccess.GenAuthorizeToken(id, fullName, email, exp)
}

func (uc *loginUsecaseImpl) GengerateRefreshToken(id, fullName, email string, exp time.Time, os string) (string, error) {
	token, err := uc.jwtRefresh.GenAuthorizeToken(id, fullName, email, exp)
	if err != nil {
		return "", err
	}

	session := entity.Session{
		Token:     token,
		UserID:    id,
		Os:        os,
		Type:      (entity.SessionTypeAuth),
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
