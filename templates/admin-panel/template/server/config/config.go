package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
	"{{ shortName }}/pkg/postgres"
	"{{ shortName }}/pkg/tracing"
	"strconv"

	"github.com/go-logr/logr"
)

type (
	Config struct {
		Server   Server
		Database postgres.Config
		Session  Session
		Tracer   tracing.Config
	}
	Server struct {
		Port int
		CSRF []byte
		CORS bool
	}
	Session struct {
		Secret string
		Domain string
	}
)

func (s Server) Address() string {
	return ":" + strconv.Itoa(s.Port)
}

func (s Server) URL() string {
	return "http://localhost:" + strconv.Itoa(s.Port)
}

func LoadConfiguration(logger logr.Logger, configPath string) (*Config, error) {
	if err := createDefaultConfig(logger, configPath); err != nil {
		return nil, err
	}
	var cfg gocfg
	if err := gcfg.ReadFileInto(&cfg, configPath); err != nil {
		return nil, fmt.Errorf("can't read config: %w", err)
	}
	return cfg.toConfig(), nil
}
