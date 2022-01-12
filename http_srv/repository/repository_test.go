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
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

	test_cases := []struct {
		test_name      string
		user           entities.User
		repository_res int
		repository_err error
		user_err       error
	}{
		{
			test_name: "user created successfully",
			user: entities.User{
				Email:    "user@email.com",
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenereateDetails(),
			},
			repository_res: 1,
			user_err:       nil,
		},
		{
			test_name: "no password error",
			user: entities.User{
				Email:   "user@email.com",
				Age:     23,
				Details: repository.GenereateDetails(),
			},
			repository_res: -1,
			user_err:       status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			test_name: "no email error",
			user: entities.User{
				Password: "qwerty",
				Age:      23,
				Details:  repository.GenereateDetails(),
			},
			repository_res: -1,
			user_err:       status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var user_res *userpb.CreateUserResponse
			var details_res *detailspb.SetUserDetailsResponse
			user_req := &userpb.CreateUserRequest{
				Email:    tc.user.Email,
				Password: tc.user.Password,
				Age:      uint32(tc.user.Age),
			}
			details_req := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.repository_res),
				Country:      tc.user.Country,
				City:         tc.user.City,
				MobileNumber: tc.user.MobileNumber,
				Married:      tc.user.Married,
				Height:       tc.user.Height,
				Weight:       tc.user.Weight,
			}

			if tc.user_err == nil {
				user_res = &userpb.CreateUserResponse{Id: int32(tc.repository_res)}
				details_res = &detailspb.SetUserDetailsResponse{Success: tc.repository_res >= 0}
			}

			user_mock.On("CreateUser", mock.Anything, user_req).Return(user_res, tc.user_err)
			details_mock.On("SetUserDetails", mock.Anything, details_req).Return(details_res, nil)

			// act
			res, err := http_repository.CreateUser(ctx, tc.user)

			// assert
			assert.Equal(res, tc.repository_res)
			assert.True(repository.TestErrors(err, tc.user_err))
		})
	}
}

