package usecase

import (
	"auth-service/domain/repository"
	"errors"
)

var (
	ErrCodeInvalid  = errors.New("mã xác thực không hợp lệ")
	ErrUserNotFound = errors.New("không tìm thấy người dùng với email này")
)

type CheckCodeUsecase interface {
	CheckCode(code, email string) (bool, error)
}

type checkCodeUsecaseImpl struct {
	userRepo repository.UserRepository
	session  repository.SessionRepository
}

func NewCheckCodeUsecase(userRepo repository.UserRepository, session repository.SessionRepository) CheckCodeUsecase {
	return &checkCodeUsecaseImpl{
		userRepo: userRepo,
		session:  session,
	}
}

func (c *checkCodeUsecaseImpl) CheckCode(code, email string) (bool, error) {
	user, err := c.userRepo.GetUserByEmail(email)
	if err != nil {
		return false, ErrUserNotFound
	}

	session, err := c.session.GetSessionForgotAliveByTokenAndIdUser(code, user.ID)
	if err != nil {
		return false, ErrCodeInvalid
	}
	if session.Token == "" {
		return false, ErrCodeInvalid
	}
	return true, nil
}
