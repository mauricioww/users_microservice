package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
)

// RepoMock type is used to mock the performance of the repository layer
type RepoMock struct {
	mock.Mock
}

// CreateUser is a mock of the real method
func (r *RepoMock) CreateUser(ctx context.Context, user entities.User) (int, error) {
	args := r.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

// Authenticate is a mock of the real method
func (r *RepoMock) Authenticate(ctx context.Context, session entities.Session) (bool, error) {
	args := r.Called(ctx, session)

	return args.Bool(0), args.Error(1)
}

// UpdateUser is a mock of the real method
func (r *RepoMock) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	args := r.Called(ctx, user)

	return args.Bool(0), args.Error(1)
}

// GetUser is a mock of the real method
func (r *RepoMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

// DeleteUser is a mock of the real method
func (r *RepoMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := r.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}

// GenenerateDetails returns mock data to use in tests
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

// TestErrors validates the code and message from the status errors
func TestErrors(err1 error, err2 error) bool {
	e1 := status.Convert(err1)
	e2 := status.Convert(err2)
	return (e1.Code() == e2.Code()) && (e1.Message() == e2.Message())
}
