package config

import (
	"fmt"
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
	Server   server
	DataBase dataBaseConfig
	Session  session
	Api      api
)

type Config struct {
	Server   server
	Database dataBaseConfig
	Session  session
	Api      api
}

type server struct {
	Port           int
	AllowedDomains []string
}

type dataBaseConfig struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

type session struct {
	Secret string
}

type api struct {
	SocketPath string
}

func (dbCfg *dataBaseConfig) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name)
}

const defaultConfig = `
[server]
port = 3000
allowedDomains = http://{{ name }}.com
allowedDomains = http://www.{{ name }}.com

[database]
host = 127.0.0.1
port = 5432
name = {{ dbName }}
user = {{ dbUser }}
pass = {{ dbPassword }}

[session]
secret = secret-key

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
	DataBase = cfg.Database
	Session = cfg.Session
	Api = cfg.Api
}
