package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/stretchr/testify/mock"
)

// ServiceMock type is used to mock the performance of the service layer
type ServiceMock struct {
	mock.Mock
}

// CreateUser is a mock of the real method
func (s *ServiceMock) CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error) {
	args := s.Called(ctx, email, pwd, age, details)

	return args.Int(0), args.Error(1)
}

// Authenticate is a mock of the real method
func (s *ServiceMock) Authenticate(ctx context.Context, email string, pwd string) (bool, error) {
	args := s.Called(ctx, email, pwd)

	return args.Bool(0), args.Error(1)
}

// UpdateUser is a mock of the real method
func (s *ServiceMock) UpdateUser(ctx context.Context, userID int, email string, pwd string, age int, details entities.Details) (bool, error) {
	args := s.Called(ctx, userID, email, pwd, age, details)

	return args.Bool(0), args.Error(1)
}

// GetUser is a mock of the real method
func (s *ServiceMock) GetUser(ctx context.Context, userID int) (entities.User, error) {
	args := s.Called(ctx, userID)

	return args.Get(0).(entities.User), args.Error(1)
}

// DeleteUser is a mock of the real method
func (s *ServiceMock) DeleteUser(ctx context.Context, userID int) (bool, error) {
	args := s.Called(ctx, userID)

	return args.Bool(0), args.Error(1)
}
