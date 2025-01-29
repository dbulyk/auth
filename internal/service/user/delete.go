package user

import (
	"context"
)

func (s *service) Delete(ctx context.Context, userID int64) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
