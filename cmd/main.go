package main

import (
	"context"
	"log"
	"nilchan-hackaton/internal/app"
	"nilchan-hackaton/internal/config"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
