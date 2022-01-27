package transport_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/transport"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestCreateUser(t *testing.T) {
	srvMock := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srvMock)
	grpcService := transport.NewGrpcUserServer(endpoints)

	testCases := []struct {
		testName string
		userReq  *userpb.CreateUserRequest
		userRes  *userpb.CreateUserResponse
		err      error
		srvRes   int
		srvErr   error
	}{
		{
			testName: "user created successfully",
			userReq: &userpb.CreateUserRequest{
				Email:    "success@email.com",
				Password: "qwerty",
				Age:      23,
			},
			srvRes: 0,
			srvErr: nil,
		},
		{
			testName: "no password error",
			userReq: &userpb.CreateUserRequest{
				Email: "success@email.com",
				Age:   23,
			},
			srvRes: -1,
			srvErr: errors.NewBadRequestPasswordError(),
		},
		{
			testName: "no email error",
			userReq: &userpb.CreateUserRequest{
				Password: "qwerty",
				Age:      23,
			},
			srvRes: -1,
			srvErr: errors.NewBadRequestEmailError(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srvErr != nil {
				tc.userRes = nil
				e, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srvErr.Error())
			} else {
				tc.userRes = &userpb.CreateUserResponse{
					Id:       int32(tc.srvRes),
					Email:    tc.userReq.GetEmail(),
					Password: tc.userReq.GetPassword(),
					Age:      tc.userReq.GetAge(),
				}
			}

			// act
			srvMock.On("CreateUser", ctx, tc.userReq.GetEmail(), tc.userReq.GetPassword(), int(tc.userReq.GetAge())).Return(tc.srvRes, tc.srvErr)
			res, err := grpcService.CreateUser(ctx, tc.userReq)

			// assert
			assert.Equal(tc.userRes, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	srvMock := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srvMock)
	grpcService := transport.NewGrpcUserServer(endpoints)

	testCases := []struct {
		testName string
		data     *userpb.AuthenticateRequest
		res      *userpb.AuthenticateResponse
		err      error
		srvRes   bool
		srvErr   error
	}{
		{
			testName: "authenticate successfully",
			data: &userpb.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "qwerty",
			},
			srvRes: true,
			srvErr: nil,
		},
		{
			testName: "no password error",
			data: &userpb.AuthenticateRequest{
				Email: "user@email.com",
			},
			srvErr: errors.NewBadRequestPasswordError(),
		},
		{
			testName: "no email error",
			data: &userpb.AuthenticateRequest{
				Password: "invalid_password",
			},
			srvErr: errors.NewBadRequestEmailError(),
		},
		{
			testName: "invalid password error",
			data: &userpb.AuthenticateRequest{
				Email:    "user@email.com",
				Password: "invalid_password",
			},
			srvErr: errors.NewUnauthenticatedError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srvErr != nil {
				tc.res = nil
				e, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srvErr.Error())

			} else {
				tc.res = &userpb.AuthenticateResponse{Success: tc.srvRes}
			}

			// act
			srvMock.On("Authenticate", ctx, tc.data.GetEmail(), tc.data.GetPassword()).Return(tc.srvRes, tc.srvErr)
			res, err := grpcService.Authenticate(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	srvMock := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srvMock)
	grpcService := transport.NewGrpcUserServer(endpoints)

	testCases := []struct {
		testName string
		data     *userpb.UpdateUserRequest
		res      *userpb.UpdateUserResponse
		err      error
		srvRes   bool
		srvErr   error
	}{
		{
			testName: "update user successfully",
			data: &userpb.UpdateUserRequest{
				Id:       1,
				Email:    "new_email@domain.com",
				Password: "new_password",
				Age:      25,
			},
			srvRes: true,
			srvErr: nil,
		},
		{
			testName: "no password error",
			data: &userpb.UpdateUserRequest{
				Id:    1,
				Email: "new_email@domain.com",
				Age:   25,
			},
			srvErr: errors.NewBadRequestPasswordError(),
		},
		{
			testName: "no email error",
			data: &userpb.UpdateUserRequest{
				Id:       1,
				Password: "new_password",
				Age:      25,
			},
			srvErr: errors.NewBadRequestEmailError(),
		},
		{
			testName: "user not found error",
			data: &userpb.UpdateUserRequest{
				Id:       1,
				Email:    "new_email@domain.com",
				Password: "new_password",
			},
			srvErr: errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srvErr != nil {
				tc.res = nil
				e, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srvErr.Error())
			} else {
				tc.res = &userpb.UpdateUserResponse{Success: tc.srvRes}
			}

			// act
			srvMock.On("UpdateUser", ctx, int(tc.data.GetId()), tc.data.GetEmail(), tc.data.GetPassword(), int(tc.data.GetAge())).Return(tc.srvRes, tc.srvErr)
			res, err := grpcService.UpdateUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	srvMock := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srvMock)
	grpcService := transport.NewGrpcUserServer(endpoints)

	testCases := []struct {
		testName string
		data     *userpb.GetUserRequest
		res      *userpb.GetUserResponse
		err      error
		srvRes   entities.User
		srvErr   error
	}{
		{
			testName: "user found",
			data: &userpb.GetUserRequest{
				Id: 0,
			},
			srvRes: entities.User{
				Email:    "user@email.com",
				Password: "password",
				Age:      20,
			},
			srvErr: nil,
		},
		{
			testName: "user not found error",
			data: &userpb.GetUserRequest{
				Id: 1,
			},
			srvErr: errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srvErr != nil {
				tc.res = nil
				e, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srvErr.Error())
			} else {
				tc.res = &userpb.GetUserResponse{
					Id:       tc.data.GetId(),
					Email:    tc.srvRes.Email,
					Password: tc.srvRes.Password,
					Age:      uint32(tc.srvRes.Age),
				}
			}

			// act
			srvMock.On("GetUser", ctx, int(tc.data.GetId())).Return(tc.srvRes, tc.srvErr)
			res, err := grpcService.GetUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	srvMock := new(transport.GrpcUserSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srvMock)
	grpcService := transport.NewGrpcUserServer(endpoints)

	testCases := []struct {
		testName string
		data     *userpb.DeleteUserRequest
		res      *userpb.DeleteUserResponse
		err      error
		srvRes   bool
		srvErr   error
	}{
		{
			testName: "delete user success",
			data: &userpb.DeleteUserRequest{
				Id: 0,
			},
			srvRes: true,
			srvErr: nil,
		},
		{
			testName: "user not found error",
			data: &userpb.DeleteUserRequest{
				Id: 1,
			},
			srvErr: errors.NewUserNotFoundError(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srvErr != nil {
				tc.res = nil
				e, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srvErr.Error())
			} else {
				tc.res = &userpb.DeleteUserResponse{Success: tc.srvRes}
			}

			// act
			srvMock.On("DeleteUser", ctx, int(tc.data.GetId())).Return(tc.srvRes, tc.srvErr)
			res, err := grpcService.DeleteUser(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
