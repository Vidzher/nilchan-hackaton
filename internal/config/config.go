package config

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server" env-required:"true"`
	Firecrawl   Firecrawl  `yaml:"firecrawl"`
	OpenRouter  OpenRouter `yaml:"openrouter"`
}

type Firecrawl struct {
	APIKey  string        `yaml:"api_key" env:"FIRECRAWL_KEY"`
	BaseURL string        `yaml:"base_url" env-default:"https://api.firecrawl.dev/v2/scrape"`
	Timeout time.Duration `yaml:"timeout" env-default:"30s"`
}

type OpenRouter struct {
	APIKey    string `env:"OPENROUTER_API_KEY" env-required:"true"`
	ModelName string `yaml:"model_name" env:"OPENROUTER_MODEL"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"0.0.0.0:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"40s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

func Load(path string) (*Config, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(absolutePath, &cfg); err != nil {
		return nil, fmt.Errorf("read config %s: %w", absolutePath, err)
	}

	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		cfg.HTTPServer.Address = net.JoinHostPort("0.0.0.0", port)
	}

	if cfg.StoragePath != "" && cfg.StoragePath != ":memory:" && !filepath.IsAbs(cfg.StoragePath) && !strings.HasPrefix(cfg.StoragePath, "file:") {
		cfg.StoragePath = filepath.Join(filepath.Dir(absolutePath), cfg.StoragePath)
	}

	return &cfg, nil
}
