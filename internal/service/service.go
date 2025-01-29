package service

import (
	"context"

	"auth/internal/model"
)

// UserService описывает контракт для сервиса пользователя
type UserService interface {
	Create(ctx context.Context, in *model.CreateUserRequest) (int64, error)
	Get(ctx context.Context, userID int64) (*model.GetUserResponse, error)
	Update(ctx context.Context, in *model.UpdateUserRequest) error
	Delete(ctx context.Context, userID int64) error
}
