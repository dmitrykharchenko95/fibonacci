package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/dmitrykharchenko95/fibonacci/config"
	"github.com/dmitrykharchenko95/fibonacci/internal/server"
	"github.com/dmitrykharchenko95/fibonacci/internal/service"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/fibonacci_config.json", "path to config file")
	flag.BoolVar(&service.UseRedis, "redis", true, "using Redis for caching")

}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if cfg.Redis.MaxErrors == 0{
		service.UseRedis = false
	}

	log.Printf("Fibonacci started. UseRaddis = %v", service.UseRedis)

	s := server.New(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	s.Start()
}
