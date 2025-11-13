package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// MigrateFS runs migrations using an embedded filesystem
func (s *service) MigrateFS(migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return s.Migrate(dir)
}

// Migrate runs database migrations from a directory
func (s *service) Migrate(dir string) error {
	db, err := sql.Open("pgx", s.connStr)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	defer db.Close()
	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func (s *service) MigrateDownFS(migrationsFS fs.FS, dir string) error {
	db, err := sql.Open("pgx", s.connStr)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}

	err = goose.Down(db, dir)
	if err != nil {
		return fmt.Errorf("goose down: %w", err)
	}
	return nil
}

// MigrateDownAllFS rolls back all migrations using embedded filesystem
func (s *service) MigrateDownAllFS(migrationsFS fs.FS, dir string) error {
	db, err := sql.Open("pgx", s.connStr)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate down all: %w", err)
	}

	err = goose.DownTo(db, dir, 0)
	if err != nil {
		return fmt.Errorf("goose down to 0: %w", err)
	}
	return nil
}
