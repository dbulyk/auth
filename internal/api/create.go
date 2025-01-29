package api

import (
	"context"

	desc "auth/pkg/auth_v1"
)

func (i *Implementation) CreateUser(ctx context.Context, in *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	userID, err := i.userService.Create(ctx, in)
	if err != nil {
		return nil, err
	}
	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}
