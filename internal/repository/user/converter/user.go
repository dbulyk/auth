package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
)

// ToUserFromRepo является конвертером получения пользователя с модели данных репозитория в сервисную модель данных
func ToUserFromRepo(user *modelRepo.GetUserResponse) *model.GetUserResponse {
	return &model.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Tag:       user.Tag,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
