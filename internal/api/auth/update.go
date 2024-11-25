package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/model"
	desc "auth/pkg/auth_v1"
)

// UpdateUser обновляет данные пользователя
func (i *Implementation) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	updateUser := model.UpdateUser{
		Name:            in.GetName(),
		Email:           in.GetEmail(),
		Tag:             in.GetTag(),
		Password:        in.GetPassword(),
		PasswordConfirm: in.GetPasswordConfirm(),
	}
	err := i.authService.UpdateUserServ(ctx, updateUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обновления данных пользователя: %v", err)
	}

	return &emptypb.Empty{}, nil
}
