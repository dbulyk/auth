package repository

import (
	"context"

	"auth/internal/model"
)

// UserRepository описывает контракт для репозитория пользователя
type UserRepository interface {
	CreateUser(ctx context.Context, in *model.CreateUserRequest) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.GetUserResponse, error)
	UpdateUser(ctx context.Context, in *model.UpdateUserRequest) error
	DeleteUser(ctx context.Context, userID int64) error
}
