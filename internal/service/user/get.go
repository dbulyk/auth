package user

import (
	"auth/internal/model"
	"context"
)

func (s *service) GetUser(ctx context.Context, userID int64) (*model.GetUserResponse, error) {
	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
