package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/adeynack/finances/pkg/api/apimodel"
	"github.com/adeynack/finances/pkg/repository/gen/finances/public/model"
	. "github.com/adeynack/finances/pkg/repository/gen/finances/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
)

// GetExchanges implements StrictServerInterface.
func (s *Service) GetExchanges(ctx context.Context, request GetExchangesRequestObject) (GetExchangesResponseObject, error) {
	exchangesFromDB, err := getExchangesFromRepo(ctx, s.DB)
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

	return GetExchanges200JSONResponse{Exchanges: exchangesForAPIResponse}, nil
}

type ExchangesWithSplits struct {
	model.Exchanges
	Splits []model.Splits
}

func getExchangesFromRepo(ctx context.Context, db *sql.DB) ([]ExchangesWithSplits, error) {
	stmt := SELECT(
		Exchanges.AllColumns,
		Splits.AllColumns,
	).FROM(
		Exchanges.
			LEFT_JOIN(Splits, Splits.ExchangeID.EQ(Exchanges.ID)),
	).ORDER_BY(Exchanges.Date)

	var result []ExchangesWithSplits
	err := stmt.QueryContext(ctx, db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
