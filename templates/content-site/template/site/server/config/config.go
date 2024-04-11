package config

import (
	"fmt"
	"site/pkg/tracing"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"gopkg.in/gcfg.v1"
)

type (
	Config struct {
		Server Server
		API    API
		Tracer tracing.Config
	}
	Server struct {
		Port int
		Cors bool
	}
	API struct {
		Address     string
		DialTimeout time.Duration
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
