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
	dialector := postgres.Open(dsn)
	if err = attemptAutoMigrate(dialector); err != nil {
		return nil, err
	}

	// Create application's Gorm DB.
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
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

func attemptAutoMigrate(dialector gorm.Dialector) error {
	allowAutoMigrate := utils.ReadEnvBoolean("DB_AUTO_MIGRATE", false)

	migrateDialector := dialector
	if !allowAutoMigrate {
		migrateDialector = NewMigrationTrapDialector(dialector)
	}

	db, err := gorm.Open(migrateDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // ensure info logging is active so the "migration trap" to work
		DryRun: true,
	})
	if err != nil {
		return fmt.Errorf("error opening database connection for auto-migrate: %v", err)
	}

	err = db.AutoMigrate(
		model.User{},
	)

	if _, ok := err.(*TrappedMigrationsError); ok {
		return err
	} else if err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	return nil
}
