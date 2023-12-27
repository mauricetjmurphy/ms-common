package db

import (
	"fmt"
	"strings"
)

const (
	DefaultHost = "127.0.0.1"
	DefaultPort = 3306
)

// Config includes required and optional parameters for OpenConnection.
type Config struct {
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	AWSRegion          string
	AWSSecretID        string
	MigrationSourceURL string
}

// NewConfigs creates Config on given configuration options.
func NewConfigs(opts ...Option) *Config {
	cf := &Config{
		Host: DefaultHost,
		Port: DefaultPort,
	}
	for _, opt := range opts {
		opt(cf)
	}
	return cf
}

func (c *Config) IsLocal() bool {
	return c != nil && c.AWSSecretID == ""
}

func (c *Config) RequiredMigration() bool {
	return c != nil && len(strings.TrimSpace(c.MigrationSourceURL)) > 0
}

func (c *Config) toDSN() *dsnConf {
	return &dsnConf{
		Host:     c.Host,
		Port:     c.Port,
		Name:     c.Name,
		User:     c.User,
		Password: c.Password,
	}
}

type dsnConf struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

func (dsn *dsnConf) toDSNConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true", dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.Name)
}
