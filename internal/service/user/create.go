package user

import (
	"auth/internal/model"
	"context"
)

func (s *service) CreateUser(ctx context.Context, in *model.CreateUserRequest) (int64, error) {
	userID, err := s.userRepo.CreateUser(ctx, in)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
