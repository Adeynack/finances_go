package db

import (
	"context"
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

type migrationTrapLogger struct {
	PendingMigrations []string
}

// Error implements logger.Interface.
func (l *migrationTrapLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	panic("unimplemented")
}

// Info implements logger.Interface.
func (l *migrationTrapLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	panic("unimplemented")
}

// LogMode implements logger.Interface.
func (l *migrationTrapLogger) LogMode(logger.LogLevel) logger.Interface {
	panic("unimplemented")
}

// Trace implements logger.Interface.
func (l *migrationTrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	l.considerQuery(sql)
}

// Warn implements logger.Interface.
func (l *migrationTrapLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	panic("unimplemented")
}

func (l *migrationTrapLogger) considerQuery(sql string) {
	upSql := strings.ToUpper(sql)
	if strings.HasPrefix(upSql, "SELECT ") {
		return
	}

	l.PendingMigrations = append(l.PendingMigrations, sql)
}
