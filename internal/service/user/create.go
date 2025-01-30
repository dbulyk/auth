package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Create(ctx context.Context, in *model.CreateUserRequest) (int64, error) {
	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		userID, txErr = s.userRepo.CreateUser(ctx, in)
		return txErr
	})

	if err != nil {
		return 0, err
	}
	return userID, nil
}
