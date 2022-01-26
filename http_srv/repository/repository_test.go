package repository_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	userMock := new(repository.GrpcUserMock)
	detailsMock := new(repository.GrpcDetailsMock)
	conn1, conn2, httpRepository := repository.InitRepoMock(userMock, detailsMock)

	defer conn1.Close()
	defer conn2.Close()

	testCases := []struct {
		testName      string
		user          entities.User
		repositoryRes int
		repositoryErr error
		userErr       error
	}{
		{
			testName: "user created successfully",
			user: entities.User{
				Email:    "user@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenerateDetails(),
			},
			repositoryRes: 1,
			userErr:       nil,
		},
		{
			testName: "no password error",
			user: entities.User{
				Email:   "user@email.com",
				Age:     23,
				Details: repository.GenerateDetails(),
			},
			repositoryRes: -1,
			userErr:       status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			testName: "no email error",
			user: entities.User{
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenerateDetails(),
			},
			repositoryRes: -1,
			userErr:       status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var userRes *userpb.CreateUserResponse
			var detailsRes *detailspb.SetUserDetailsResponse
			userReq := &userpb.CreateUserRequest{
				Email:    tc.user.Email,
				Password: tc.user.Password,
				Age:      uint32(tc.user.Age),
			}
			detailsReq := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.repositoryRes),
				Country:      tc.user.Country,
				City:         tc.user.City,
				MobileNumber: tc.user.MobileNumber,
				Married:      tc.user.Married,
				Height:       tc.user.Height,
				Weight:       tc.user.Weight,
			}

			if tc.userErr == nil {
				userRes = &userpb.CreateUserResponse{Id: int32(tc.repositoryRes)}
				detailsRes = &detailspb.SetUserDetailsResponse{Success: tc.repositoryRes >= 0}
			}

			userMock.On("CreateUser", mock.Anything, userReq).Return(userRes, tc.userErr)
			detailsMock.On("SetUserDetails", mock.Anything, detailsReq).Return(detailsRes, nil)

			// act
			res, err := httpRepository.CreateUser(ctx, tc.user)

			// assert
			assert.Equal(res, tc.repositoryRes)
			assert.True(repository.TestErrors(err, tc.userErr))
		})
	}
}

