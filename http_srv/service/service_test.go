package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	var http_service service.HTTPServicer
	repository_mock := new(service.RepoMock)
	logger := log.NewLogfmtLogger(os.Stderr)
	http_service = service.NewHTTPService(repository_mock, logger)

	testCases := []struct {
		testName string
		data     entities.User
		res      int
		err      error
	}{
		{
			testName: "user created successfully",
			data: entities.User{
				Email:    "success@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  service.GenenerateDetails(),
			},
			res: 1,
			err: nil,
		},
		{
			testName: "no email error",
			data: entities.User{
				Password: "qwerty",
				Age:      23,
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			testName: "no password error",
			data: entities.User{
				Email: "success@email.com",
				Age:   23,
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repository_mock.On("CreateUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := http_service.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details)

			// assert
			assert.True(service.TestErrors(err, tc.err))
			assert.Equal(tc.res, res)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	var http_service service.HTTPServicer
	repository_mock := new(service.RepoMock)
	logger := log.NewLogfmtLogger(os.Stderr)
	http_service = service.NewHTTPService(repository_mock, logger)

	testCases := []struct {
		testName string
		data     entities.Session
		res      bool
		err      error
	}{
		{
			testName: "success authenticate",
			data: entities.Session{
				Email:    "fake_email@email.com",
				Password: "fake_password",
			},
			res: true,
			err: nil,
		},
		{
			testName: "no email error",
			data: entities.Session{
				Password: "fake_password",
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
		{
			testName: "no password error",
			data: entities.Session{
				Email: "fake_email@email.com",
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			testName: "user not found error",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			err: status.Error(codes.NotFound, "User not found"),
		},
		{
			testName: "invalid password error",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			err: status.Error(codes.Unauthenticated, "Password or email error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repository_mock.On("Authenticate", ctx, tc.data).Return(tc.res, tc.err)
			res, err := http_service.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.True(service.TestErrors(err, tc.err))
		})
	}
}

func TestUpdateUser(t *testing.T) {
	var http_service service.HTTPServicer
	repository_mock := new(service.RepoMock)
	logger := log.NewLogfmtLogger(os.Stderr)
	http_service = service.NewHTTPService(repository_mock, logger)

	testCases := []struct {
		testName string
		data     entities.UserUpdate
		res      bool
		err      error
	}{
		{
			testName: "update user success",
			data: entities.UserUpdate{
				UserID: 1,
				User: entities.User{
					Email:    "new_email@domain.com",
					Password: "new_password",
					Age:      23,
					Details:  service.GenenerateDetails(),
				},
			},
			res: true,
			err: nil,
		},
		{
			testName: "no email error",
			data: entities.UserUpdate{
				UserID: 1,
				User: entities.User{
					Password: "new_password",
					Age:      23,
					Details:  service.GenenerateDetails(),
				},
			},
			res: false,
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
		{
			testName: "no password error",
			data: entities.UserUpdate{
				UserID: 1,
				User: entities.User{
					Email:   "new_email@domain.com",
					Age:     23,
					Details: service.GenenerateDetails(),
				},
			},
			res: false,
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			testName: "user not found error",
			data: entities.UserUpdate{
				UserID: 2,
				User: entities.User{
					Email:    "new_email@domain.com",
					Password: "new_password",
					Age:      23,
					Details:  service.GenenerateDetails(),
				},
			},
			res: false,
			err: status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		// act
		repository_mock.On("UpdateUser", ctx, tc.data).Return(tc.res, tc.err)
		res, err := http_service.UpdateUser(ctx, tc.data.UserID, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details)

		// assert
		assert.Equal(tc.res, res)
		assert.True(service.TestErrors(err, tc.err))
	}
}

func TestGetUser(t *testing.T) {
	var http_service service.HTTPServicer
	repository_mock := new(service.RepoMock)
	logger := log.NewLogfmtLogger(os.Stderr)
	http_service = service.NewHTTPService(repository_mock, logger)

	testCases := []struct {
		testName string
		data     int
		res      entities.User
		err      error
	}{
		{
			testName: "user found success",
			data:     1,
			res: entities.User{
				Email:    "email@domain.com",
				Password: "password",
				Age:      23,
			},
			err: nil,
		},

		{
			testName: "user found success",
			data:     2,
			err:      status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		// act
		repository_mock.On("GetUser", ctx, tc.data).Return(tc.res, tc.err)
		res, err := http_service.GetUser(ctx, tc.data)

		// assert
		assert.Equal(tc.res, res)
		assert.True(service.TestErrors(err, tc.err))
	}
}

func TestDeleteUser(t *testing.T) {
	var http_service service.HTTPServicer
	repository_mock := new(service.RepoMock)
	logger := log.NewLogfmtLogger(os.Stderr)
	http_service = service.NewHTTPService(repository_mock, logger)

	testCases := []struct {
		testName string
		data     int
		res      bool
		err      error
	}{
		{
			testName: "user deleted success",
			data:     1,
			res:      true,
			err:      nil,
		},
		{
			testName: "user not found error",
			data:     2,
			res:      false,
			err:      status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		// act
		repository_mock.On("DeleteUser", ctx, tc.data).Return(tc.res, tc.err)
		res, err := http_service.DeleteUser(ctx, tc.data)

		// assert
		assert.Equal(tc.res, res)
		assert.True(service.TestErrors(err, tc.err))
	}
}
