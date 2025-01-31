package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"auth/internal/client/db"
	mocks2 "auth/internal/client/db/mocks"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/mocks"
	"auth/internal/service/user"
	desc "auth/pkg/auth_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req *model.UpdateUserRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name            = gofakeit.Name()
		email           = gofakeit.Email()
		tag             = gofakeit.Gamertag()
		password        = gofakeit.Name()
		passwordConfirm = password
		role            = desc.Role_ROLE_USER
		userID          = gofakeit.Int64()

		req = &model.UpdateUserRequest{
			ID:              userID,
			Name:            name,
			Email:           email,
			Tag:             tag,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role.String(),
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name          string
		args          args
		err           error
		userRepoMock  chatRepositoryMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "Успешное обновление данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks2.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
		},
		{
			name: "Ошибка обновления данных пользователя",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(repoErr)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := mocks2.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMock(mc)
			txManagerMock := tt.txManagerMock(mc)
			serv := user.NewUserService(userRepoMock, txManagerMock)

			err := serv.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
