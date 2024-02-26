package db

import (
	"fmt"
	"os"

	"github.com/adeynack/finances/model"
	"github.com/adeynack/finances/utils"
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
		if err := db.AutoMigrate(listGormModels()...); err != nil {
			return fmt.Errorf("error auto-migrating: %v", err)
		}
		return nil
	}

	// Temporary silence the Gorm logger (avoiding seeind the pending migrations twice).
	originalDbLogger := db.Logger
	db.Logger = db.Logger.LogMode(logger.Warn)
	defer func() { db.Logger = originalDbLogger }()

	// Performing AutoMigrate inside of a DryRun session to only collects the
	// migration's SQL (through the "trap logger") without performing it.
	txLogger := &migrationTrapLogger{}
	tx := db.Session(&gorm.Session{DryRun: true, Logger: txLogger})

	err := tx.AutoMigrate(listGormModels()...)
	if err != nil {
		return fmt.Errorf("error auto-migrating in dry run: %v", err)
	}
	if len(txLogger.PendingMigrations) > 0 {
		return &TrappedMigrationsError{PendingMigrations: txLogger.PendingMigrations}
	}

	return nil
}

func listGormModels() []any {
	return []any{
		model.User{},
	}
}
