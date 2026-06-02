package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	sqlmigrations "github.com/nenrinyear/small-library-system/api/migrations"
)

func RunMigrations(sqlDB *sql.DB) (err error) {
	sourceDriver, err := iofs.New(sqlmigrations.Files, ".")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}

	databaseDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("create mysql migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "mysql", databaseDriver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		sourceErr, databaseErr := migrator.Close()
		if err == nil {
			if sourceErr != nil {
				err = fmt.Errorf("close migration source: %w", sourceErr)
			} else if databaseErr != nil {
				err = fmt.Errorf("close migration database: %w", databaseErr)
			}
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
