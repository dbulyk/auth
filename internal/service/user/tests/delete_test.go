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
	"auth/internal/repository"
	"auth/internal/repository/mocks"
	"auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		userID  = gofakeit.Int64()
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
			name: "Успешное удаление пользователя",
			args: args{
				ctx: ctx,
				req: userID,
			},
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, userID).Return(nil)
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
			name: "Ошибка удаления пользователя",
			args: args{
				ctx: ctx,
				req: userID,
			},
			err: repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, userID).Return(repoErr)
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

			err := serv.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
