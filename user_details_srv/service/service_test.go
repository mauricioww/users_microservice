package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
	"github.com/stretchr/testify/assert"
)

func TestSetUserDetails(t *testing.T) {
	var srv service.GrpcUserDetailsServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserDetailsRepositoryMock)
	srv = service.NewGrpcUserDetailsService(repoMock, logger)

	testCases := []struct {
		testName string
		data     entities.UserDetails
		res      bool
		err      error
	}{
		{
			testName: "set user details which no exists success",
			data: entities.UserDetails{
				UserID:       1,
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

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("SetUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := srv.SetUserDetails(ctx, tc.data.UserID, tc.data.Country, tc.data.City,
				tc.data.MobileNumber, tc.data.Married, tc.data.Height, tc.data.Weight)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	var srv service.GrpcUserDetailsServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserDetailsRepositoryMock)
	srv = service.NewGrpcUserDetailsService(repoMock, logger)

	testCases := []struct {
		testName string
		data     int
		res      entities.UserDetails
		err      error
	}{
		{
			testName: "get user details success",
			data:     0,
			res: entities.UserDetails{
				UserID:       0,
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
			testName: "get user details which does not exist error",
			data:     1,
			err:      errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("GetUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := srv.GetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUserDetails(t *testing.T) {
	var srv service.GrpcUserDetailsServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserDetailsRepositoryMock)
	srv = service.NewGrpcUserDetailsService(repoMock, logger)

	testCases := []struct {
		testName string
		data     int
		res      bool
		err      error
	}{
		{
			testName: "delete user details success",
			data:     0,
			res:      true,
			err:      nil,
		},
		{
			testName: "delete user details which does not exist error",
			data:     1,
			err:      errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("DeleteUserDetails", ctx, tc.data).Return(tc.res, tc.err)
			res, err := srv.DeleteUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
