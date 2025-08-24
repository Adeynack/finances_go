package repository

import (
	"context"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/google/uuid"
)

type R interface {
	GetBooks(ctx context.Context) ([]apimodel.Book, error)
	GetBookByID(ctx context.Context, bookId uuid.UUID) (*apimodel.Book, error)
	CreateBook(ctx context.Context, body apimodel.BookProperties) (apimodel.Book, error)
	GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error)
}
