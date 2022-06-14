package database

import (
	"database/sql"
	"testing"
)

type SQLQueryExecutor func(query string, args ...interface{})

type TestDatabase struct {
	DB *sql.DB
}

func (db *TestDatabase) ExecForCase(t *testing.T, name string) SQLQueryExecutor {
	return func(query string, args ...interface{}) {
		_, err := db.DB.Exec(query, args...)
		if nil != err {
			t.Fatalf("%s: exec SQL \"%s\": %s", name, query, err)
		}
	}
}

func (db *TestDatabase) Clear() {
	db.DB.Exec("DELETE FROM http_sessions")
	db.DB.Exec("DELETE FROM members")
}

func NewTestDatabase(db *sql.DB) *TestDatabase {
	return &TestDatabase{DB: db}
}
