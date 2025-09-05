package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/platform/ctxval"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	const query = `-- GetBooks
		select
			books.created_at,
			books.default_currency_iso_code,
			books.id,
			books.name,
			books.owner_id,
			users.display_name,
			books.updated_at
		from books
		inner join users on users.id = books.owner_id
		order by books.name, books.id
	`

	booksForAPIResponse := make([]apimodel.Book, 0)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying GetBooks: %w", err)
	}

	for rows.Next() {
		var b apimodel.Book
		err = rows.Scan(
			&b.CreatedAt,
			&b.DefaultCurrencyIsoCode,
			&b.Id,
			&b.Name,
			&b.OwnerId,
			&b.OwnerDisplayName,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning GetBooks: %w", err)
		}

		booksForAPIResponse = append(booksForAPIResponse, b)
	}

	return booksForAPIResponse, nil
}

func (r *implementation) GetBookByID(ctx context.Context, bookId uuid.UUID) (*apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	const query = `-- GetBookByID
		select
			books.created_at,
			books.default_currency_iso_code,
			books.id,
			books.name,
			books.owner_id,
			users.display_name,
			books.updated_at
		from books
		inner join users on users.id = books.owner_id
		where books.id = $1
		limit 1
	`

	row := db.QueryRowContext(ctx, query, bookId)
	if row.Err() != nil {
		return nil, fmt.Errorf("querying GetBookByID: %w", row.Err())
	}

	var b apimodel.Book
	err = row.Scan(
		&b.CreatedAt,
		&b.DefaultCurrencyIsoCode,
		&b.Id,
		&b.Name,
		&b.OwnerId,
		&b.OwnerDisplayName,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning GetBookByID: %w", err)
	}

	return &b, nil
}

func (r *implementation) CreateBook(ctx context.Context, props apimodel.BookProperties) (apimodel.Book, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return apimodel.Book{}, err
	}

	// Validate if the user exists
	owner, err := r.GetUserByID(ctx, props.OwnerId)
	if err != nil {
		if r.errorIsEmptyResultSet(err) {
			return apimodel.Book{}, fmt.Errorf("%w: owner does not exist", ErrValidation)
		}

		return apimodel.Book{}, fmt.Errorf("checking existence of book owner: %w", err)
	}

	const query = `-- CreateBook
		insert into books(
			created_at,
			updated_at,
			name,
			owner_id,
			default_currency_iso_code
		) values (
			$1, -- created_at,
			$2, -- updated_at,
			$3, -- name,
			$4, -- owner_id,
			$5  -- default_currency_iso_code
		)
		returning
			id,
			created_at,
			updated_at,
			name,
			owner_id,
			default_currency_iso_code
	`

	row := db.QueryRowContext(ctx, query,
		time.Now(),
		time.Now(),
		props.Name,
		props.OwnerId,
		props.DefaultCurrencyIsoCode,
	)
	if row.Err() != nil {
		return apimodel.Book{}, fmt.Errorf("querying CreateBook: %w", row.Err())
	}

	result := apimodel.Book{
		OwnerDisplayName: owner.DisplayName,
	}
	err = row.Scan(
		&result.Id,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Name,
		&result.OwnerId,
		&result.DefaultCurrencyIsoCode,
	)
	if err != nil {
		return apimodel.Book{}, fmt.Errorf("scanning CreateBook: %w", err)
	}

	return result, nil
}

func (r *implementation) GetExchanges(ctx context.Context) ([]apimodel.Exchange, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	const query = `-- GetExchanges
		select
			exchanges.cheque,
			exchanges.created_at,
			exchanges.date,
			exchanges.description,
			exchanges.id,
			exchanges.memo,
			exchanges.register_id,
			exchanges.status,
			exchanges.updated_at
		from
			exchanges
		order by
			exchanges.date,
			exchanges.created_at,
			exchanges.id
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying GetExchanges: %w", err)
	}

	result := make([]apimodel.Exchange, 0)
	var exchange apimodel.Exchange

	for rows.Next() {
		err = rows.Scan(
			&exchange.Cheque,
			&exchange.CreatedAt,
			&exchange.Date.Time,
			&exchange.Description,
			&exchange.Id,
			&exchange.Memo,
			&exchange.RegisterId,
			&exchange.Status,
			&exchange.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning GetExchanges: %w", err)
		}

		result = append(result, exchange)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("getting next row for GetExchangesWithSplits: %w", err)
	}

	return result, nil
}

