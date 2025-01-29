package user

import "auth/internal/repository"

type service struct {
	userRepo repository.UserRepository
}

// NewUserService создаёт и возвращает новый объект сервиса пользователя
func NewUserService(userRepo repository.UserRepository) *service {
	return &service{userRepo: userRepo}
}
