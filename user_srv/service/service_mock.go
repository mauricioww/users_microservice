package service

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) CreateUser(ctx context.Context, user entities.User) (int, error) {
	args := r.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

func (r *UserRepositoryMock) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	args := r.Called(ctx, session)

	return args.String(0), args.Error(1)
}

func (r *UserRepositoryMock) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	args := r.Called(ctx, update)

	return args.Get(0).(entities.User), args.Error(1)
}

func (r *UserRepositoryMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (r *UserRepositoryMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := r.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

func InitLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"user_grpc",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}
	return logger
}
