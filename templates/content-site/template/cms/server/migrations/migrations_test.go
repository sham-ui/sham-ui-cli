package migrations

import (
	"cms/pkg/asserts"
	"cms/pkg/postgres/testingdb"
	"testing"

	"github.com/go-logr/logr/testr"
)

func TestUp(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)

	// Action
	err := Up(testr.New(t).V(0), db)

	// Assert
	asserts.NoError(t, err)
}

func TestDownTo(t *testing.T) {
	// Arrange
	db := testingdb.Connect(t)
	log := testr.New(t).V(0)
	t.Cleanup(func() {
		asserts.NoError(t, Up(log, db))
	})

	// Action
	err := DownTo(log, db, 1)

	// Assert
	asserts.NoError(t, err)
}
