package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	slogctx "github.com/veqryn/slog-context"
)

func MustConnectDatabase() *bun.DB {
	dsn, dsnPresent := os.LookupEnv("DATABASE_URL")
	if !dsnPresent {
		panic(errors.New("environment DATABASE_URL must be set"))
	}

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunDB := bun.NewDB(db, pgdialect.New())
	bunDB.AddQueryHook(&QueryHook{})

	return bunDB
}

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	switch event.Operation() {
	case "INSERT", "UPDATE":
		slogctx.Info(ctx, "Saving model",
			slog.String("model_type", fmt.Sprintf("%T", event.Model.Value())),
			slog.Any("model", event.Model.Value()),
		)
	}

	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	// fmt.Println(time.Since(event.StartTime), string(event.Query))
}
