package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/service/mocks"
	desc "auth/pkg/auth_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = gofakeit.Int64()
		name            = gofakeit.Name()
		email           = gofakeit.Email()
		tag             = gofakeit.Gamertag()
		password        = gofakeit.Password(false, false, false, false, false, 6)
		passwordConfirm = password
		role            = desc.Role_ROLE_USER

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateUserRequest{
			Name:            name,
			Email:           email,
			Tag:             tag,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		modelReq = &model.CreateUserRequest{
			Name:            name,
			Email:           email,
			Tag:             tag,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role.String(),
		}

		res = &desc.CreateUserResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Успешное создание пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, modelReq).Return(id, nil)
				return mock
			},
		},
		{
			name: "Ошибка создания пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, modelReq).Return(0, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(chatServiceMock)

			newID, err := api.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
