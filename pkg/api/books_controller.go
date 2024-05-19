package api

import (
	"context"

	"github.com/adeynack/finances/pkg/repository/gen/finances/public/model"
	. "github.com/adeynack/finances/pkg/repository/gen/finances/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/samber/lo"
)

func (s *Service) GetBooks(ctx context.Context, request GetBooksRequestObject) (GetBooksResponseObject, error) {
	stmt := SELECT(Books.AllColumns).FROM(Books).ORDER_BY(Books.Name) // todo: Move to service/repo
	var result []model.Books
	err := stmt.QueryContext(ctx, s.DB, &result)
	if err != nil {
		return GetBooks200JSONResponse{}, err
	}
	books := lo.Map(result, func(b model.Books, _ int) Book {
		return Book{
			CreatedAt:              b.CreatedAt,
			DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
			Id:                     b.ID.String(),
			Name:                   b.Name,
			OwnerId:                b.OwnerID.String(),
			UpdatedAt:              b.UpdatedAt,
		}
	})
	return GetBooks200JSONResponse{Books: books}, nil
}
