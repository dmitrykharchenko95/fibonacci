package rds

import (
	"log"
	"net"
	"time"

	"github.com/dmitrykharchenko95/fibonacci/config"
	"github.com/go-redis/redis/v8"
)

const defaultExpiration = 12 * time.Hour

type Client struct {
	Cl         *redis.Client
	Expiration time.Duration
	MaxErrors  int
}

func NewRedisClient(cfg config.RedisConfig) *Client {
	exp, err := time.ParseDuration(cfg.Expiration)
	if err != nil {
		log.Printf("parse redis Expiration fail: %v", err)
		log.Printf("use default value - %v", defaultExpiration)
		exp = defaultExpiration
	}

	return &Client{
		Cl: redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
			Password: "",
			DB:       0,
		}),
		Expiration: exp,
		MaxErrors:  cfg.MaxErrors,
	}
}
