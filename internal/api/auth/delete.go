package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "auth/pkg/auth_v1"
)

// DeleteUser удаляет пользователя по id
func (i *Implementation) DeleteUser(ctx context.Context, in *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	err := i.authService.DeleteUserServ(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления пользователя: %v", err)
	}

	return &emptypb.Empty{}, nil
}
