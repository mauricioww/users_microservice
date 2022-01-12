package service_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.User
		res       int
		err       error
	}{
		{
			test_name: "user created successfully",
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
			test_name: "no email error",
			data: entities.User{
				Password: "qwerty",
				Age:      23,
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			test_name: "no password error",
			data: entities.User{
				Email: "success@email.com",
				Age:   23,
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
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
	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.Session
		res       int
		err       error
	}{
		{
			test_name: "success authenticate",
			data: entities.Session{
				Email:    "fake_email@email.com",
				Password: "fake_password",
			},
			res: 0,
			err: nil,
		},
		{
			test_name: "no email error",
			data: entities.Session{
				Password: "fake_password",
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
		{
			test_name: "no password error",
			data: entities.Session{
				Email: "fake_email@email.com",
			},
			res: -1,
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			test_name: "user not found error",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			res: -1,
			err: status.Error(codes.NotFound, "User not found"),
		},
		{
			test_name: "invalid password error",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			res: -1,
			err: status.Error(codes.Unauthenticated, "Password or email error"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repository_mock.On("Authenticate", ctx, tc.data).Return(tc.res, tc.err)
			res, err := http_service.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			if tc.res >= 0 {
				assert.NotEmpty(res)
			} else {
				assert.Empty(res)
			}
			assert.True(service.TestErrors(err, tc.err))
		})
	}
}

func TestUpdateUser(t *testing.T) {
	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.UserUpdate
		res       bool
		err       error
	}{
		{
			test_name: "update user success",
			data: entities.UserUpdate{
				UserId: 1,
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
			test_name: "no email error",
			data: entities.UserUpdate{
				UserId: 1,
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
			test_name: "no password error",
			data: entities.UserUpdate{
				UserId: 1,
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
			test_name: "user not found error",
			data: entities.UserUpdate{
				UserId: 2,
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

	for _, tc := range test_cases {
		// prepare
		ctx := context.Background()
		assert := assert.New(t)

		// act
		repository_mock.On("UpdateUser", ctx, tc.data).Return(tc.res, tc.err)
		res, err := http_service.UpdateUser(ctx, tc.data.UserId, tc.data.Email, tc.data.Password, tc.data.Age, tc.data.Details)

		// assert
		assert.Equal(tc.res, res)
		assert.True(service.TestErrors(err, tc.err))
	}
}

func TestGetUser(t *testing.T) {
	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       entities.User
		err       error
	}{
		{
			test_name: "user found success",
			data:      1,
			res: entities.User{
				Email:    "email@domain.com",
				Password: "password",
				Age:      23,
			},
			err: nil,
		},

		{
			test_name: "user found success",
			data:      2,
			err:       status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range test_cases {
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
	repository_mock := new(service.RepoMock)
	http_service := service.NewHttpService(repository_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       bool
		err       error
	}{
		{
			test_name: "user deleted success",
			data:      1,
			res:       true,
			err:       nil,
		},
		{
			test_name: "user not found error",
			data:      2,
			res:       false,
			err:       status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range test_cases {
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
