package transport

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/stretchr/testify/mock"
)

type GrpcUserDetailsSrvMock struct {
	mock.Mock
}

func (g *GrpcUserDetailsSrvMock) SetUserDetails(ctx context.Context, user_id int, country string, city string, mobile_number string, married bool, height float32, weigth float32) (bool, error) {
	args := g.Called(ctx, user_id, country, city, mobile_number, married, height, weigth)

	return args.Bool(0), args.Error(1)
}

func (g *GrpcUserDetailsSrvMock) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	args := g.Called(ctx, user_id)

	return args.Get(0).(entities.UserDetails), args.Error(1)
}

func (g *GrpcUserDetailsSrvMock) DeleteUserDetails(ctx context.Context, user_id int) (bool, error) {
	args := g.Called(ctx, user_id)

	return args.Bool(0), args.Error(1)
}
