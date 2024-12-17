package config

import (
	"errors"
	"os"
)

var _ HashConfig = (*hashConfig)(nil)

const key = "HASH_KEY"

type hashConfig struct {
	key string
}

// NewHashConfig создает конфиг для хранения соли для хеширования паролей
func NewHashConfig() (*hashConfig, error) {
	k := os.Getenv(key)
	if len(k) == 0 {
		return nil, errors.New("key for hash not found")
	}

	return &hashConfig{key: k}, nil
}

func (h *hashConfig) Key() string { return h.key }
