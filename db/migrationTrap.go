package db

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TrappedMigrationsError struct {
	PendingMigrations []string
}

// Error implements error.
func (*TrappedMigrationsError) Error() string {
	return "pending migrations detected" // TODO: Instruct to use `cmd/tool` (still to be coded) to output the missing SQL migrations
}

type MigrationTrapDialector struct {
	gorm.Dialector
	migrator *MigrationTrapMigrator
}

func NewMigrationTrapDialector(baseDialector gorm.Dialector) *MigrationTrapDialector {
	return &MigrationTrapDialector{Dialector: baseDialector}
}

func (d *MigrationTrapDialector) Migrator(db *gorm.DB) gorm.Migrator {
	if d.migrator == nil {
		d.migrator = &MigrationTrapMigrator{
			Migrator: d.Dialector.Migrator(db),
			db:       db,
		}
	}
	return d.migrator
}

type migrationTrapLogger struct {
	BaseLogger        logger.Interface
	PendingMigrations []string
}

// Error implements logger.Interface.
func (l *migrationTrapLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	l.BaseLogger.Error(ctx, fmt, args...)
}

// Info implements logger.Interface.
func (l *migrationTrapLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	l.BaseLogger.Info(ctx, fmt, args...)
}

// LogMode implements logger.Interface.
func (l *migrationTrapLogger) LogMode(logger.LogLevel) logger.Interface {
	panic("unimplemented")
}

// Trace implements logger.Interface.
func (l *migrationTrapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	l.BaseLogger.Trace(
		ctx,
		begin,
		func() (sql string, rowsAffected int64) {
			sql, rowsAffected = fc()
			l.considerQuery(sql)
			return
		},
		err,
	)
}

// Warn implements logger.Interface.
func (l *migrationTrapLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	l.BaseLogger.Warn(ctx, fmt, args...)
}

func (l *migrationTrapLogger) considerQuery(sql string) {
	upSql := strings.ToUpper(sql)
	if strings.HasPrefix(upSql, "SELECT ") {
		return
	}

	l.PendingMigrations = append(l.PendingMigrations, sql)
}

type MigrationTrapMigrator struct {
	gorm.Migrator
	db *gorm.DB
}

func (m *MigrationTrapMigrator) AutoMigrate(dst ...interface{}) error {
	// Logger is bypassed to be able to trap migrations' SQL queries.
	trapLogger := &migrationTrapLogger{BaseLogger: m.db.Logger}
	m.db.Logger = trapLogger

	err := m.Migrator.AutoMigrate(dst...)
	if err != nil {
		return err
	}
	if len(trapLogger.PendingMigrations) > 0 {
		return &TrappedMigrationsError{PendingMigrations: trapLogger.PendingMigrations}
	}
	return nil
}
