package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"nilchan-hackaton/internal/app"
	"nilchan-hackaton/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env.local")

	cfg, err := config.Load("internal/config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

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
