package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

type CLI struct {
	URN string `help:"urn for the database" default:"file:test.db?cache=shared&mode=memory"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	cli := &CLI{}
	ctx := kong.Parse(cli)

	err := ctx.Run()
	if err != nil {
		log.Fatalf("could not execute: %s", err)
	}
}


func (c *CLI) Run() error {
	db, err := sql.Open("libsql", c.URN)
	if err != nil {
		return fmt.Errorf("could open database: %w", err)
	}

	migrationsFS, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("could not wrap migrations: %w", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not wrap driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance(
		"iofs", migrationsFS,
		"ql", driver,
	)
	if err != nil {
		return fmt.Errorf("could not setup migrations: %w", err)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not migrate up: %w", err)
	}

	return nil
}