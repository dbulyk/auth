package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
)

func ToUserFromRepo(user *modelRepo.GetUserResponse) *model.GetUserResponse {
	return &model.GetUserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Tag:       user.Tag,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
