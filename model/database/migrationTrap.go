package database

import (
	"context"
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

func FatalLogIfTrappedMigrationError(err error) bool {
	if trappedMigrationsError, ok := err.(*TrappedMigrationsError); ok {
		// TODO: Move to future `cmd/tool` or `cmd/dev` dev-ops binary
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
