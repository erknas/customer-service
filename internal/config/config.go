package config

import (
	"fmt"
	"net/url"
	"os"
)

type (
	Config struct {
		GRPC struct {
			Port        string `env:"GRPC_PORT"`
			GatewayPort string `env:"GRPC_GATEWAY_PORT"`
		}

		Postgres struct {
			URL      string
			User     string `env:"POSTGRES_USER"`
			Password string `env:"POSTGRES_PASSWORD"`
			DB       string `env:"POSTGRES_DB"`
			Host     string `env:"POSTGRES_HOST"`
			Port     string `env:"POSTGRES_PORT"`
			MaxConns string `env:"MAX_CONNS"`
		}
	}
)

func Load() *Config {
	cfg := new(Config)

	cfg.GRPC.Port = os.Getenv("GRPC_PORT")
	cfg.GRPC.GatewayPort = os.Getenv("GRPC_GATEWAY_PORT")

	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.DB = os.Getenv("POSTGRES_DB")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.MaxConns = os.Getenv("MAX_CONNS")

	cfg.Postgres.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%s",
		url.QueryEscape(cfg.Postgres.User),
		url.QueryEscape(cfg.Postgres.Password),
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
		cfg.Postgres.MaxConns,
	)

	return cfg
}
