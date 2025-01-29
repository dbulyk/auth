package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, in *model.CreateUserRequest) (int64, error)
	Get(ctx context.Context, userID int64) (*model.GetUserResponse, error)
	Update(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error)
	Delete(ctx context.Context, userID int64) (*emptypb.Empty, error)
}
