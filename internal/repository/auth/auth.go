package auth

import (
	"auth/internal/client/db"
	"auth/internal/repository"
)

var _ repository.AuthRepository = (*Repo)(nil)

// Repo подключает модули для взаимодействия с пользователем
type Repo struct {
	repoUser
}

// NewRepository создаёт репозиторий для действий с пользователями
func NewRepository(db db.Client, key string) *Repo {
	return &Repo{
		repoUser: repoUser{db: db, hashKey: key},
	}
}
