package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/platform/ctxval"
	"github.com/adeynack/finances/pkg/repository/dbmodel"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func New() (R, error) {
	return &implementation{}, nil
}

type implementation struct{}

func (r *implementation) errorIsEmptyResultSet(err error) bool {
	return errors.Is(err, sql.ErrNoRows) ||
		strings.Contains(err.Error(), "qrm: no rows in result set")
}

func (r *implementation) GetBooks(ctx context.Context) ([]apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	var books []dbmodel.Book
	err = db.NewSelect().
		Model(&books).
		Relation("Owner").
		Order("book.name", "book.id").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("quering GetBooks: %w", err)
	}

	booksForAPIResponse := lo.Map(books, func(b dbmodel.Book, _ int) apimodel.Book {
		return apimodel.Book{
			CreatedAt:              b.CreatedAt,
			DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
			Id:                     b.ID,
			Name:                   b.Name,
			OwnerDisplayName:       b.Owner.DisplayName,
			OwnerId:                b.OwnerID,
			UpdatedAt:              b.UpdatedAt,
		}
	})

	return booksForAPIResponse, nil
}

func (r *implementation) GetBookByID(ctx context.Context, bookId uuid.UUID) (*apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	var book dbmodel.Book
	err = db.NewSelect().
		Model(&book).
		Relation("Owner").
		Where("book.id = ?", bookId).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying BookByID: %w", err)
	}

	b := apimodel.Book{
		CreatedAt:              book.CreatedAt,
		DefaultCurrencyIsoCode: book.DefaultCurrencyIsoCode,
		Id:                     book.ID,
		Name:                   book.Name,
		OwnerDisplayName:       book.Owner.DisplayName,
		OwnerId:                book.OwnerID,
		UpdatedAt:              book.UpdatedAt,
	}

	return &b, nil
}

func (r *implementation) CreateBook(ctx context.Context, props apimodel.BookProperties) (apimodel.Book, error) {
	return apimodel.Book{}, errors.New("TODO")
	// db, err := r.resolveSqlcDb(ctx)
	// if err != nil {
	// 	return apimodel.Book{}, err
	// }

	// // Validate if the user exists
	// owner, err := r.GetUserByID(ctx, props.OwnerId)
	// if err != nil {
	// 	if r.errorIsEmptyResultSet(err) {
	// 		return apimodel.Book{}, fmt.Errorf("%w: owner does not exist", ErrValidation)
	// 	}

	// 	return apimodel.Book{}, fmt.Errorf("checking existence of book owner: %w", err)
	// }

	// book, err := db.CreateBook(ctx, sqlcdb.CreateBookParams{
	// 	CreatedAt:              time.Now(),
	// 	UpdatedAt:              time.Now(),
	// 	Name:                   props.Name,
	// 	OwnerID:                props.OwnerId,
	// 	DefaultCurrencyIsoCode: props.DefaultCurrencyIsoCode,
	// })
	// if err != nil {
	// 	return apimodel.Book{}, fmt.Errorf("querying CreateBook: %w", err)
	// }

	// result := apimodel.Book{
	// 	CreatedAt:              book.CreatedAt,
	// 	DefaultCurrencyIsoCode: book.DefaultCurrencyIsoCode,
	// 	Id:                     book.ID,
	// 	Name:                   book.Name,
	// 	OwnerDisplayName:       owner.DisplayName,
	// 	OwnerId:                book.OwnerID,
	// 	UpdatedAt:              book.UpdatedAt,
	// }

	// return result, nil
}

func (r *implementation) GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error) {
	return nil, errors.New("TODO")
	// db, err := r.resolveSqlcDb(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// exchanges, err := db.GetExchanges(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("querying GetExchanges: %w", err)
	// }

	// exchangeIds := lo.Map(exchanges, func(e sqlcdb.Exchange, _ int) uuid.UUID { return e.ID })
	// splits, err := db.GetSplitsForExchangeIds(ctx, exchangeIds)
	// if err != nil {
	// 	return nil, fmt.Errorf("querying GetSplitsForExchangeIds: %w", err)
	// }

	// splitsByExchangeId := lo.GroupByMap(splits, func(s sqlcdb.GetSplitsForExchangeIdsRow) (uuid.UUID, apimodel.Split) {
	// 	return s.ExchangeID, apimodel.Split{
	// 		Amount:                s.Amount,
	// 		CounterpartAmount:     int64Ptr(s.CounterpartAmount),
	// 		CreatedAt:             s.CreatedAt,
	// 		DestinationRegisterId: s.DestinationRegisterID,
	// 		ExchangeId:            s.ExchangeID,
	// 		Id:                    s.ID,
	// 		Memo:                  strPtr(s.Memo),
	// 		Status:                apimodel.ExchangeStatus(s.Status),
	// 		UpdatedAt:             s.UpdatedAt,
	// 	}
	// })

	// result := lo.Map(exchanges, func(e sqlcdb.Exchange, _ int) apimodel.ExchangeWithSplits {
	// 	splitsForExchange := splitsByExchangeId[e.ID]

	// 	return apimodel.ExchangeWithSplits{
	// 		Cheque:      strPtr(e.Cheque),
	// 		CreatedAt:   e.CreatedAt,
	// 		Date:        types.Date{Time: e.Date},
	// 		Description: e.Description,
	// 		Id:          e.ID,
	// 		Memo:        strPtr(e.Memo),
	// 		RegisterId:  e.RegisterID,
	// 		Splits:      splitsForExchange,
	// 		Status:      apimodel.ExchangeStatus(e.Status),
	// 		UpdatedAt:   e.UpdatedAt,
	// 	}
	// })

	// return result, nil
}

func (r *implementation) GetUserByID(ctx context.Context, userID uuid.UUID) (apimodel.User, error) {
	return apimodel.User{}, errors.New("TODO")
	// 	db, err := r.resolveSqlcDb(ctx)
	// 	if err != nil {
	// 		return apimodel.User{}, err
	// 	}

	// 	user, err := db.GetUserByID(ctx, userID)
	// 	if err != nil {
	// 		return apimodel.User{}, fmt.Errorf("querying GetUserByID: %w", err)
	// 	}

	// 	result := apimodel.User{
	// 		CreatedAt:   user.CreatedAt,
	// 		DisplayName: user.DisplayName,
	// 		Id:          user.ID,
	// 		UpdatedAt:   user.UpdatedAt,
	// 	}

	// return result, nil
}
