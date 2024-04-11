package postgres

import (
	"cms/pkg/asserts"
	"testing"
)

func TestDatabase_URL(t *testing.T) {
	d := Config{
		Host: "localhost",
		Port: 5432,
		Name: "db",
		User: "user",
		Pass: "pass",
	}
	asserts.Equals(t, "postgres://user:pass@localhost:5432/db", d.URL())
}
