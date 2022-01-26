package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	var srv service.GrpcUserServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserRepositoryMock)
	srv = service.NewGrpcUserService(repoMock, logger)

	testCases := []struct {
		testName string
		data     entities.User
		res      int
		err      error
	}{
		{
			testName: "create user successfully",
			data: entities.User{
				Email:    "email@domain.com",
				Password: "qwerty",
				Age:      23,
			},
			err: nil,
		},
		{
			testName: "empty email error",
			data: entities.User{
				Password: "qwerty",
				Age:      23,
			},
			res: -1,
			err: errors.NewBadRequestEmailError(),
		},
		{
			testName: "empty password error",
			data: entities.User{
				Email: "email@domain.com",
				Age:   23,
			},
			res: -1,
			err: errors.NewBadRequestPasswordError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("CreateUser", ctx, mock.AnythingOfType("entities.User")).Return(tc.res, nil)
			res, err := srv.CreateUser(ctx, tc.data.Email, tc.data.Password, tc.data.Age)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	var srv service.GrpcUserServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserRepositoryMock)
	srv = service.NewGrpcUserService(repoMock, logger)

	testCases := []struct {
		testName string
		data     *entities.Session
		repoPwd  string
		repoErr  error
		res      bool
		err      error
	}{
		{
			testName: "authenticate successfully",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "secret",
			},
			res:     true,
			repoPwd: "secret",
		},
		{
			testName: "no pasword error",
			data: &entities.Session{
				Email: "user@email.com",
			},
			err: errors.NewBadRequestPasswordError(),
		},
		{
			testName: "no email error",
			data: &entities.Session{
				Password: "password",
			},
			err: errors.NewBadRequestEmailError(),
		},
		{
			testName: "user not found error",
			data: &entities.Session{
				Email:    "use@email.com",
				Password: "password",
			},
			repoPwd: "password",
			repoErr: errors.NewUserNotFoundError(),
			err:     errors.NewUserNotFoundError(),
		},
		{
			testName: "invalid pasword error",
			data: &entities.Session{
				Email:    "user@email.com",
				Password: "password",
			},
			repoPwd: "incorrect_password",
			err:     errors.NewUnauthenticatedError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			tc.repoPwd = helpers.Cipher(tc.repoPwd)

			// act
			repoMock.On("Authenticate", ctx, tc.data).Return(tc.repoPwd, tc.repoErr)
			res, err := srv.Authenticate(ctx, tc.data.Email, tc.data.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	var srv service.GrpcUserServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserRepositoryMock)
	srv = service.NewGrpcUserService(repoMock, logger)

	testCases := []struct {
		testName string
		data     entities.Update
		repo_res entities.User
		res      bool
		err      error
	}{
		{
			testName: "update all fields",
			data: entities.Update{
				UserID: 0,
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
			testName: "no password error",
			data: entities.Update{
				UserID: 0,
				User: entities.User{
					Email: "new_email@domain.com",
					Age:   20,
				},
			},
			res: false,
			err: errors.NewBadRequestPasswordError(),
		},
		{
			testName: "no email error",
			data: entities.Update{
				UserID: 0,
				User: entities.User{
					Password: "new_password",
					Age:      20,
				},
			},
			res: false,
			err: errors.NewBadRequestEmailError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)
			tc.repo_res.Password = helpers.Cipher(tc.data.Password)

			// act
			repoMock.On("UpdateUser", ctx, mock.Anything).Return(tc.repo_res, tc.err)
			res, err := srv.UpdateUser(ctx, tc.data.UserID, tc.data.Email, tc.data.Password, tc.data.Age)

			// assert
			assert.Equal(tc.err, err)
			assert.Equal(tc.res, res)
		})
	}
}

func TestGetUser(t *testing.T) {
	var srv service.GrpcUserServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserRepositoryMock)
	srv = service.NewGrpcUserService(repoMock, logger)

	testCases := []struct {
		testName string
		data     int
		res      entities.User
		err      error
	}{
		{
			testName: "user found",
			data:     0,
			res: entities.User{
				Email:    "user@email.com",
				Password: "password",
				Age:      20,
			},
			err: nil,
		},
		{
			testName: "user not found error",
			data:     -1,
			res:      entities.User{},
			err:      errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("GetUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := srv.GetUser(ctx, tc.data)
			tc.res.Password = helpers.Decipher(tc.res.Password)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	var srv service.GrpcUserServicer
	logger := log.NewLogfmtLogger(os.Stderr)
	repoMock := new(service.UserRepositoryMock)
	srv = service.NewGrpcUserService(repoMock, logger)

	testCases := []struct {
		testName string
		data     int
		res      bool
		err      error
	}{
		{
			testName: "delete user success",
			data:     1,
			res:      true,
			err:      nil,
		},
		{
			testName: "user does not exist error",
			data:     -1,
			res:      false,
			err:      errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			ctx := context.Background()
			assert := assert.New(t)

			// act
			repoMock.On("DeleteUser", ctx, tc.data).Return(tc.res, tc.err)
			res, err := srv.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
