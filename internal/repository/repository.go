package repository

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, in *model.CreateUserRequest) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.GetUserResponse, error)
	UpdateUser(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}
