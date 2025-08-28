package entity

import (
	"time"
)

type UserRole struct {
	tableName struct{}  `pg:"user_roles,alias:ur"`
	UserID    string    `pg:"user_id,pk"`
	RoleID    string    `pg:"role_id,pk"`
	CreatedBy string    `pg:"created_by"`
	CreatedAt time.Time `pg:"created_at"`
}

func (ur *UserRole) GetNameTable() any {
	return ur.tableName
}
