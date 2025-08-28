package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"
	"errors"

	"github.com/anhvanhoa/service-core/domain/cache"
	hashpass "github.com/anhvanhoa/service-core/domain/hash_pass"
	"github.com/anhvanhoa/service-core/domain/token"
)

type ResetPasswordByCodeUsecase interface {
	VerifySession(code, email string) (string, error)
	ResetPass(IdUser, Password, NewPassword string) error
}

type ResetPasswordByCodeUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	cache       cache.CacheI
	jwt         token.TokenForgotPasswordI
	hashPass    hashpass.HashPassI
}

var (
	ErrNotFoundUser   = errors.New("không tìm thấy người dùng")
	ErrHashPassword   = errors.New("không thể mã hóa mật khẩu")
	ErrUpdatePassword = errors.New("không thể cập nhật mật khẩu")
)

func NewResetPasswordCodeUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	cache cache.CacheI,
	token token.TokenForgotPasswordI,
	hashPass hashpass.HashPassI,
) ResetPasswordByCodeUsecase {
	return &ResetPasswordByCodeUsecaseImpl{
		userRepo,
		sessionRepo,
		cache,
		token,
		hashPass,
	}
}

func (uc *ResetPasswordByCodeUsecaseImpl) VerifySession(code, email string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrNotFoundUser
	}
	key := code + user.ID
	if _, err := uc.cache.Get(key); err != nil {
		if _, err := uc.sessionRepo.GetSessionForgotAliveByTokenAndIdUser(code, user.ID); err != nil {
			return "", ErrNotFoundSession
		}
	}
	go func() {
		uc.sessionRepo.DeleteSessionForgotByTokenAndIdUser(context.Background(), code, user.ID)
		uc.cache.Delete(key)
	}()
	return user.ID, nil
}

func (uc *ResetPasswordByCodeUsecaseImpl) ResetPass(IdUser, Password, ConfirmPassword string) error {
	ConfirmPassword, err := uc.hashPass.HashPassword(ConfirmPassword)
	if err != nil {
		return ErrHashPassword
	}

	if _, err = uc.userRepo.UpdateUser(IdUser, entity.User{Password: ConfirmPassword}); err != nil {
		return ErrUpdatePassword
	}
	return nil
}
