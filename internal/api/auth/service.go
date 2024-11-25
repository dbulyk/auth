package auth

import (
	"auth/internal/service"
	desc "auth/pkg/auth_v1"
)

// Implementation является объектом сервера
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation создаёт объект сервера
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
