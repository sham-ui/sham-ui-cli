package migrations

import (
	"database/sql"
	"fmt"
)

type Migration interface {
	Name() string
	Apply() error
}

type DBMigration struct {
	db       *sql.DB
	name     string
	sqlQuery string
}

func (m *DBMigration) Name() string {
	return m.name
}

func (m *DBMigration) Apply() error {
	_, err := m.db.Exec(m.sqlQuery)
	if nil != err {
		return fmt.Errorf("exec SQL: %s", err)
	}
	return nil
}

func NewDBMigration(db *sql.DB, name, sqlQuery string) Migration {
	return &DBMigration{
		db:       db,
		name:     name,
		sqlQuery: sqlQuery,
	}
}
