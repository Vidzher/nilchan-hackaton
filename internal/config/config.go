package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server" env-required:"true"`
	Firecrawl   Firecrawl  `yaml:"firecrawl"`
}

type Firecrawl struct {
	APIKey  string        `yaml:"api_key" env:"FIRECRAWL_KEY"`
	BaseURL string        `yaml:"base_url" env-default:"https://api.firecrawl.dev/v2/scrape"`
	Timeout time.Duration `yaml:"timeout" env-default:"30s"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env.local")
	if err != nil {
		fmt.Printf("%s", err.Error())
		log.Fatalf("Cannot load .env")
	}

	cfgPath := os.Getenv("CFG_PATH")
	if cfgPath == "" {
		log.Fatalf("CFG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err.Error())
	}

	return &cfg
}
