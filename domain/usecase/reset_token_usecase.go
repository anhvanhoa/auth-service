package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"

	"github.com/anhvanhoa/service-core/domain/cache"
	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
	"github.com/anhvanhoa/service-core/domain/token"
)

type ResetPasswordByTokenUsecase interface {
	VerifySession(token string) (string, error)
	ResetPass(IdUser, Password, NewPassword string) error
}

type ResetPasswordByTokenUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cache       cache.CacheI
	jwt         token.TokenForgotPasswordI
	hashPass    hashpass.HashPassI
}

func NewResetPasswordTokenUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cache cache.CacheI,
	token token.TokenForgotPasswordI,
	hashPass hashpass.HashPassI,
) ResetPasswordByTokenUsecase {
	return &ResetPasswordByTokenUsecaseImpl{
		userRepo,
		sessionRepo,
		cache,
		token,
		hashPass,
	}
}

func (uc *ResetPasswordByTokenUsecaseImpl) VerifySession(token string) (string, error) {
	if _, err := uc.cache.Get(token); err != nil {
		if isExist := uc.sessionRepo.TokenExists(token); !isExist {
			return "", ErrNotFoundSession
		}
	}
	go func() {
		uc.sessionRepo.DeleteSessionForgotByToken(context.Background(), token)
		uc.cache.Delete(token)
	}()
	claim, err := uc.jwt.VerifyForgotPasswordToken(token)
	if err != nil {
		return "", err
	}

	return claim.Data.Id, nil
}

func (uc *ResetPasswordByTokenUsecaseImpl) ResetPass(IdUser, Password, ConfirmPassword string) error {
	user, err := uc.userRepo.GetUserByID(IdUser)
	if err != nil {
		return ErrNotFoundUser
	}

	ConfirmPassword, err = uc.hashPass.HashPassword(ConfirmPassword)
	if err != nil {
		return ErrHashPassword
	}

	if _, err = uc.userRepo.UpdateUser(user.ID, entity.User{Password: ConfirmPassword}); err != nil {
		return ErrUpdatePassword
	}
	return nil
}
