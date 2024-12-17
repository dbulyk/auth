package user

import (
	"context"

	"auth/internal/model"
)

// GetUserServ является сервисной прослойкой для получения данных о пользователе
func (s *serv) GetUserServ(ctx context.Context, id int64) (user *model.User, err error) {
	user, err = s.authRepository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
