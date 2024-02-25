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
	dialector, err := buildDialector()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}
	if err = migrate(db); err != nil {
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

func buildDialector() (gorm.Dialector, error) {
	dsn, err := determineDsn()
	if err != nil {
		return nil, fmt.Errorf("error establishing database connection parameters: %v", err)
	}

	baseDialector := postgres.Open(dsn)
	if utils.ReadEnvBoolean("DB_AUTO_MIGRATE", false) {
		return baseDialector, nil
	}

	return NewMigrationTrapDialector(baseDialector), nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		model.User{},
	)
	if _, ok := err.(*TrappedMigrationsError); ok {
		return err
	} else if err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	// dialector, ok := db.Dialector.(*MigrationTrapDialector)
	// if ok && len(dialector.MigrationChanges) > 0 {
	// 	ctx := context.Background()
	// 	db.Logger.Warn(ctx, "Pending changes detected by Gorm's AutoMigrate:")
	// 	for index, change := range dialector.MigrationChanges {
	// 		db.Logger.Warn(ctx, "%d: %s", index, change)
	// 	}
	// 	return fmt.Errorf("in %q, the database is expected to have been updated through migrations", os.Getenv("APP_ENV"))
	// }

	return nil
}
