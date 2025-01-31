package user

import (
	"context"

	"auth/internal/converter"
	desc "auth/pkg/auth_v1"
)

// CreateUser является имплементацией api для создания пользователя
func (i *Implementation) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	userID, err := i.userService.Create(ctx, converter.ToCreateUserRequestFromAPI(in))
	if err != nil {
		return nil, err
	}
	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}
