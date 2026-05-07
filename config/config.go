package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type (
	// Config stores runtime settings for the Gophermart service.
	Config struct {
		RunAddress           string `env:"RUN_ADDRESS"`
		DatabaseURI          string `env:"DATABASE_URI"`
		AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
		JWTSecret            string `env:"JWT_SECRET" envDefault:"gophermart-secret"`
		LogLevel             string `env:"LOG_LEVEL" envDefault:"info"`
		HTTP                 HTTPConfig
		PG                   PGConfig
		JWT                  JWTConfig
		Accrual              AccrualConfig
	}

	// HTTPConfig stores HTTP server settings.
	HTTPConfig struct {
		Port           string
		UsePreforkMode bool `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// PGConfig stores PostgreSQL connection settings.
	PGConfig struct {
		PoolMax int `env:"PG_POOL_MAX" envDefault:"2"`
	}

	// JWTConfig stores JWT token settings.
	JWTConfig struct {
		TokenExpiry time.Duration `env:"JWT_TOKEN_EXPIRY" envDefault:"24h"`
	}

	// AccrualConfig stores background accrual polling settings.
	AccrualConfig struct {
		PollInterval time.Duration `env:"ACCRUAL_POLL_INTERVAL" envDefault:"5s"`
	}
)

// NewConfig returns application config.
func NewConfig() (*Config, error) {
	cfg := &Config{
		RunAddress: "localhost:8080",
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "server run address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "postgres connection string")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual service address")
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	cfg.HTTP.Port = cfg.RunAddress
	cfg.PG = PGConfig{PoolMax: cfg.PG.PoolMax}
	cfg.JWT = JWTConfig{TokenExpiry: cfg.JWT.TokenExpiry}
	cfg.Accrual = AccrualConfig{PollInterval: cfg.Accrual.PollInterval}

	return cfg, nil
}
