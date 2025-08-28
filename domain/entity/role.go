package entity

import (
	"time"

	"github.com/anhvanhoa/service-core/common"
)

type Role struct {
	tableName   struct{}      `pg:"roles,alias:r"`
	ID          string        `pg:"id,pk"`
	Name        string        `pg:"name"`
	Description string        `pg:"description"`
	Variant     string        `pg:"variant"`
	Status      common.Status `pg:"status"`
	CreatedBy   string        `pg:"created_by"`
	CreatedAt   time.Time     `pg:"created_at"`
	UpdatedAt   *time.Time    `pg:"updated_at"`
}

func (r *Role) GetNameTable() any {
	return r.tableName
}
