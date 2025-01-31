package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/converter"
	desc "auth/pkg/auth_v1"
)

// UpdateUser является имплементацией api для обновления данных пользователя
func (i *Implementation) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	err := i.userService.Update(ctx, converter.ToUpdateUserRequestFromAPI(in))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil

}
