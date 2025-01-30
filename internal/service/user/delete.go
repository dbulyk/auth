package user

import (
	"context"
)

func (s *service) Delete(ctx context.Context, userID int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		txErr = s.userRepo.DeleteUser(ctx, userID)
		return txErr
	})
	if err != nil {
		return err
	}
	return nil
}
