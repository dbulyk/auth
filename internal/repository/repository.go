package repository

import (
	"context"

	"auth/internal/model"
)

// AuthRepository определяет взаимодействие с бд
type AuthRepository interface {
	User
}

// User описывает взаимодействие с репозиторием пользователя
type User interface {
	CreateUser(ctx context.Context, user model.CreateUser) (id int64, err error)
	UpdateUser(ctx context.Context, user model.UpdateUser) (err error)
	GetUser(ctx context.Context, id int64) (user *model.User, err error)
	DeleteUser(ctx context.Context, id int64) (err error)
}
