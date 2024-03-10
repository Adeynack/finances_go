package database

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type TrappedMigrationsError struct {
	PendingMigrations []string
}

// Error implements error.
func (*TrappedMigrationsError) Error() string {
	return "pending migrations detected" // TODO: Instruct to use `cmd/tool` (still to be coded) to output the missing SQL migrations
}

// FatalLogIfTrappedMigrationError will call `log.Fatal` if the received err is
// a TrappedMigrationsError, listing the migrations SQL command that are missing
// to be in sync with the declared Gorm structure.
//
// It returns `true` if it failed; otherwise, `false` if it did not fail.
func FatalLogIfTrappedMigrationError(err error) bool {
	var trappedMigrationsError *TrappedMigrationsError
	if errors.As(err, &trappedMigrationsError) {
		log.Fatalf("Database migration is missing those elements to be in sync with the actual Gorm declared model:\n\n%s;\n\n", strings.Join(trappedMigrationsError.PendingMigrations, ";\n\n"))
		return true
	}

	return false
}

type migrationTrapLogger struct {
	logger.Interface
	PendingMigrations []string
}

// Trace implements logger.Interface.
func (l *migrationTrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	l.considerQuery(sql)
}

func (l *migrationTrapLogger) considerQuery(sql string) {
	upSql := strings.ToUpper(sql)
	if strings.HasPrefix(upSql, "SELECT ") {
		return
	}

	l.PendingMigrations = append(l.PendingMigrations, sql)
}
