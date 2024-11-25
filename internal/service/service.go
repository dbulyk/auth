package service

import (
	"context"

	"auth/internal/model"
)

// AuthService является сервисной прослойкой для пользовательского репозитория
type AuthService interface {
	CreateUserServ(ctx context.Context, user model.CreateUser) (id int64, err error)
	UpdateUserServ(ctx context.Context, user model.UpdateUser) (err error)
	GetUserServ(ctx context.Context, id int64) (user *model.User, err error)
	DeleteUserServ(ctx context.Context, id int64) (err error)
}
