package api

import (
	"auth/internal/service"
	desc "auth/pkg/auth_v1"
)

// Implementation описывает модель данных сервиса
type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
}

// NewImplementation возвращает объект имплементации сервиса
func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
