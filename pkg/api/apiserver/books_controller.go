package apiserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/repository/gen/finances/public/model"
	. "github.com/adeynack/finances/pkg/repository/gen/finances/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
	slogctx "github.com/veqryn/slog-context"
)

func (s *Service) GetBooks(ctx context.Context, request GetBooksRequestObject) (GetBooksResponseObject, error) {
	booksFromDB, err := getBooksFromRepo(ctx, s.DB)
	if err != nil {
		return nil, fmt.Errorf("fetching books from database: %w", err)
	}
	slogctx.Info(ctx, fmt.Sprintf("getBooksFromRepo returned %d entries", len(booksFromDB)))

	booksForAPIResponse := lo.Map(booksFromDB, func(b BooksWithOwner, _ int) apimodel.Book {
		return apimodel.Book{
			CreatedAt:              b.CreatedAt,
			DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
			Id:                     b.ID.String(),
			Name:                   b.Name,
			OwnerId:                b.OwnerID.String(),
			OwnerDisplayName:       b.Owner.DisplayName,
			UpdatedAt:              b.UpdatedAt,
		}
	})

	return GetBooks200JSONResponse{Books: booksForAPIResponse}, nil
}

type BooksWithOwner struct {
	model.Books
	Owner model.Users
}

func getBooksFromRepo(ctx context.Context, db *sql.DB) ([]BooksWithOwner, error) {
	stmt := SELECT(
		Books.AllColumns,
		Users.DisplayName,
	).FROM(
		Books.
			INNER_JOIN(Users, Users.ID.EQ(Books.OwnerID)),
	).ORDER_BY(Books.Name)

	var result []BooksWithOwner
	err := stmt.QueryContext(ctx, db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
