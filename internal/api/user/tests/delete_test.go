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
	"auth/internal/service"
	"auth/internal/service/mocks"
	desc "auth/pkg/auth_v1"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.DeleteUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID = gofakeit.Int64()

		req = &desc.DeleteUserRequest{Id: userID}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		chatServiceMock userServiceMockFunc
	}{
		{
			name: "Успешное удаление данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			chatServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, userID).Return(nil)
				return mock
			},
		},
		{
			name: "Ошибка удаления данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, userID).Return(serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatServiceMock := tt.chatServiceMock(mc)
			api := user.NewImplementation(chatServiceMock)

			newID, err := api.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
