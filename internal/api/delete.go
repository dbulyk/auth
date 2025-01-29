package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "auth/pkg/auth_v1"
)

func (i *Implementation) DeleteUser(ctx context.Context, in *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
