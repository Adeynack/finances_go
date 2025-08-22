package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/platform/ctxval"
	"github.com/adeynack/finances/pkg/repository/gen/finances/public/model"
	. "github.com/adeynack/finances/pkg/repository/gen/finances/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
	slogctx "github.com/veqryn/slog-context"
)

func New() (R, error) {
	return &implementation{}, nil
}

type implementation struct{}

func (r *implementation) GetBooks(ctx context.Context) ([]apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	stmt := SELECT(
		Books.AllColumns,
		Users.DisplayName,
	).FROM(
		Books.
			INNER_JOIN(Users, Users.ID.EQ(Books.OwnerID)),
	).ORDER_BY(Books.Name)

	type BooksWithOwner struct {
		model.Books
		Owner model.Users
	}

	var booksFromDB []BooksWithOwner
	err = stmt.QueryContext(ctx, db, &booksFromDB)
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

	// Temporary code to prove the embedded transaction simulation works.
	InTransaction(ctx, db, func(ctx context.Context, db DB) (bool, error) {
		const query = `insert into books(created_at, updated_at, default_currency_iso_code, name, owner_id) values ($1, $2, $3, $4, $5)`
		_, err := db.ExecContext(ctx, query,
			time.Now(),
			time.Now(),
			"CAD",
			fmt.Sprintf("Foo %s", uuid.NewString()),
			"569bcfdd-4056-42cd-af9c-285fa5ce92c8",
		)
		return true, err
	})

	return booksForAPIResponse, nil
}

func (r *implementation) GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	stmt := SELECT(
		Exchanges.AllColumns,
		Splits.AllColumns,
	).FROM(
		Exchanges.
			LEFT_JOIN(Splits, Splits.ExchangeID.EQ(Exchanges.ID)),
	).ORDER_BY(Exchanges.Date)

	type ExchangesWithSplits struct {
		model.Exchanges
		Splits []model.Splits
	}

	var exchangesFromDB []ExchangesWithSplits
	err = stmt.QueryContext(ctx, db, &exchangesFromDB)
	if err != nil {
		return nil, fmt.Errorf("fetching exchanges from database: %w", err)
	}

	exchangesForAPIResponse := lo.Map(exchangesFromDB, func(e ExchangesWithSplits, _ int) apimodel.ExchangeWithSplits {
		return apimodel.ExchangeWithSplits{
			Cheque:      e.Cheque,
			CreatedAt:   e.CreatedAt,
			Date:        types.Date{Time: e.Date},
			Description: e.Description,
			Id:          e.ID.String(),
			Memo:        e.Memo,
			RegisterId:  e.RegisterID.String(),
			Splits: lo.Map(e.Splits, func(s model.Splits, _ int) apimodel.Split {
				return apimodel.Split{
					Amount:                int(s.Amount),
					CounterpartAmount:     nil,
					CreatedAt:             time.Time{},
					DestinationRegisterId: "",
					ExchangeId:            "",
					Id:                    "",
					Memo:                  new(string),
					Status:                "",
					UpdatedAt:             time.Time{},
				}
			}),
			Status:    apimodel.ExchangeStatus(e.Status),
			UpdatedAt: time.Time{},
		}
	})

	return exchangesForAPIResponse, nil
}
