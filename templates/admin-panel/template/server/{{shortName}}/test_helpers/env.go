package test_helpers

import (
	"github.com/go-logr/logr/testr"
    "github.com/urfave/negroni"
	"{{shortName}}/app"
	"{{shortName}}/test_helpers/client"
	testDB "{{shortName}}/test_helpers/database"
	"path"
	"testing"
)

type TestEnv struct {
	DB  *testDB.TestDatabase
	API *client.ApiClient
	T   *testing.T
}

func (env *TestEnv) Default() func() {
	env.DB.Clear()
	return func() {
		env.DB.DB.Close()
	}
}

func (env *TestEnv) CreateUser() {
	_, err := env.DB.DB.Exec("INSERT INTO public.members (id, name, email, password) VALUES (1, 'test', 'email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO')")
	if nil != err {
		env.T.Fatalf("can't create test user: %s", err)
	}
}

func (env *TestEnv) CreateSuperUser() {
	_, err := env.DB.DB.Exec("INSERT INTO public.members (id, name, email, password, is_superuser) VALUES (1, 'test', 'email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO', TRUE)")
	if nil != err {
		env.T.Fatalf("can't create test super user: %s", err)
	}
}

func NewTestEnv(t *testing.T) *TestEnv {
	env := &TestEnv{T: t}
	n := negroni.New()
	log := testr.New(t).V(1)
	db := app.StartApplication(log, path.Join("testdata", "config.cfg"), n)
	env.DB = testDB.NewTestDatabase(db)
	env.API = client.NewApiClient(n)
	return env
}
