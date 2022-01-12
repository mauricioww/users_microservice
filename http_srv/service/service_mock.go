package service

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
)

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) CreateUser(ctx context.Context, user entities.User) (int, error) {
	args := r.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

func (r *RepoMock) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	args := r.Called(ctx, session)

	return args.Int(0), args.Error(1)
}

func (r *RepoMock) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	args := r.Called(ctx, user)

	return args.Bool(0), args.Error(1)
}

func (r *RepoMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (r *RepoMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := r.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

func GenenerateDetails() entities.Details {
	return entities.Details{
		Country:      "Mexico",
		City:         "CDMX",
		MobileNumber: "11223344",
		Married:      false,
		Height:       1.75,
		Weight:       76.0,
	}
}

func InitLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"account",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}
	return logger
}

func TestErrors(err1 error, err2 error) bool {
	e1 := status.Convert(err1)
	e2 := status.Convert(err2)
	return (e1.Code() == e2.Code()) && (e1.Message() == e2.Message())
}
