package test_helpers

import (
	"github.com/go-logr/logr/testr"
	"github.com/urfave/negroni"
	"path"
	"site/app"
	"site/config"
	"site/ssr"
	"site/test_helpers/client"
	"site/test_helpers/cms"
	"testing"
)

type TestEnv struct {
	API *client.ApiClient
	CMS *cms.MockCMSClient
}

func (env *TestEnv) Default() func() {
	return func() {}
}

func NewTestEnv(render ssr.Render, t *testing.T) *TestEnv {
	env := &TestEnv{}
	n := negroni.New()
	logger := testr.New(t).V(1)
	config.LoadConfiguration(logger, path.Join("testdata", "config.cfg"))
	env.CMS = &cms.MockCMSClient{}
	app.StartApplication(logger, n, env.CMS, render)
	env.API = client.NewApiClient(n)
	return env
}
