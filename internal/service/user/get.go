package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Get(ctx context.Context, userID int64) (*model.GetUserResponse, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
