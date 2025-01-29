package user

import (
	"context"

	"auth/internal/model"
)

func (s *service) Update(ctx context.Context, in *model.UpdateUserRequest) error {
	err := s.userRepo.UpdateUser(ctx, in)
	if err != nil {
		return err
	}
	return err
}
