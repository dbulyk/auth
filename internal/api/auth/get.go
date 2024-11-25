package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "auth/pkg/auth_v1"
)

// GetUser получает пользователя по id
func (i *Implementation) GetUser(ctx context.Context, in *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.authService.GetUserServ(ctx, in.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка получения данных пользователя: %v", err)
	}

	createdAt := timestamppb.New(user.CreatedAt)
	updatedAt := timestamppb.New(user.UpdatedAt.Time)
	res := &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Tag:       user.Tag,
		Role:      desc.Role(user.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return res, nil
}
