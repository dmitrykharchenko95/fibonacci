package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	HTTP  HTTPConfig
	GRPC  GRPCConfig
	Redis RedisConfig
}

type HTTPConfig struct {
	Host    string        `config:"http_host"`
	Port    string        `config:"http_port"`
	Timeout string `config:"http_timeout"`
}

type GRPCConfig struct {
	Host    string        `config:"grpc_host"`
	Port    string        `config:"grpc_port"`
	Timeout string `config:"grpc_timeout"`
}

type RedisConfig struct {
	Host       string        `config:"redis_host"`
	Port       string        `config:"redis_port"`
	Expiration string `config:"redis_expiration"`
	MaxErrors  int           `config:"redis_max_errors"`
}

func New(configFile string) (*Config, error) {
	cfg := &Config{}

	l := confita.NewLoader(file.NewBackend(configFile))

	if err := l.Load(context.Background(), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
