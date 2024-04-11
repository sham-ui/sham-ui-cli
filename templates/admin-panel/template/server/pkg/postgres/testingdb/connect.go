package testingdb

import (
	"context"
	"{{ shortName }}/config"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/postgres"
	"path"
	"runtime"
	"testing"

	"github.com/go-logr/logr/testr"
)

func Connect(t *testing.T) *postgres.Database {
	t.Helper()

	// Get file path to config
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled
	filepath := path.Join(path.Dir(filename), "../../../testdata/config.cfg")

	log := testr.New(t)
	cfg, err := config.LoadConfiguration(log, filepath)
	asserts.NoError(t, err)
	db, err := postgres.New(log, cfg.Database.URL())
	asserts.NoError(t, err)
	t.Cleanup(func() {
		asserts.NoError(t, db.GracefulShutdown(context.Background()))
	})
	return db
}
