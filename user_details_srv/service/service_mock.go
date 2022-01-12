package service

import (
	"context"
	"os"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/stretchr/testify/mock"
)

type UserDetailsRepositoryMock struct {
	mock.Mock
}

func (r *UserDetailsRepositoryMock) SetUserDetails(ctx context.Context, information entities.UserDetails) (bool, error) {
	args := r.Called(ctx, information)

	return args.Bool(0), args.Error(1)
}

func (r *UserDetailsRepositoryMock) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	args := r.Called(ctx, user_id)

	return args.Get(0).(entities.UserDetails), args.Error(1)
}

func (r *UserDetailsRepositoryMock) DeleteUserDetails(ctx context.Context, user_id int) (bool, error) {
	args := r.Called(ctx, user_id)

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
			"user_details",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}
	return logger
}
