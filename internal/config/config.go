package config

import (
	"time"

	"github.com/joho/godotenv"
)

// PGConfig хранит в себе строку подключения к бд
type PGConfig interface {
	DSN() string
}

// GRPCConfig хранит адрес, на котором поднимается сервер
type GRPCConfig interface {
	Address() string
}

// HashConfig хранит ключ для хеширования паролей
type HashConfig interface {
	Key() string
}

// RedisConfig описывает контракт взаимодействия с конфигом редиса
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// Load загружает во флаг считываемый путь
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
