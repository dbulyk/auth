package user

import (
	"context"

	"auth/internal/model"
)

// UpdateUserServ является сервисной прослойкой для удаления пользователя
func (s *serv) UpdateUserServ(ctx context.Context, user model.UpdateUser) (err error) {
	err = s.authRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
