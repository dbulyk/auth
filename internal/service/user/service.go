package user

import (
	"github.com/dbulyk/platform_common/pkg/db"

	"auth/internal/repository/user"
)

type service struct {
	userRepo  user.Repository
	txManager db.TxManager
	userCache user.Cache
}

// NewUserService создаёт и возвращает новый объект сервиса пользователя
func NewUserService(userRepo user.Repository, txManager db.TxManager, userCache user.Cache) *service {
	return &service{userRepo: userRepo, txManager: txManager, userCache: userCache}
}
