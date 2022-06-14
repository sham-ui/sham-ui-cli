package test_helpers

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"cms/app"
	"cms/test_helpers/client"
	testDB "cms/test_helpers/database"
	"path"
)

type TestEnv struct {
	DB  *testDB.TestDatabase
	API *client.ApiClient
}

func (env *TestEnv) DisableLogger() {
	log.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.FatalLevel)
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
		log.Fatalf("can't create test user: %s", err)
	}
}

func (env *TestEnv) CreateSuperUser() {
	_, err := env.DB.DB.Exec("INSERT INTO public.members (id, name, email, password, is_superuser) VALUES (1, 'test', 'email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO', TRUE)")
	if nil != err {
		log.Fatalf("can't create test super user: %s", err)
	}
}

func NewTestEnv() *TestEnv {
	env := &TestEnv{}
	env.DisableLogger()
	n := negroni.New()
	db := app.StartApplication(path.Join("testdata", "config.cfg"), n)
	env.DB = testDB.NewTestDatabase(db)
	env.API = client.NewApiClient(n)
	return env
}
