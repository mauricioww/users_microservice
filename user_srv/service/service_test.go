package service_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.User
		res       int
		err       error
	}{
		{
			test_name: "create user successfully",
			data: entities.User{
				Email:    "email@domain.com",
				Password: "qwerty",
				Age:      23,
			},
			err: nil,
		},
		{
			test_name: "empty email error",
			data: entities.User{
				Password: "qwerty",
				Age:      23,
			},
			res: -1,
			err: errors.NewBadRequestEmailError(),
		},
		{
			test_name: "empty password error",
			data: entities.User{
				Email: "email@domain.com",
				Age:   23,
			},
			res: -1,
			err: errors.NewBadRequestPasswordError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("CreateUser", ctx, mock.AnythingOfType("entities.User")).Return(tc.res, nil)
			res, err := grpc_user_srv.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.Age)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      *entities.Session
		repo_pwd  string
		repo_err  error
		res       int
		err       error
	}{
		{
			test_name: "authenticate successfully",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "secret",
			},
			repo_pwd: "secret",
		},
		{
			test_name: "no pasword error",
			data: &entities.Session{
				Email: "user@email.com",
			},
			res: -1,
			err: errors.NewBadRequestPasswordError(),
		},
		{
			test_name: "no email error",
			data: &entities.Session{
				Password: "password",
			},
			res: -1,
			err: errors.NewBadRequestEmailError(),
		},
		{
			test_name: "user not found error",
			data: &entities.Session{
				Email:    "use@email.com",
				Password: "password",
			},
			repo_pwd: "password",
			repo_err: errors.NewUserNotFoundError(),
			res:      -1,
			err:      errors.NewUserNotFoundError(),
		},
		{
			test_name: "invalid pasword error",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "password",
			},
			repo_pwd: "incorrect_password",
			res:      -1,
			err:      errors.NewUnauthenticatedError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			tc.repo_pwd = helpers.Cipher(tc.repo_pwd)

			// act
			user_repo_mock.On("Authenticate", ctx, tc.data).Return(tc.repo_pwd, tc.repo_err)
			res, err := grpc_user_srv.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      entities.Update
		repo_res  entities.User
		res       bool
		err       error
	}{
		{
			test_name: "update all fields",
			data: entities.Update{
				UserId: 0,
				User: entities.User{
					Email:    "new_email@domain.com",
					Password: "new_password",
					Age:      20,
				},
			},
			repo_res: entities.User{
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      20,
			},
			res: true,
			err: nil,
		},
		{
			test_name: "no password error",
			data: entities.Update{
				UserId: 0,
				User: entities.User{
					Email: "new_email@domain.com",
					Age:   20,
				},
			},
			res: false,
			err: errors.NewBadRequestPasswordError(),
		},
		{
			test_name: "no email error",
			data: entities.Update{
				UserId: 0,
				User: entities.User{
					Password: "new_password",
					Age:      20,
				},
			},
			res: false,
			err: errors.NewBadRequestEmailError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			tc.repo_res.Password = helpers.Cipher(tc.data.Password)

			// act
			user_repo_mock.On("UpdateUser", ctx, mock.Anything).Return(tc.repo_res, tc.err)
			res, err := grpc_user_srv.UpdateUser(ctx, tc.data.UserId, tc.data.Email, tc.data.Password, tc.data.Age)

			// assert
			assert.Equal(tc.err, err)
			assert.Equal(tc.res, res)
		})
	}
}

func TestGetUser(t *testing.T) {
	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       entities.User
		err       error
	}{
		{
			test_name: "user found",
			data:      0,
			res: entities.User{
				Email:    "user@email.com",
				Password: "password",
				Age:      20,
			},
			err: nil,
		},
		{
			test_name: "user not found error",
			data:      -1,
			res:       entities.User{},
			err:       errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("GetUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_srv.GetUser(ctx, tc.data)
			tc.res.Password = helpers.Decipher(tc.res.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	var grpc_user_srv service.GrpcUserService

	user_repo_mock := new(service.UserRepositoryMock)
	grpc_user_srv = service.NewGrpcUserService(user_repo_mock, service.InitLogger())

	test_cases := []struct {
		test_name string
		data      int
		res       bool
		err       error
	}{
		{
			test_name: "delete user success",
			data:      1,
			res:       true,
			err:       nil,
		},
		{
			test_name: "user does not exist error",
			data:      -1,
			res:       false,
			err:       errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			user_repo_mock.On("DeleteUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := grpc_user_srv.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
