package service_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
	"github.com/stretchr/testify/assert"
)

func TestSetUserDetails(t *testing.T) {
	var grpc_user_details_srv service.GrpcUserDetailsService

	user_details_repo_mock := new(service.UserDetailsRepositoryMock)
	grpc_user_details_srv = service.NewGrpcUserDetailsService(user_details_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.UserDetails
		res       bool
		err       error
	}{
		{
			test_name: "set user details which no exists success",
			data: entities.UserDetails{
				UserId:       1,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			res: true,
			err: nil,
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_details_repo_mock.On("SetUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_details_srv.SetUserDetails(ctx, tc.data.UserId, tc.data.Country, tc.data.City,
				tc.data.MobileNumber, tc.data.Married, tc.data.Height, tc.data.Weight)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	var grpc_user_details_srv service.GrpcUserDetailsService

	user_details_repo_mock := new(service.UserDetailsRepositoryMock)
	grpc_user_details_srv = service.NewGrpcUserDetailsService(user_details_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       entities.UserDetails
		err       error
	}{
		{
			test_name: "get user details success",
			data:      0,
			res: entities.UserDetails{
				UserId:       0,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			err: nil,
		},
		{
			test_name: "get user details which does not exist error",
			data:      1,
			err:       errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_details_repo_mock.On("GetUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_details_srv.GetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUserDetails(t *testing.T) {
	var grpc_user_details_srv service.GrpcUserDetailsService

	user_details_repo_mock := new(service.UserDetailsRepositoryMock)
	grpc_user_details_srv = service.NewGrpcUserDetailsService(user_details_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       bool
		err       error
	}{
		{
			test_name: "delete user details success",
			data:      0,
			res:       true,
			err:       nil,
		},
		{
			test_name: "delete user details which does not exist error",
			data:      1,
			err:       errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_details_repo_mock.On("DeleteUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_details_srv.DeleteUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
