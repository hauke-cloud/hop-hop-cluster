package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/hauke-cloud/hop-hop-cluster/internal/config"
	di "github.com/hauke-cloud/hop-hop-cluster/pkg/di"
)

func main() {
	v, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	cfg, err := config.ParseConfig(v)
	if err != nil {
		log.Fatal("cannot parse config: ", err)
	}

	config, err := config.LoadClientCertificates(cfg)
	if err != nil {
		log.Fatal("cannot load client certificates: ", err)
	}

	startApp, diErr := di.Initialize(config)
	if diErr != nil {
		log.Fatal("cannot load apps: ", diErr)
	}

	// Separat initialization step
	err = startApp.Initialize()
	if err != nil {
		log.Fatal("cannot initialize apps: ", err)
	}

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	startApp.Run(ctx)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	if err := startApp.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shut down app: %v", err)
	}
}