func TestAuthenticate(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

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
			res: 1,
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
			test_name: "user not found",
			data: entities.Session{
				Email:    "no_real@email.com",
				Password: "fake_password",
			},
			res: -1,
			err: status.Error(codes.NotFound, "User not found"),
		},
		{
			test_name: "invalid password or email",
			data: entities.Session{
				Email:    "user@email.com",
				Password: "password",
			},
			res: -1,
			err: status.Error(codes.Unauthenticated, "Password or email error"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var grpc_res *userpb.AuthenticateResponse
			grpc_req := &userpb.AuthenticateRequest{
				Email:    tc.data.Email,
				Password: tc.data.Password,
			}
			if tc.err == nil {
				grpc_res = &userpb.AuthenticateResponse{UserId: int32(tc.res)}
			}

			// act
			user_mock.On("Authenticate", mock.Anything, grpc_req).Return(grpc_res, tc.err)
			res, err := http_repository.Authenticate(ctx, tc.data)

			// assert
			assert.True(repository.TestErrors(err, tc.err))
			assert.Equal(res, tc.res)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

	test_cases := []struct {
		test_name string
		data      entities.UserUpdate
		res       bool
		err       error
	}{
		{
			test_name: "Update email successfully",
			data: entities.UserUpdate{
				UserId: 0,
				User: entities.User{
					Email:    "email@domian.com",
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenereateDetails(),
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
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenereateDetails(),
				},
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'email'"),
		},
		{
			test_name: "no password error",
			data: entities.UserUpdate{
				UserId: 0,
				User: entities.User{
					Email:   "email@domian.com",
					Age:     23,
					Details: repository.GenereateDetails(),
				},
			},
			err: status.Error(codes.FailedPrecondition, "Missing field 'password'"),
		},
		{
			test_name: "user not found error",
			data: entities.UserUpdate{
				UserId: 1,
				User: entities.User{
					Email:    "email@domian.com",
					Password: "qwerty",
					Age:      23,
					Details:  repository.GenereateDetails(),
				},
			},
			err: status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			ctx := context.Background()
			assert := assert.New(t)
			var user_res *userpb.UpdateUserResponse
			var details_res *detailspb.SetUserDetailsResponse
			user_req := &userpb.UpdateUserRequest{
				Id:       uint32(tc.data.UserId),
				Email:    tc.data.Email,
				Password: tc.data.Password,
				Age:      uint32(tc.data.Age),
			}
			details_req := &detailspb.SetUserDetailsRequest{
				UserId:       uint32(tc.data.UserId),
				Country:      tc.data.Country,
				City:         tc.data.City,
				MobileNumber: tc.data.MobileNumber,
				Married:      tc.data.Married,
				Height:       tc.data.Height,
				Weight:       tc.data.Weight,
			}
			if tc.err == nil {
				user_res = &userpb.UpdateUserResponse{Success: tc.res}
				details_res = &detailspb.SetUserDetailsResponse{Success: tc.res}
			}

			// act
			user_mock.On("UpdateUser", mock.Anything, user_req).Return(user_res, tc.err)
			details_mock.On("SetUserDetails", mock.Anything, details_req).Return(details_res, nil)
			res, err := http_repository.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}

func TestGetUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

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
				Email:    "email@domain.com",
				Password: "password",
				Age:      10,
			},
			err: nil,
		},
		{
			test_name: "user not found",
			data:      1,
			res:       entities.User{},
			err:       status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var user_res *userpb.GetUserResponse
			var details_res *detailspb.GetUserDetailsResponse
			user_req := &userpb.GetUserRequest{
				Id: uint32(tc.data),
			}
			details_req := &detailspb.GetUserDetailsRequest{
				UserId: uint32(tc.data),
			}
			if tc.err == nil {
				user_res = &userpb.GetUserResponse{
					Email:    tc.res.Email,
					Password: tc.res.Password,
					Age:      uint32(tc.res.Age),
				}
				details_res = &detailspb.GetUserDetailsResponse{}
			}

			// act
			user_mock.On("GetUser", mock.Anything, user_req).Return(user_res, tc.err)
			details_mock.On("GetUserDetails", mock.Anything, details_req).Return(details_res, tc.err)
			res, err := http_repository.GetUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}

func TestDeleteUser(t *testing.T) {
	user_mock := new(repository.GrpcUserMock)
	details_mock := new(repository.GrpcDetailsMock)
	conn1, conn2, http_repository := repository.InitRepoMock(user_mock, details_mock)

	defer conn1.Close()
	defer conn2.Close()

	test_cases := []struct {
		test_name string
		data      int
		res       bool
		err       error
	}{
		{
			test_name: "user deleted success",
			data:      0,
			res:       true,
			err:       nil,
		},
		{
			test_name: "user delete error",
			data:      1,
			res:       false,
			err:       status.Error(codes.NotFound, "User not found"),
		},
	}

	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			//  prepare
			assert := assert.New(t)
			ctx := context.Background()
			var user_res *userpb.DeleteUserResponse
			var details_res *detailspb.DeleteUserDetailsResponse
			user_req := &userpb.DeleteUserRequest{Id: uint32(tc.data)}
			details_req := &detailspb.DeleteUserDetailsRequest{UserId: uint32(tc.data)}
			if tc.err == nil {
				user_res = &userpb.DeleteUserResponse{Success: tc.res}
				details_res = &detailspb.DeleteUserDetailsResponse{Success: tc.res}
			}

			// act
			user_mock.On("DeleteUser", mock.Anything, user_req).Return(user_res, tc.err)
			details_mock.On("DeleteUserDetails", mock.Anything, details_req).Return(details_res, nil)
			res, err := http_repository.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(res, tc.res)
			assert.True(repository.TestErrors(err, tc.err))
		})
	}
}
