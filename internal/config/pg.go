package config

import (
	"errors"
	"os"
)

var _ PGConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig инициализирует и возвращает конфиг для подключения к бд
func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN возвращает строку для подключения к бд
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
