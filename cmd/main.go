package main

import (
	"log"
	"nilchan-hackaton/internal/app"
	"nilchan-hackaton/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
