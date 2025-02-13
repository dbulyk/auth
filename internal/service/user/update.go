package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Update(ctx context.Context, in *model.UpdateUserRequest) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		txErr = s.userRepo.UpdateUser(ctx, in)
		return txErr
	})
	if err != nil {
		return err
	}
	return nil
}
