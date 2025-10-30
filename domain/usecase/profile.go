package usecase

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"

	"github.com/anhvanhoa/service-core/domain/cache"
	"github.com/anhvanhoa/service-core/domain/user_context"
)

type ProfileUsecase interface {
	Execute(ctx context.Context, token string) (*entity.UserInfor, error)
}

type profileUsecaseImpl struct {
	userRepo repository.UserRepository
	cache    cache.CacheI
}

func NewProfileUsecase(userRepo repository.UserRepository, cache cache.CacheI) ProfileUsecase {
	return &profileUsecaseImpl{
		userRepo: userRepo,
		cache:    cache,
	}
}

func (uc *profileUsecaseImpl) Execute(ctx context.Context, token string) (*entity.UserInfor, error) {
	userData, err := uc.cache.Get(token)
	if err != nil || userData == nil {
		return nil, ErrUserNotFound
	}
	user := user_context.NewUserContext()
	user.FromBytes(userData)

	userEntity, err := uc.userRepo.GetUserByID(user.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	userInfor := userEntity.GetInfor()
	return &userInfor, nil
}
