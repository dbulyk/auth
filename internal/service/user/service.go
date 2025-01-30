package user

import (
	"auth/internal/client/db"
	"auth/internal/repository"
)

type service struct {
	userRepo  repository.UserRepository
	txManager db.TxManager
}

// NewUserService создаёт и возвращает новый объект сервиса пользователя
func NewUserService(userRepo repository.UserRepository, txManager db.TxManager) *service {
	return &service{userRepo: userRepo, txManager: txManager}
}
