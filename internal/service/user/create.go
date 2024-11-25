package user

import (
	"context"

	"auth/internal/model"
)

// CreateUserServ является сервисной прослойкой для создания пользователя
func (s *serv) CreateUserServ(ctx context.Context, user model.CreateUser) (id int64, err error) {
	id, err = s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return -1, err
	}
	return id, nil
}