func (r *implementation) GetSplitsForExchangeIds(ctx context.Context, exchangeIds []uuid.UUID) ([]apimodel.Split, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return nil, err
	}

	const query = `-- GetSplitsForExchangeIds
		select
			splits.amount,
			splits.counterpart_amount,
			splits.created_at,
			splits.destination_register_id,
			splits.exchange_id,
			splits.id,
			splits.memo,
			splits.status,
			splits.updated_at
		from
			splits
		where
			splits.exchange_id = ANY($1)
		order by
			splits.created_at,
			splits.id
	`
	rows, err := db.QueryContext(ctx, query, pq.Array(exchangeIds))
	if err != nil {
		return nil, fmt.Errorf("querying GetSplitsForExchangeIds: %w", err)
	}

	result := make([]apimodel.Split, 0)
	var split apimodel.Split

	for rows.Next() {
		err = rows.Scan(
			&split.Amount,
			&split.CounterpartAmount,
			&split.CreatedAt,
			&split.DestinationRegisterId,
			&split.ExchangeId,
			&split.Id,
			&split.Memo,
			&split.Status,
			&split.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning GetSplitsForExchangeIds: %w", err)
		}

		result = append(result, split)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("getting next row for GetSplitsForExchangeIds: %w", err)
	}

	return result, nil
}

func (r *implementation) GetExchangesWithSplits(ctx context.Context) ([]apimodel.ExchangeWithSplits, error) {
	exchanges, err := r.GetExchanges(ctx)
	if err != nil {
		return nil, err
	}

	exchangeIds := lo.Map(exchanges, func(e apimodel.Exchange, _ int) uuid.UUID { return e.Id })
	splits, err := r.GetSplitsForExchangeIds(ctx, exchangeIds)
	if err != nil {
		return nil, err
	}

	splitsByExchangeId := lo.GroupBy(splits, func(s apimodel.Split) uuid.UUID { return s.ExchangeId })

	result := lo.Map(exchanges, func(e apimodel.Exchange, _ int) apimodel.ExchangeWithSplits {
		splitsForExchange := splitsByExchangeId[e.Id]

		return apimodel.ExchangeWithSplits{
			Cheque:      e.Cheque,
			CreatedAt:   e.CreatedAt,
			Date:        e.Date,
			Description: e.Description,
			Id:          e.Id,
			Memo:        e.Memo,
			RegisterId:  e.RegisterId,
			Splits:      splitsForExchange,
			Status:      e.Status,
			UpdatedAt:   e.UpdatedAt,
		}
	})

	return result, nil
}

func (r *implementation) GetUserByID(ctx context.Context, userID uuid.UUID) (apimodel.User, error) {
	db, err := ctxval.Resolve[DB](ctx)
	if err != nil {
		return apimodel.User{}, err
	}

	const query = `-- GetUserByID
		select
			users.created_at,
			users.display_name,
			users.id,
			users.updated_at
		from
			users
		where
			users.id = $1
		limit 1
	`

	row := db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return apimodel.User{}, fmt.Errorf("querying GetUserByID: %w", row.Err())
	}

	var user apimodel.User
	err = row.Scan(
		&user.CreatedAt,
		&user.DisplayName,
		&user.Id,
		&user.UpdatedAt,
	)
	if err != nil {
		return apimodel.User{}, fmt.Errorf("scanning GetUserByID: %w", err)
	}

	return user, nil
}
