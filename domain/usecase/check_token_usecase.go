package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"errors"
)

var (
	ErrTokenInvalid = errors.New("phiên làm việc không hợp lệ hoặc đã hết hạn")
)

type CheckTokenUsecase interface {
	CheckToken(token string) (bool, error)
}

type checkTokenUsecaseImpl struct {
	sessionRepo repository.SessionRepository
}

func NewCheckTokenUsecase(sessionRepo repository.SessionRepository) CheckTokenUsecase {
	return &checkTokenUsecaseImpl{
		sessionRepo: sessionRepo,
	}
}

func (c *checkTokenUsecaseImpl) CheckToken(token string) (bool, error) {
	session, err := c.sessionRepo.GetSessionAliveByToken(entity.SessionTypeForgot, token)
	if err != nil {
		return false, ErrTokenInvalid
	}
	if session.Token == "" {
		return false, ErrTokenInvalid
	}
	return true, nil
}
