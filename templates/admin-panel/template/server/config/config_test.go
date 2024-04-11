package config

import (
	"os"
	"{{ shortName }}/pkg/asserts"
	"{{ shortName }}/pkg/postgres"
	"{{ shortName }}/pkg/tracing"
	"path"
	"strings"
	"testing"

	"github.com/go-logr/logr/testr"
)

func TestCreateConfigIfNotExists(t *testing.T) {
	configPath := path.Join(t.TempDir(), "config.cfg")

	cfg, err := LoadConfiguration(testr.New(t).V(1), configPath)
	asserts.NoError(t, err)
	asserts.Equals(t, &Config{
		Server: Server{
			Port: 3002,
			CSRF: []byte("32-byte-long-auth-key"),
			CORS: true,
		},
		Database: postgres.Config{
			Host: "127.0.0.1",
			Port: 5432,
			Name: "pa",
			User: "pauser",
			Pass: "pauser",
		},
		Session: Session{
			Domain: "localhost",
			Secret: "secret-key",
		},
		Tracer: tracing.Config{
			Endpoint:      "",
			Path:          "/api/default/v1/traces",
			Authorization: "cm9vdEBleGFtcGxlLmNvbTpETktERTFKNkJTSE9DTlVa",
			ServiceName:   "site",
			Version:       "0.0.1",
			Environment:   "prod",
		},
	}, cfg, "config")
	content, err := os.ReadFile(configPath)
	asserts.Equals(t, strings.TrimSpace(defaultGoCfgConfig), string(content), "file content")
	asserts.NoError(t, err)
}

func TestReadConfig(t *testing.T) {
	configPath := path.Join("testdata", "config.cfg")
	cfg, err := LoadConfiguration(testr.New(t).V(1), configPath)
	asserts.NoError(t, err)
	asserts.Equals(t, &Config{
		Server: Server{
			Port: 3001,
			CSRF: []byte("32-byte-long-auth-key"),
			CORS: true,
		},
		Database: postgres.Config{
			Host: "127.0.0.1",
			Port: 5432,
			Name: "pa",
			User: "pauser",
			Pass: "pauser",
		},
		Session: Session{
			Domain: "localhost",
			Secret: "secret-key",
		},
		Tracer: tracing.Config{
			Endpoint:      "localhost:5080",
			Path:          "/api/default/v1/traces",
			Authorization: "cm9vdEBleGFtcGxlLmNvbTpETktERTFKNkJTSE9DTlVa",
			ServiceName:   "site",
			Version:       "0.0.1",
			Environment:   "dev",
		},
	}, cfg)
}

func TestServer(t *testing.T) {
	s := Server{
		Port: 1234,
		CSRF: nil,
		CORS: false,
	}
	asserts.Equals(t, ":1234", s.Address())
	asserts.Equals(t, "http://localhost:1234", s.URL())
}
