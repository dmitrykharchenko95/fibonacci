package server

import (
	"log"
	"sync"
	"time"

	"github.com/dmitrykharchenko95/fibonacci/config"
	"github.com/dmitrykharchenko95/fibonacci/internal/rds"
	grpcserver "github.com/dmitrykharchenko95/fibonacci/internal/server/grpc"
	httpserver "github.com/dmitrykharchenko95/fibonacci/internal/server/http"
)

const (
	defaultTimeout = 10 * time.Second
)

type Sever struct {
	http *httpserver.Server
	grpc *grpcserver.Server
}

func New(cfg *config.Config) *Sever {
	httpTimeout, err := time.ParseDuration(cfg.HTTP.Timeout)
	if err != nil {
		log.Printf("parse grpc timeout fail: %v", err)
		log.Printf("use default value - %v", defaultTimeout)
		httpTimeout = defaultTimeout
	}

	grpcTimeout, err := time.ParseDuration(cfg.GRPC.Timeout)
	if err != nil {
		log.Printf("parse grpc timeout fail: %v", err)
		log.Printf("use default value - %v", defaultTimeout)
		grpcTimeout = defaultTimeout
	}

	rdb := rds.NewRedisClient(cfg.Redis)

	return &Sever{
		http: httpserver.New(cfg.HTTP.Host, cfg.HTTP.Port, httpTimeout, rdb),
		grpc: grpcserver.New(cfg.GRPC.Host, cfg.GRPC.Port, grpcTimeout, rdb),
	}
}

func (s *Sever) Start() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := s.http.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := s.grpc.Start()
		if err != nil {
			log.Fatal(err)
		}
	}(&wg)

	wg.Wait()
}

func (s *Sever) Stop() {
	s.grpc.Stop()

	err := s.http.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
