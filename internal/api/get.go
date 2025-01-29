package api

import (
	"context"

	desc "auth/pkg/auth_v1"
)

func (i *Implementation) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.userService.Get(ctx, in)
	if err != nil {
		return nil, err
	}
	return &desc.GetUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Tag:       user.Tag,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
