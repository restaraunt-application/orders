package config

import (
	"fmt"
)

type PostgreSQLConfig struct {
	Host               string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port               int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User               string `env:"POSTGRES_USER" envDefault:"orders"`
	Password           string `env:"POSTGRES_PASSWORD" envDefault:"orders"`
	DBName             string `env:"POSTGRES_DB_NAME" envDefault:"orders"`
	MaxOpenConnections int    `env:"POSTGRES_MAX_OPEN_CONNECTIONS" envDefault:"10"`
}

func (c *PostgreSQLConfig) ConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}
