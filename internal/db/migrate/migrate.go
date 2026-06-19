package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

const migrationsDir = "migrations"

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("migrate: open database: %w", err)
	}
	return db, nil
}

func setup() error {
	goose.SetBaseFS(migrations)
	return goose.SetDialect("postgres")
}

func Up(dsn string) error {
	db, err := openDB(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := setup(); err != nil {
		return fmt.Errorf("migrate: setup: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("migrate: up: %w", err)
	}

	return nil
}

func Down(dsn string) error {
	db, err := openDB(dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := setup(); err != nil {
		return fmt.Errorf("migrate: setup: %w", err)
	}

	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("migrate: down: %w", err)
	}

	return nil
}
