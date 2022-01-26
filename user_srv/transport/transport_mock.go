package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/stretchr/testify/mock"
)

// GrpcUserSrvMock type is used to mock the performance of the service layer
type GrpcUserSrvMock struct {
	mock.Mock
}

// CreateUser is a mock of the real method
func (s *GrpcUserSrvMock) CreateUser(ctx context.Context, email string, pwd string, age int) (int, error) {
	args := s.Called(ctx, email, pwd, age)

	return args.Int(0), args.Error(1)
}

// Authenticate is a mock of the real method
func (s *GrpcUserSrvMock) Authenticate(ctx context.Context, email string, pwd string) (bool, error) {
	args := s.Called(ctx, email, pwd)

	return args.Bool(0), args.Error(1)
}

// UpdateUser is a mock of the real method
func (s *GrpcUserSrvMock) UpdateUser(ctx context.Context, id int, email string, pwd string, age int) (bool, error) {
	args := s.Called(ctx, id, email, pwd, age)

	return args.Bool(0), args.Error(1)
}

// GetUser is a mock of the real method
func (s *GrpcUserSrvMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := s.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

// DeleteUser is a mock of the real method
func (s *GrpcUserSrvMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := s.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}
