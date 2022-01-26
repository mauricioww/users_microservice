package service

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/stretchr/testify/mock"
)

// UserDetailsRepositoryMock type is used to mock the performance of the repository layer
type UserDetailsRepositoryMock struct {
	mock.Mock
}

// SetUserDetails is a mock of the real method
func (r *UserDetailsRepositoryMock) SetUserDetails(ctx context.Context, information entities.UserDetails) (bool, error) {
	args := r.Called(ctx, information)

	return args.Bool(0), args.Error(1)
}

// GetUserDetails is a mock of the real method
func (r *UserDetailsRepositoryMock) GetUserDetails(ctx context.Context, userID int) (entities.UserDetails, error) {
	args := r.Called(ctx, userID)

	return args.Get(0).(entities.UserDetails), args.Error(1)
}

// DeleteUserDetails is amock of the real method
func (r *UserDetailsRepositoryMock) DeleteUserDetails(ctx context.Context, userID int) (bool, error) {
	args := r.Called(ctx, userID)

	return args.Bool(0), args.Error(1)
}
