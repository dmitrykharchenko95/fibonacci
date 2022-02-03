package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {

	want := &Config{
		HTTP: HTTPConfig{
			Host:    "0.0.0.0",
			Port:    "8080",
			Timeout: "10s",
		},
		GRPC: GRPCConfig{
			Host:    "0.0.0.0",
			Port:    "50052",
			Timeout: "10s",
		},
		Redis: RedisConfig{
			Host:       "0.0.0.0",
			Port:       "6379",
			Expiration: "12h",
			MaxErrors:  6,
		},
	}

	t.Run("base", func(t *testing.T) {
		got, err := New("../configs/fibonacci_config.json")
		require.NoError(t, err)
		require.Equal(t, got, want)
	})
}
