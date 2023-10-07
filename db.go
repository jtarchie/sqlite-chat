package main

import (
	"errors"
	"fmt"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewDB(urn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("libsql", urn)
	if err != nil {
		return nil, fmt.Errorf("could open database: %w", err)
	}

	migrationsFS, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not wrap migrations: %w", err)
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not wrap driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance(
		"iofs", migrationsFS,
		"ql", driver,
	)
	if err != nil {
		return nil, fmt.Errorf("could not setup migrations: %w", err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil,fmt.Errorf("could not migrate up: %w", err)
	}

	return db, nil
}