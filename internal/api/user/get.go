package user

import (
	"context"

	"auth/internal/converter"
	desc "auth/pkg/auth_v1"
)

// GetUser является имплементацией api для получения данных пользователя
func (i *Implementation) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.userService.Get(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return converter.ToGetUserResponseFromService(user), nil
}
