package config

import (
	"errors"
	"fmt"
	"net"

	"github.com/caarlos0/env/v10"
)

type (
	Config struct {
		Settings struct {
			StorageType string `env:"STORAGE_TYPE" envDefault:"postgres"`
		}

		GraphQL struct {
			Port string `env:"GRAPHQL_PORT" envDefault:"8080"`
		}

		PG struct {
			Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
			Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
			DB       string `env:"POSTGRES_DB" envDefault:"posts"`
			User     string `env:"POSTGRES_USER" envDefault:"posts_user"`
			Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
		}
	}
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}

func (c *Config) ConstructPostgresURL() string {
	hostPort := net.JoinHostPort(c.PG.Host, c.PG.Port)
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PG.User,
		c.PG.Password,
		hostPort,
		c.PG.DB,
	)
}

func (c *Config) GetStorageType() (StorageType, error) {
	switch c.Settings.StorageType {
	case "postgres":
		return Postgres, nil
	case "inmemory":
		return InMemory, nil
	default:
		return Unknown, errors.New("unknown storage type")
	}
}

type StorageType int

const (
	Postgres StorageType = iota
	InMemory
	Unknown
)
