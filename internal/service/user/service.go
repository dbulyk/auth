package user

import (
	"github.com/dbulyk/platform_common/pkg/db"

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
