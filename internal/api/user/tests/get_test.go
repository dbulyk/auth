package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/service/mocks"
	desc "auth/pkg/auth_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		tag       = gofakeit.Gamertag()
		role      = desc.Role_ROLE_USER
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		serviceErr = fmt.Errorf("service error")

		res = &desc.GetUserResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Tag:       tag,
			Role:      role,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}

		modelRes = &model.GetUserResponse{
			ID:        id,
			Name:      name,
			Email:     email,
			Tag:       tag,
			Role:      role.String(),
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}

		req = &desc.GetUserRequest{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetUserResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Успешное получение данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(modelRes, nil)
				return mock
			},
		},
		{
			name: "Ошибка получения данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
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

			newID, err := api.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
