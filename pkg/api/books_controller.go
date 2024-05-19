package api

import (
	"context"

	"github.com/adeynack/finances/pkg/repository/gen/finances/public/model"
	. "github.com/adeynack/finances/pkg/repository/gen/finances/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
)

func (s *Service) GetBooks(ctx context.Context, request GetBooksRequestObject) (GetBooksResponseObject, error) {
	// todo: Move SQL part to service/repo
	stmt := SELECT(
		Books.AllColumns,
		Users.DisplayName,
	).FROM(
		Books.
			INNER_JOIN(Users, Users.ID.EQ(Books.OwnerID)),
	).ORDER_BY(Books.Name)

	type Result struct {
		model.Books
		Owner model.Users
	}
	var result []Result
	err := stmt.QueryContext(ctx, s.DB, &result)
	if err != nil {
		return GetBooks200JSONResponse{}, err
	}

	books := lo.Map(result, func(b Result, _ int) Book {
		return Book{
			CreatedAt:              b.CreatedAt,
			DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
			Id:                     b.ID.String(),
			Name:                   b.Name,
			OwnerId:                b.OwnerID.String(),
			OwnerDisplayName:       b.Owner.DisplayName,
			UpdatedAt:              b.UpdatedAt,
		}
	})
	return GetBooks200JSONResponse{Books: books}, nil
}
