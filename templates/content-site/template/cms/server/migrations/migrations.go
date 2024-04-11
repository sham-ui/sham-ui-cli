package migrations

import (
	"cms/pkg/goose_logr"
	"cms/pkg/postgres"
	"embed"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func setupGoose(log logr.Logger) error {
	goose.SetLogger(goose_logr.NewLogger(log))
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	return nil
}

func Up(log logr.Logger, db *postgres.Database) error {
	if err := setupGoose(log); err != nil {
		return fmt.Errorf("setup goose: %w", err)
	}
	if err := goose.Up(db.Postgres(), "."); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}
	return nil
}

func DownTo(log logr.Logger, db *postgres.Database, version int64) error {
	if err := setupGoose(log); err != nil {
		return fmt.Errorf("setup goose: %w", err)
	}
	if err := goose.DownTo(db.Postgres(), ".", version); err != nil {
		return fmt.Errorf("down to migrations: %w", err)
	}
	return nil
}
