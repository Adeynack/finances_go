package app

import (
	"database/sql"
	"errors"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func MustConnectDatabase() *bun.DB {
	dsn, dsnPresent := os.LookupEnv("DATABASE_URL")
	if !dsnPresent {
		panic(errors.New("environment DATABASE_URL must be set"))
	}

	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	bunDB := bun.NewDB(db, pgdialect.New())

	return bunDB
}
