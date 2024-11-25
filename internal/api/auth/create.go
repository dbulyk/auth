package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth/internal/model"
	desc "auth/pkg/auth_v1"
)

// CreateUser создает пользователя
func (i *Implementation) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	user := model.CreateUser{
		Name:            in.GetName(),
		Email:           in.GetEmail(),
		Tag:             in.GetTag(),
		Role:            int32(in.GetRole()),
		Password:        in.GetPassword(),
		PasswordConfirm: in.GetPasswordConfirm(),
	}

	userID, err := i.authService.CreateUserServ(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания пользователя: %v", err)
	}

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}