func TestAuthenticate(t *testing.T) {
	userMock := new(repository.GrpcUserMock)
	detailsMock := new(repository.GrpcDetailsMock)
	conn1, conn2, httpRepository := repository.InitRepoMock(userMock, detailsMock)

	defer conn1.Close()
	defer conn2.Close()

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
			testName: "user not found",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			err: status.Error(codes.NotFound, "User not found"),
		},
		{
			testName: "invalid password or email",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "password",
			},
			err: status.Error(codes.Unauthenticated, "Password or email error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var userRes *userpb.AuthenticateResponse
			userReq := &userpb.AuthenticateRequest{
				Email:    tc.data.Email,
				Password: tc.data.Password,
			}
			if tc.err == nil {
				userRes = &userpb.AuthenticateResponse{Success: tc.res}
			}

			// act
			userMock.On("Authenticate", mock.Anything, userReq).Return(userRes, tc.err)
			res, err := httpRepository.Authenticate(ctx, tc.data)

			// assert
			assert.True(repository.TestErrors(err, tc.err))
			assert.Equal(res, tc.res)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	userMock := new(repository.GrpcUserMock)
	detailsMock := new(repository.GrpcDetailsMock)
	conn1, conn2, httpRepository := repository.InitRepoMock(userMock, detailsMock)

	defer conn1.Close()
	defer conn2.Close()

	testCases := []struct {
		testName string
		data     entities.UserUpdate
		res      bool
		err      error
	}{
		{
			testName: "Update email successfully",
			data: entities.UserUpdate{
				UserID: 0,
				User: entities.User{
					Email:    "email@domian.com",
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenerateDetails(),
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
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenerateDetails(),
				},
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
		{
			testName: "no password error",
			data: entities.UserUpdate{
				UserID: 0,
				User: entities.User{
					Email:   "email@domian.com",
					Age:     23,
					Details: repository.GenerateDetails(),
				},
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			testName: "user not found error",
			data: entities.UserUpdate{
				UserID: 1,
				User: entities.User{
					Email:    "email@domian.com",
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenerateDetails(),
				},
			},
			err: status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)
			var userRes *userpb.UpdateUserResponse
			var detailsRes *detailspb.SetUserDetailsResponse
			userReq := &userpb.UpdateUserRequest{
				Id:       uint32(tc.data.UserID),
				Email:    tc.data.Email,
				Password: tc.data.Password,
				Age:      uint32(tc.data.Age),
			}
			detailsReq := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.data.UserID),
				Country:      tc.data.Country,
				City:         tc.data.City,
				MobileNumber: tc.data.MobileNumber,
				Married:      tc.data.Married,
				Height:       tc.data.Height,
				Weight:       tc.data.Weight,
			}
			if tc.err == nil {
				userRes = &userpb.UpdateUserResponse{Success: tc.res}
				detailsRes = &detailspb.SetUserDetailsResponse{Success: tc.res}
			}

			// act
			userMock.On("UpdateUser", mock.Anything, userReq).Return(userRes, tc.err)
			detailsMock.On("SetUserDetails", mock.Anything, detailsReq).Return(detailsRes, nil)
			res, err := httpRepository.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}

func TestGetUser(t *testing.T) {
	userMock := new(repository.GrpcUserMock)
	detailsMock := new(repository.GrpcDetailsMock)
	conn1, conn2, httpRepository := repository.InitRepoMock(userMock, detailsMock)

	defer conn1.Close()
	defer conn2.Close()

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
				Email:    "email@domain.com",
				Password: "password",
				Age:      10,
			},
			err: nil,
		},
		{
			testName: "user not found",
			data:     1,
			res:      entities.User{},
			err:      status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var userRes *userpb.GetUserResponse
			var detailsRes *detailspb.GetUserDetailsResponse
			userReq := &userpb.GetUserRequest{
				Id: uint32(tc.data),
			}
			detailsReq := &detailspb.GetUserDetailsRequest{
				UserId: uint32(tc.data),
			}
			if tc.err == nil {
				userRes = &userpb.GetUserResponse{
					Email:    tc.res.Email,
					Password: tc.res.Password,
					Age:      uint32(tc.res.Age),
				}
				detailsRes = &detailspb.GetUserDetailsResponse{}
			}

			// act
			userMock.On("GetUser", mock.Anything, userReq).Return(userRes, tc.err)
			detailsMock.On("GetUserDetails", mock.Anything, detailsReq).Return(detailsRes, tc.err)
			res, err := httpRepository.GetUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}

func TestDeleteUser(t *testing.T) {
	userMock := new(repository.GrpcUserMock)
	detailsMock := new(repository.GrpcDetailsMock)
	conn1, conn2, httpRepository := repository.InitRepoMock(userMock, detailsMock)

	defer conn1.Close()
	defer conn2.Close()

	testCases := []struct {
		testName string
		data     int
		res      bool
		err      error
	}{
		{
			testName: "user deleted success",
			data:     0,
			res:      true,
			err:      nil,
		},
		{
			testName: "user delete error",
			data:     1,
			res:      false,
			err:      status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var userRes *userpb.DeleteUserResponse
			var detailsRes *detailspb.DeleteUserDetailsResponse
			userReq := &userpb.DeleteUserRequest{Id: uint32(tc.data)}
			detailsReq := &detailspb.DeleteUserDetailsRequest{UserId: uint32(tc.data)}
			if tc.err == nil {
				userRes = &userpb.DeleteUserResponse{Success: tc.res}
				detailsRes = &detailspb.DeleteUserDetailsResponse{Success: tc.res}
			}

			// act
			userMock.On("DeleteUser", mock.Anything, userReq).Return(userRes, tc.err)
			detailsMock.On("DeleteUserDetails", mock.Anything, detailsReq).Return(detailsRes, nil)
			res, err := httpRepository.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}
