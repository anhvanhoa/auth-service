package entity

import (
	"time"
)

type SessionType string

const (
	SessionTypeAuth   SessionType = "authorization"
	SessionTypeForgot SessionType = "forgot"
	SessionTypeReset  SessionType = "reset"
	SessionTypeVerify SessionType = "verify"
)

type Session struct {
	tableName struct{}    `pg:"sessions,alias:s"`
	Token     string      `pg:"token,pk"`
	UserID    string      `pg:"user_id,pk"`
	User      *User       `pg:"rel:has-one"`
	Type      SessionType `pg:"type"`
	Os        string      `pg:"os"`
	ExpiredAt time.Time   `pg:"expired_at"`
	CreatedAt time.Time   `pg:"created_at"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiredAt)
}

func (s *Session) NameTable() any {
	return s.tableName
}
