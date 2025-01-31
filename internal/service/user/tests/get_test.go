package tests

import (
	"auth/internal/client/db"
	mocks2 "auth/internal/client/db/mocks"
	"auth/internal/model"
	"auth/internal/repository"
	"auth/internal/repository/mocks"
	"auth/internal/service/user"
	desc "auth/pkg/auth_v1"
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name      = gofakeit.Name()
		email     = gofakeit.Email()
		tag       = gofakeit.Gamertag()
		role      = desc.Role_ROLE_USER
		userID    = gofakeit.Int64()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		res = &model.GetUserResponse{
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
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name          string
		args          args
		want          *model.GetUserResponse
		err           error
		userRepoMock  chatRepositoryMockFunc
		txManagerMock txManagerMockFunc
	}{
		{
			name: "Успешное получение данных пользователя",
			args: args{
				ctx: ctx,
				req: userID,
			},
			want: res,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, userID).Return(res, nil)
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
			name: "Ошибка получения данных пользователя",
			args: args{
				ctx: ctx,
				req: userID,
			},
			err: repoErr,
			userRepoMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, userID).Return(nil, repoErr)
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

			chatId, err := serv.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, chatId)
		})
	}
}
