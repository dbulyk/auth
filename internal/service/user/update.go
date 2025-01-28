package user

import (
	"auth/internal/model"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *service) UpdateUser(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepo.UpdateUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}
