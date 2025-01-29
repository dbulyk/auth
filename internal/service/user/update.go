package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/model"
)

func (s *service) Update(ctx context.Context, in *model.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepo.UpdateUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, err
}
