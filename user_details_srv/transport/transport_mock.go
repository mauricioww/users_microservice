package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/stretchr/testify/mock"
)

// GrpcUserDetailsSrvMock type is used to mock the performance of the service layer
type GrpcUserDetailsSrvMock struct {
	mock.Mock
}

// SetUserDetails is a mock of the real method
func (g *GrpcUserDetailsSrvMock) SetUserDetails(ctx context.Context, userID int, country string, city string, number string, married bool, height float32, weigth float32) (bool, error) {
	args := g.Called(ctx, userID, country, city, number, married, height, weigth)

	return args.Bool(0), args.Error(1)
}

// GetUserDetails is a mock of the real method
func (g *GrpcUserDetailsSrvMock) GetUserDetails(ctx context.Context, userID int) (entities.UserDetails, error) {
	args := g.Called(ctx, userID)

	return args.Get(0).(entities.UserDetails), args.Error(1)
}

// DeleteUserDetails is a mock of the real method
func (g *GrpcUserDetailsSrvMock) DeleteUserDetails(ctx context.Context, userID int) (bool, error) {
	args := g.Called(ctx, userID)

	return args.Bool(0), args.Error(1)
}
