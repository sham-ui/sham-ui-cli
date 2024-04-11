package main

import (
	"os"
	"{{ shortName }}/config"
	"{{ shortName }}/internal/app"
	"{{ shortName }}/pkg/logger"
)

func main() {
	log := logger.NewLogger(128) //nolint:gomnd
	cfg, err := config.LoadConfiguration(log, "config.cfg")
	if err != nil {
		log.Error(err, "can't load config")
		os.Exit(1)
	}
	os.Exit(app.Run(log, cfg))
}
