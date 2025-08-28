package repo

import (
	"auth-service/domain/entity"
	"auth-service/domain/repository"
	"context"

	"github.com/go-pg/pg/v10"
)

type sessionRepositoryImpl struct {
	db pg.DBI
}

func NewSessionRepository(db *pg.DB) repository.SessionRepository {
	return &sessionRepositoryImpl{
		db: db,
	}
}

func (sr *sessionRepositoryImpl) CreateSession(data entity.Session) error {
	_, err := sr.db.Model(&data).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) GetSessionAliveByToken(typeSession entity.SessionType, token string) (entity.Session, error) {
	var session entity.Session
	err := sr.db.Model(&session).Where("token = ?", token).Where("type = ?", typeSession).
		Where("expired_at > NOW()").
		Select()
	if err != nil {
		return session, err
	}
	return session, nil
}

func (sr *sessionRepositoryImpl) GetSessionAliveByTokenAndIdUser(typeSession entity.SessionType, token, idUser string) (entity.Session, error) {
	var session entity.Session
	err := sr.db.Model(&session).
		Where("token = ?", token).
		Where("type = ?", typeSession).
		Where("user_id = ?", idUser).
		Where("expired_at > NOW()").
		Select()
	if err != nil {
		return session, err
	}
	return session, nil
}

func (sr *sessionRepositoryImpl) GetSessionForgotAliveByTokenAndIdUser(token, idUser string) (entity.Session, error) {
	return sr.GetSessionAliveByTokenAndIdUser(entity.SessionTypeForgot, token, idUser)
}

func (sr *sessionRepositoryImpl) TokenExists(token string) bool {
	count, err := sr.db.Model(&entity.Session{}).Where("token = ?", token).
		Where("expired_at > NOW()").
		Count()
	if err != nil {
		return false
	}
	return count > 0
}

func (sr *sessionRepositoryImpl) DeleteSessionVerifyByUserID(ctx context.Context, userID string) error {
	return sr.DeleteSessionByTypeAndUserID(ctx, entity.SessionTypeVerify, userID)
}

func (sr *sessionRepositoryImpl) DeleteSessionByTypeAndUserID(ctx context.Context, sessionType entity.SessionType, userID string) error {
	_, err := sr.db.ModelContext(ctx, &entity.Session{}).
		Where("type = ? AND user_id = ?", sessionType, userID).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) DeleteSessionByTypeAndToken(ctx context.Context, sessionType entity.SessionType, token string) error {
	_, err := sr.db.ModelContext(ctx, &entity.Session{}).
		Where("type = ? AND token = ?", sessionType, token).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) DeleteSessionAuthByToken(ctx context.Context, token string) error {
	return sr.DeleteSessionByTypeAndToken(ctx, entity.SessionTypeAuth, token)
}

func (sr *sessionRepositoryImpl) DeleteSessionVerifyByToken(ctx context.Context, token string) error {
	return sr.DeleteSessionByTypeAndToken(ctx, entity.SessionTypeVerify, token)
}

func (sr *sessionRepositoryImpl) DeleteSessionForgotByToken(ctx context.Context, token string) error {
	return sr.DeleteSessionByTypeAndToken(ctx, entity.SessionTypeForgot, token)
}

func (sr *sessionRepositoryImpl) DeleteAllSessionsExpired(ctx context.Context) error {
	_, err := sr.db.ModelContext(ctx, &entity.Session{}).
		Where("expired_at < NOW()").
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) DeleteSessionForgotByTokenAndIdUser(ctx context.Context, token, idUser string) error {
	_, err := sr.db.ModelContext(ctx, &entity.Session{}).
		Where("type = ? AND token = ? AND user_id = ?", entity.SessionTypeForgot, token, idUser).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) DeleteAllSessionsForgot(ctx context.Context) error {
	_, err := sr.db.ModelContext(ctx, &entity.Session{}).
		Where("type = ?", entity.SessionTypeForgot).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (sr *sessionRepositoryImpl) Tx(ctx context.Context) repository.SessionRepository {
	tx := getTx(ctx, sr.db)
	return &sessionRepositoryImpl{
		db: tx,
	}
}
