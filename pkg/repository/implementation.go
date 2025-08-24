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
	).ORDER_BY(Books.Name, Books.ID)

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
			Id:                     b.ID,
			Name:                   b.Name,
			OwnerId:                b.OwnerID,
			OwnerDisplayName:       b.Owner.DisplayName,
			UpdatedAt:              b.UpdatedAt,
		}
	})

	// // Temporary code to prove the embedded transaction simulation works.
	// InTransaction(ctx, db, func(ctx context.Context, db DB) (bool, error) {
	// 	const query = `insert into books(created_at, updated_at, default_currency_iso_code, name, owner_id) values ($1, $2, $3, $4, $5)`
	// 	_, err := db.ExecContext(ctx, query,
	// 		time.Now(),
	// 		time.Now(),
	// 		"CAD",
	// 		fmt.Sprintf("Foo %s", uuid.NewString()),
	// 		"569bcfdd-4056-42cd-af9c-285fa5ce92c8",
	// 	)
	// 	return true, err
	// })

	return booksForAPIResponse, nil
}

func (r *implementation) GetBookByID(ctx context.Context, bookId uuid.UUID) (*apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	stmt := SELECT(
		Books.AllColumns,
		Users.DisplayName.AS("user_display_name"),
	).FROM(
		Books.
			INNER_JOIN(Users, Users.ID.EQ(Books.OwnerID)),
	).
		WHERE(
			Books.ID.EQ(UUID(bookId)),
		)

	var book struct {
		model.Books
		UserDisplayName string
	}
	err = stmt.QueryContext(ctx, db, &book)
	if err != nil {
		return nil, fmt.Errorf("fetching the book by its ID: %w", err)
	}

	return &apimodel.Book{
		CreatedAt:              book.CreatedAt,
		DefaultCurrencyIsoCode: book.DefaultCurrencyIsoCode,
		Id:                     book.ID,
		Name:                   book.Name,
		OwnerDisplayName:       book.UserDisplayName,
		OwnerId:                book.OwnerID,
		UpdatedAt:              book.UpdatedAt,
	}, nil
}

func (r *implementation) CreateBook(ctx context.Context, props apimodel.BookProperties) (apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return apimodel.Book{}, err
	}

	stmt := Books.INSERT(
		Books.CreatedAt,
		Books.UpdatedAt,
		Books.Name,
		Books.OwnerID,
		Books.DefaultCurrencyIsoCode,
	).VALUES(
		time.Now(),
		time.Now(),
		props.Name,
		props.OwnerId,
		props.DefaultCurrencyIsoCode,
	).RETURNING(
		Books.AllColumns,
	)

	var insertedBooks []model.Books
	err = stmt.QueryContext(ctx, db, &insertedBooks)
	if err != nil {
		return apimodel.Book{}, fmt.Errorf("inserting a new book: %w", err)
	}
	if len(insertedBooks) != 1 {
		return apimodel.Book{}, fmt.Errorf("expected 1 book to be inserted, got %d", len(insertedBooks))
	}
	book := insertedBooks[0]

	result := apimodel.Book{
		CreatedAt:              book.CreatedAt,
		DefaultCurrencyIsoCode: book.DefaultCurrencyIsoCode,
		Id:                     book.ID,
		Name:                   book.Name,
		OwnerDisplayName:       "TODO",
		OwnerId:                book.OwnerID,
		UpdatedAt:              book.UpdatedAt,
	}

	return result, nil
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
			Id:          e.ID,
			Memo:        e.Memo,
			RegisterId:  e.RegisterID,
			Splits: lo.Map(e.Splits, func(s model.Splits, _ int) apimodel.Split {
				return apimodel.Split{
					Amount:                s.Amount,
					CounterpartAmount:     s.CounterpartAmount,
					CreatedAt:             s.CreatedAt,
					DestinationRegisterId: s.DestinationRegisterID,
					ExchangeId:            s.ExchangeID,
					Id:                    s.ID,
					Memo:                  s.Memo,
					Status:                apimodel.ExchangeStatus(s.Status),
					UpdatedAt:             s.UpdatedAt,
				}
			}),
			Status:    apimodel.ExchangeStatus(e.Status),
			UpdatedAt: time.Time{},
		}
	})

	return exchangesForAPIResponse, nil
}
