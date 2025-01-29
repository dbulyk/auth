package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Create(ctx context.Context, in *model.CreateUserRequest) (int64, error) {
	userID, err := s.userRepo.CreateUser(ctx, in)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
