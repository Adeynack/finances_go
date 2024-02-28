package database

import (
	"fmt"
	"os"

	"github.com/adeynack/finances/app/utils"
	"github.com/adeynack/finances/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	dsn, err := determineDsn()
	if err != nil {
		return nil, fmt.Errorf("error establishing database connection parameters: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err = attemptAutoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func errExpectedEnv(expectedEnv string) error {
	return fmt.Errorf(
		"in application environment %q, expected environment variable %q to be defined",
		os.Getenv("APP_ENV"),
		expectedEnv,
	)
}

func determineDsn() (string, error) {
	// When DATABASE_URL is set, use it as-is.
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl != "" {
		return databaseUrl, nil
	}

	// Otherwise, build the connection string from individual ENV values.
	pgHost := os.Getenv("PGHOST")
	if pgHost == "" {
		return "", errExpectedEnv("PGHOST")
	}
	pgUser := os.Getenv("PGUSER")
	if pgUser == "" {
		return "", errExpectedEnv("PGUSER")
	}
	pgDatabaseName := os.Getenv("PGDB")
	if pgDatabaseName == "" {
		return "", errExpectedEnv("PGDB")
	}
	pgPort := os.Getenv("PGPORT")
	if pgPort == "" {
		return "", errExpectedEnv("PGPORT")
	}
	pgSSLMode := os.Getenv("PGSSLMODE")
	if pgSSLMode == "" {
		pgSSLMode = "prefer" // default (https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNECT-SSLMODE)
	}
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s", pgHost, pgUser, pgDatabaseName, pgPort, pgSSLMode)
	return dsn, nil
}

func attemptAutoMigrate(db *gorm.DB) error {
	allowAutoMigrate := utils.ReadEnvBoolean("DB_AUTO_MIGRATE", false)
	if allowAutoMigrate {
		if err := db.AutoMigrate(model.All()...); err != nil {
			return fmt.Errorf("error auto-migrating: %v", err)
		}
		return nil
	}

	// Temporary switch the Gorm logger (avoiding seeind the pending migrations twice).
	originalDbLogger := db.Logger
	txLogger := &migrationTrapLogger{Interface: originalDbLogger}
	db.Logger = txLogger
	defer func() { db.Logger = originalDbLogger }()

	// Performing AutoMigrate inside of an automatically rolled-back transaction.
	// Sadly, Gorm's `DryRun` session does not work as expected (crashes sometimes with `dryrun not supported` or segfaults.
	tx := db.Begin()
	defer tx.Rollback()

	err := tx.AutoMigrate(model.All()...)
	if err != nil {
		return fmt.Errorf("error auto-migrating in dry run: %v", err)
	}
	if len(txLogger.PendingMigrations) > 0 {
		return &TrappedMigrationsError{PendingMigrations: txLogger.PendingMigrations}
	}

	return nil
}
