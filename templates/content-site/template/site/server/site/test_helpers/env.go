package test_helpers

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"path"
	"site/app"
	"site/config"
	"site/ssr"
	"site/test_helpers/client"
	"site/test_helpers/cms"
)

type TestEnv struct {
	API *client.ApiClient
	CMS *cms.MockCMSClient
}

func (env *TestEnv) DisableLogger() {
	log.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.FatalLevel)
}

func (env *TestEnv) Default() func() {
	return func() {}
}

func NewTestEnv(render ssr.Render) *TestEnv {
	env := &TestEnv{}
	env.DisableLogger()
	n := negroni.New()
	config.LoadConfiguration(path.Join("testdata", "config.cfg"))
	env.CMS = &cms.MockCMSClient{}
	app.StartApplication(n, env.CMS, render)
	env.API = client.NewApiClient(n)
	return env
}
