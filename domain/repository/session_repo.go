package repository

import (
	"auth-service/domain/entity"
	"context"
)

type SessionRepository interface {
	CreateSession(data entity.Session) error
	GetSessionAliveByToken(typeSession entity.SessionType, token string) (entity.Session, error)
	GetSessionAliveByTokenAndIdUser(typeSession entity.SessionType, token, idUser string) (entity.Session, error)
	GetSessionForgotAliveByTokenAndIdUser(token, idUser string) (entity.Session, error)
	TokenExists(token string) bool
	DeleteSessionByTypeAndUserID(ctx context.Context, sessionType entity.SessionType, userID string) error
	DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error
	DeleteSessionVerifyByUserID(ctx context.Context, userID string) error
	DeleteSessionAuthByToken(ctx context.Context, token string) error
	DeleteSessionVerifyByToken(ctx context.Context, token string) error
	DeleteSessionForgotByToken(ctx context.Context, token string) error
	DeleteAllSessionsExpired(ctx context.Context) error
	DeleteAllSessionsForgot(ctx context.Context) error
	DeleteSessionForgotByTokenAndIdUser(ctx context.Context, token, idUser string) error
	Tx(ctx context.Context) SessionRepository
}
