package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/gcfg.v1"
	stdIOutls "io/ioutil"
	stdOS "os"
	"strings"
)

// Interface for "os" package (for mocking in tests)
type operationSystem interface {
	Stat(name string) (stdOS.FileInfo, error)
	IsNotExist(err error) bool
}
type originalOS struct{}

func (originalOS) Stat(name string) (stdOS.FileInfo, error) { return stdOS.Stat(name) }
func (originalOS) IsNotExist(err error) bool                { return stdOS.IsNotExist(err) }

// Interface for "io/ioutil" package
type osIOutil interface {
	WriteFile(string, []byte, stdOS.FileMode) error
}
type originalIOutil struct{}

func (originalIOutil) WriteFile(filename string, data []byte, perm stdOS.FileMode) error {
	return stdIOutls.WriteFile(filename, data, perm)
}

var os operationSystem = originalOS{}
var ioutil osIOutil = originalIOutil{}

var (
	Server server
	Api    api
)

type Config struct {
	Server server
	Api    api
}

type server struct {
	Port           int
	AllowedDomains []string
}

type api struct {
	SocketPath string
}

const defaultConfig = `
[server]
port = 3002
allowedDomains = http://127.0.0.1:3000
allowedDomains = http://localhost:3000

[api]
socketPath = /tmp/{{ name }}/cms.sock'
`

func LoadConfiguration(configFilename string) {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		err := ioutil.WriteFile(configFilename, []byte(strings.TrimSpace(defaultConfig)), 0644)
		if nil != err {
			logrus.WithError(err).Fatal("Fail write config", logrus.Fields{"configFilename": configFilename})
		} else {
			logrus.Info("Create config file")
		}
	}

	var cfg Config
	err := gcfg.ReadFileInto(&cfg, configFilename)
	if nil != err {
		logrus.WithError(err).Fatal("Fail read config", logrus.Fields{"configFilename": configFilename})
	}
	Server = cfg.Server
	Api = cfg.Api
}
