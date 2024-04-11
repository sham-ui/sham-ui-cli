package config

import (
	"cms/pkg/postgres"
	"cms/pkg/tracing"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
)

const defaultGoCfgConfig = `
[server]
port = 3002
csrf = 32-byte-long-auth-key
cors = true

[api]
address = /tmp/{{ name }}/cms.sock

[database]
host = 127.0.0.1
port = 5432
name = {{ dbName }}
user = {{ dbUser }}
pass = {{ dbPassword}}

[session]
domain = localhost
secret = secret-key

[upload]
path = ./upload

[tracer]
; Tracer disabled by default, for enable it uncomment next line & update authorization key
;endpoint = localhost:5080
path = /api/default/v1/traces
authorization = cm9vdEBleGFtcGxlLmNvbTpETktERTFKNkJTSE9DTlVa
serviceName = site
version = 0.0.1
environment = prod
`
const defaultConfigFilePermission = 0o600

// gocfg is special type for parsing .cfg files.
type gocfg struct {
	Server struct {
		Port int
		CSRF string
		Cors bool
	}
	API struct {
		Address string
	}
	Database struct {
		Host string
		Port int
		Name string
		User string
		Pass string
	}
	Session struct {
		Domain string
		Secret string
	}
	Upload struct {
		Path string
	}
	Tracer struct {
		Endpoint      string
		Path          string
		Authorization string
		ServiceName   string
		Version       string
		Environment   string
	}
}

func (g gocfg) toConfig() *Config {
	return &Config{
		Server: Server{
			Port: g.Server.Port,
			CSRF: []byte(g.Server.CSRF),
			CORS: g.Server.Cors,
		},
		API: API{
			Address: g.API.Address,
		},
		Database: postgres.Config{
			Host: g.Database.Host,
			Port: g.Database.Port,
			Name: g.Database.Name,
			User: g.Database.User,
			Pass: g.Database.Pass,
		},
		Session: Session{
			Domain: g.Session.Domain,
			Secret: g.Session.Secret,
		},
		Upload: Upload{
			Path: g.Upload.Path,
		},
		Tracer: tracing.Config{
			Endpoint:      g.Tracer.Endpoint,
			Path:          g.Tracer.Path,
			Authorization: g.Tracer.Authorization,
			ServiceName:   g.Tracer.ServiceName,
			Version:       g.Tracer.Version,
			Environment:   g.Tracer.Environment,
		},
	}
}

func createDefaultConfig(logger logr.Logger, configPath string) error {
	_, err := os.Stat(configPath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("check config file stats: %w", err)
	}
	content := []byte(strings.TrimSpace(defaultGoCfgConfig))
	if err := os.WriteFile(configPath, content, os.FileMode(defaultConfigFilePermission)); err != nil {
		return fmt.Errorf("create default config: %w", err)
	}
	logger.Info("config created", "configPath", configPath)
	return nil
}
