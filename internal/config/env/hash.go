package env

import (
	"errors"
	"os"

	"auth/internal/config"
)

var _ config.HashConfig = (*hashConfig)(nil)

const key = "HASH_KEY"

type hashConfig struct {
	key string
}

// NewHashConfig создает конфиг для хранения соли для хеширования паролей
func NewHashConfig() (*hashConfig, error) {
	key := os.Getenv(key)
	if len(key) == 0 {
		return nil, errors.New("key for hash not found")
	}

	return &hashConfig{key: key}, nil
}

func (h *hashConfig) Key() string { return h.key }
