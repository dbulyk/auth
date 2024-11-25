package user

import (
	"auth/internal/repository"
	"auth/internal/service"
)

var _ service.AuthService = (*serv)(nil)

type serv struct {
	authRepository repository.AuthRepository
}

// NewAuthService создаёт новый объект сервиса пользователей
func NewAuthService(authRepository repository.AuthRepository) *serv {
	return &serv{authRepository: authRepository}
}
