package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"auth/internal/repository"
)

var _ repository.AuthRepository = (*Repo)(nil)

// Repo подключает модули для взаимодействия с пользователем
type Repo struct {
	repoUser
}

// NewRepository создаёт репозиторий для действий с пользователями
func NewRepository(db *pgxpool.Pool, key string) *Repo {
	return &Repo{
		repoUser: repoUser{db: db, hashKey: key},
	}
}
