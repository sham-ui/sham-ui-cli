package config

import (
	"fmt"
	"github.com/go-logr/logr"
	"gopkg.in/gcfg.v1"
	stdIOutls "io/ioutil"
	stdOS "os"
	"strings"
)

// Interface for "os" package (for mocking in tests)
type operationSystem interface {
	Stat(name string) (stdOS.FileInfo, error)
	IsNotExist(err error) bool
	Exit(code int)
}
type originalOS struct{}

func (originalOS) Stat(name string) (stdOS.FileInfo, error) { return stdOS.Stat(name) }
func (originalOS) IsNotExist(err error) bool                { return stdOS.IsNotExist(err) }
func (originalOS) Exit(code int)                            { stdOS.Exit(code) }

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
	Upload   upload
	Api      api
)

type Config struct {
	Server   server
	Database dataBaseConfig
	Session  session
	Upload   upload
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

func (dbCfg *dataBaseConfig) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name)
}

type session struct {
	Secret string
}

type api struct {
	SocketPath string
}

type upload struct {
	Path string
}

const defaultConfig = `
[server]
port = 3003
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

[upload]
path = ./upload

[api]
socketPath = /tmp/{{ name }}/cms.sock
`

func LoadConfiguration(logger logr.Logger, configFilename string) {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		err := ioutil.WriteFile(configFilename, []byte(strings.TrimSpace(defaultConfig)), 0644)
		if nil != err {
			logger.Error(err, "Fail write config", "filename", configFilename)
			os.Exit(1)
		}
		logger.Info("Config file created", "filename", configFilename)
	}

	var cfg Config
	err := gcfg.ReadFileInto(&cfg, configFilename)
	if nil != err {
		logger.Error(err, "Fail read config", "filename", configFilename)
		os.Exit(1)
	}
	Server = cfg.Server
	DataBase = cfg.Database
	Session = cfg.Session
	Upload = cfg.Upload
	Api = cfg.Api
}
