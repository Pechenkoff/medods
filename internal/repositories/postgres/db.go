package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// RunMigration - do migration for PostgreSQL
func MustRunMigration(dbURL, migrationPath string) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		panic(fmt.Errorf("error connecting to DB: %v", err))
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("error creating driver for migration: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver,
	)

	if err != nil {
		panic(fmt.Errorf("error creating migrations: %v", err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Errorf("failed to run migrations^ %v", err))
	}
}
