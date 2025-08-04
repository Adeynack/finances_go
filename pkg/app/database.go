package app

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func mustConnectDatabase() *sql.DB {
	dsn, dsnPresent := os.LookupEnv("DATABASE_URL")
	if !dsnPresent {
		panic(errors.New("environment DATABASE_URL must be set"))
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("opening connection to database: %w", err))
	}

	return db
}
