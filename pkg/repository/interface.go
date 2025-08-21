package repository

import (
	"context"

	"github.com/adeynack/finances/pkg/api/apimodel"
)

type R interface {
	GetBooks(ctx context.Context) ([]apimodel.Book, error)
	GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error)
}
