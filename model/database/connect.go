package database

import (
	"fmt"
	"os"

	"github.com/adeynack/finances/app/appenv"
	"github.com/adeynack/finances/app/utils"
	"github.com/adeynack/finances/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	EnvDatabaseUrl     = "DATABASE_URL"
	EnvDatabaseHost    = "PGHOST"
	EnvDatabaseUser    = "PGUSER"
	EnvDatabaseName    = "PGDATABASE"
	EnvDatabasePort    = "PGPORT"
	EnvDatabaseSSLMode = "PGSSLMODE"
	EnvDbAutoMigrate   = "DB_AUTO_MIGRATE"
)

func init() {
	appenv.Init()
}

func Connect() (*gorm.DB, error) {
	dsn, err := determineDsn()
	if err != nil {
		return nil, fmt.Errorf("error establishing database connection parameters: %w", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err = attemptAutoMigrate(db); err != nil {
		return nil, err
	}

	err = db.Callback().Create().Before("gorm:before_save").Register("validate", validateCallback)
	if err != nil {
		return nil, fmt.Errorf("error registering \"validate\" callback")
	}

	return db, nil
}

func errExpectedEnv(expectedEnv string) error {
	return fmt.Errorf(
		"in application environment %q, expected environment variable %q to be defined",
		appenv.Env,
		expectedEnv,
	)
}

func determineDsn() (string, error) {
	// When DATABASE_URL is set, use it as-is.
	databaseUrl := os.Getenv(EnvDatabaseUrl)
	if databaseUrl != "" {
		return databaseUrl, nil
	}

	// Otherwise, build the connection string from individual ENV values.
	pgHost := os.Getenv(EnvDatabaseHost)
	if pgHost == "" {
		return "", errExpectedEnv(EnvDatabaseHost)
	}
	pgUser := os.Getenv(EnvDatabaseUser)
	if pgUser == "" {
		return "", errExpectedEnv(EnvDatabaseUser)
	}
	pgDatabaseName := os.Getenv(EnvDatabaseName)
	if pgDatabaseName == "" {
		return "", errExpectedEnv(EnvDatabaseName)
	}
	pgPort := os.Getenv(EnvDatabasePort)
	if pgPort == "" {
		return "", errExpectedEnv(EnvDatabasePort)
	}
	pgSSLMode := os.Getenv(EnvDatabaseSSLMode)
	if pgSSLMode == "" {
		pgSSLMode = "prefer" // default (https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNECT-SSLMODE)
	}
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s", pgHost, pgUser, pgDatabaseName, pgPort, pgSSLMode)
	return dsn, nil
}

func attemptAutoMigrate(db *gorm.DB) error {
	allowAutoMigrate := utils.ReadEnvBoolean(EnvDbAutoMigrate, false)
	if allowAutoMigrate {
		if err := db.AutoMigrate(model.All()...); err != nil {
			return fmt.Errorf("error auto-migrating: %w", err)
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
		return fmt.Errorf("error auto-migrating in dry run: %w", err)
	}
	if len(txLogger.PendingMigrations) > 0 {
		return &TrappedMigrationsError{PendingMigrations: txLogger.PendingMigrations}
	}

	return nil
}
