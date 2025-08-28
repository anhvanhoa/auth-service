package repository

import (
	"context"
)

type ManagerTransaction interface {
	RunInTransaction(fn func(ctx context.Context) error) error
	Begin() (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
