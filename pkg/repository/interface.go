package repository

import (
	"context"
	"database/sql"

	"github.com/adeynack/finances/pkg/api/apimodel"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type R interface {
	GetBooks(ctx context.Context) ([]apimodel.Book, error)
	GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error)
}
