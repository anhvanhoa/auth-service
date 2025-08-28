package repo

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type txContextKey struct{}

func getTx(ctx context.Context, db pg.DBI) pg.DBI {
	if tx, ok := ctx.Value(txContextKey{}).(*pg.Tx); ok {
		return tx
	}
	return db
}
