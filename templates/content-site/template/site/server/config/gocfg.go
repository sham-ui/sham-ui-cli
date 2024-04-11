package config

import (
	"fmt"
	"os"
	"site/pkg/tracing"
	"strings"
	"time"

	"github.com/go-logr/logr"
)

const defaultGoCfgConfig = `
[server]
port = 3000
cors = true

[api]
address = /tmp/site/cms.sock
dialTimeout = 5s

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
		Cors bool
	}
	API struct {
		Address     string
		DialTimeout gocfgDuration
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
			Cors: g.Server.Cors,
		},
		API: API{
			Address:     g.API.Address,
			DialTimeout: time.Duration(g.API.DialTimeout),
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

type gocfgDuration time.Duration

func (d *gocfgDuration) UnmarshalText(text []byte) error {
	v, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = gocfgDuration(v)
	return nil
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
