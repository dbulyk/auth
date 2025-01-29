package api

import (
	"auth/internal/service"
	desc "auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.UserService
	hashKey     string
}

func NewImplementation(userService service.UserService, hashKey string) *Implementation {
	return &Implementation{
		userService: userService,
		hashKey:     hashKey,
	}
}
