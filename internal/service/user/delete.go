package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *service) Delete(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	_, err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
