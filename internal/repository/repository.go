package repository

import (
	"auth/internal/repository/user/model"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserRepository interface {
	CreateUser(ctx context.Context, in *model.CreateUserRequest) (*model.CreateUserResponse, error)
	GetUser(ctx context.Context, in *model.GetUserRequest) (*model.GetUserResponse, error)
	UpdateUser(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, in *model.DeleteUserRequest) (*emptypb.Empty, error)
}
