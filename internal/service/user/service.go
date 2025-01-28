package user

import "auth/internal/repository"

type service struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *service {
	return &service{userRepo: userRepo}
}
