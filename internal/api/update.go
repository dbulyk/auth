package api

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "auth/pkg/auth_v1"
)

func (i *Implementation) UpdateUser(ctx context.Context, in *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.Update(ctx, in)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil

}
