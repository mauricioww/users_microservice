package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/stretchr/testify/mock"
)

// UserRepositoryMock type is used to mock the performance of the repository layer
type UserRepositoryMock struct {
	mock.Mock
}

// CreateUser is a mock of the real method
func (r *UserRepositoryMock) CreateUser(ctx context.Context, user entities.User) (int, error) {
	args := r.Called(ctx, user)

	return args.Int(0), args.Error(1)
}

// Authenticate is a mock of the real method
func (r *UserRepositoryMock) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	args := r.Called(ctx, session)

	return args.String(0), args.Error(1)
}

// UpdateUser is a mock of the real method
func (r *UserRepositoryMock) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	args := r.Called(ctx, update)

	return args.Get(0).(entities.User), args.Error(1)
}

// GetUser is a mock of the real method
func (r *UserRepositoryMock) GetUser(ctx context.Context, id int) (entities.User, error) {
	args := r.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

// DeleteUser is a mock of the real method
func (r *UserRepositoryMock) DeleteUser(ctx context.Context, id int) (bool, error) {
	args := r.Called(ctx, id)

	return args.Bool(0), args.Error(1)
}
