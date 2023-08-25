package migrations

import (
	"database/sql"
	"fmt"
	"github.com/go-logr/logr"
)

type Migrator struct {
	db            *sql.DB
	executedNames map[string]struct{}
	logger        logr.Logger
}

func (m *Migrator) Apply(migrations ...Migration) error {
	for _, migration := range migrations {
		if _, ok := m.executedNames[migration.Name()]; !ok {
			err := m.applyMigration(migration)
			if nil != err {
				return fmt.Errorf("apply migration %s fail: %s", migration.Name(), err)
			}
			m.logger.Info("apply migration", "name", migration.Name())
		}
	}
	return nil
}

func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations(
			"id" SERIAL UNIQUE PRIMARY KEY,
			"name" varchar(250),
			"applied_at" timestamp NOT NULL DEFAULT now()
		);
		ALTER TABLE ONLY migrations DROP CONSTRAINT IF EXISTS migrations_name_unique;
		ALTER TABLE ONLY migrations ADD CONSTRAINT migrations_name_unique UNIQUE (name);
	`)
	if nil != err {
		return fmt.Errorf("create migration table: %s", err)
	}

	return nil
}

func (m *Migrator) loadAppliedMigrationsFromDb() error {
	rows, err := m.db.Query("SELECT name FROM migrations")
	defer rows.Close()
	if nil != err {
		return fmt.Errorf("select: %s", err)
	}
	executedNames := make(map[string]struct{})
	for rows.Next() {
		name := ""
		err := rows.Scan(&name)
		if nil != err {
			return fmt.Errorf("scan row: %s", err)
		}
		executedNames[name] = struct{}{}
	}
	m.executedNames = executedNames
	return nil
}

func (m *Migrator) applyMigration(migration Migration) error {
	err := migration.Apply()
	if nil != err {
		return fmt.Errorf("apply: %s", err)
	}
	_, err = m.db.Query("INSERT INTO migrations(name) VALUES ($1)", migration.Name())
	if nil != err {
		return fmt.Errorf("insert into migrations: %s", err)
	}
	m.executedNames[migration.Name()] = struct{}{}
	return nil
}

func NewMigrator(logger logr.Logger, db *sql.DB) (*Migrator, error) {
	m := &Migrator{
		db:     db,
		logger: logger,
	}
	err := m.createMigrationsTable()
	if nil != err {
		return nil, fmt.Errorf("create migrations table: %s", err)
	}
	err = m.loadAppliedMigrationsFromDb()
	if nil != err {
		return nil, fmt.Errorf("load applied migrations from db: %s", err)
	}
	return m, nil
}
