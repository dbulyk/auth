package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"auth/internal/model"
	desc "auth/pkg/auth_v1"
)

// ToCreateUserRequestFromAPI конвертирует запрос создания пользователя из модели протобафа в сервисную модель
func ToCreateUserRequestFromAPI(in *desc.CreateUserRequest) *model.CreateUserRequest {
	return &model.CreateUserRequest{
		Email:           in.Email,
		Name:            in.Name,
		Tag:             in.Tag,
		Password:        in.Password,
		PasswordConfirm: in.PasswordConfirm,
		Role:            in.GetRole().String(),
	}
}

// ToUpdateUserRequestFromAPI конвертирует запрос на обновление данных пользователя
// с модели протобафа в сервисную модель
func ToUpdateUserRequestFromAPI(in *desc.UpdateUserRequest) *model.UpdateUserRequest {
	return &model.UpdateUserRequest{
		ID:              in.GetId(),
		Email:           in.GetEmail(),
		Name:            in.GetName(),
		Tag:             in.GetTag(),
		Password:        in.GetPassword(),
		PasswordConfirm: in.GetPasswordConfirm(),
		Role:            in.GetRole().String(),
	}
}

var roleMapping = map[string]desc.Role{
	"ENUM_NAME_UNSPECIFIED": desc.Role_ENUM_NAME_UNSPECIFIED,
	"ROLE_USER":             desc.Role_ROLE_USER,
	"ROLE_ADMIN":            desc.Role_ROLE_ADMIN,
}

// ToGetUserResponseFromService конвертирует запрос на получение данных пользователя
// из сервисной модели в модель протобафа
func ToGetUserResponseFromService(response *model.GetUserResponse) *desc.GetUserResponse {
	var updatedAt *timestamppb.Timestamp
	if response.UpdatedAt.Valid {
		updatedAt = timestamppb.New(response.UpdatedAt.Time)
	}
	role := roleMapping[response.Role]

	return &desc.GetUserResponse{
		Id:        response.ID,
		Email:     response.Email,
		Name:      response.Name,
		Tag:       response.Tag,
		Role:      role,
		CreatedAt: timestamppb.New(response.CreatedAt),
		UpdatedAt: updatedAt,
	}
}
