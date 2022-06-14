package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // github.com/lib/pq
)

// ConnToDB connects the database.
func ConnToDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open: %s", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %s", err)
	}
	return db, nil
}
