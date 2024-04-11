package postgres

import "fmt"

type Config struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func (d Config) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, d.Pass, d.Host, d.Port, d.Name) //nolint:nosprintfhostport
}
