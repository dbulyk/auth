package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Get(ctx context.Context, userID int64) (*model.GetUserResponse, error) {
	var user *model.GetUserResponse
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		user, txErr = s.userRepo.GetUser(ctx, userID)
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
