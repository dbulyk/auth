package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/service/mocks"
	desc "auth/pkg/auth_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name            = gofakeit.Name()
		email           = gofakeit.Email()
		tag             = gofakeit.Gamertag()
		password        = gofakeit.Password(false, false, false, false, false, 6)
		passwordConfirm = password
		role            = desc.Role_ROLE_USER

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateUserRequest{
			Name:            name,
			Email:           email,
			Tag:             tag,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		modelReq = &model.UpdateUserRequest{
			Name:            name,
			Email:           email,
			Tag:             tag,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role.String(),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Успешное апдейт пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelReq).Return(nil)
				return mock
			},
		},
		{
			name: "Ошибка апдейта данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, modelReq).Return(serviceErr)
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

			newID, err := api.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
