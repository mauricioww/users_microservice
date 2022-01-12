package transport_test

import (
	"context"
	"testing"

	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/transport"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestSetUserDetails(t *testing.T) {
	mock_srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserDetailsServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *detailspb.SetUserDetailsRequest
		res       *detailspb.SetUserDetailsResponse
		srv_res   bool
		err       error
	}{
		{
			test_name: "set details success",
			data: &detailspb.SetUserDetailsRequest{
				UserId:       1,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			srv_res: true,
			err:     nil,
		},
		{
			test_name: "update details success",
			data: &detailspb.SetUserDetailsRequest{
				UserId:       1,
				MobileNumber: "12345789",
				Married:      false,
				Height:       1.75,
			},
			srv_res: true,
			err:     nil,
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &detailspb.SetUserDetailsResponse{
					Success: tc.srv_res,
				}
			}

			// act
			mock_srv.On("SetUserDetails", ctx, int(tc.data.GetUserId()), tc.data.GetCountry(), tc.data.GetCity(),
				tc.data.GetMobileNumber(), tc.data.GetMarried(), tc.data.GetHeight(), tc.data.GetWeight()).Return(tc.srv_res, tc.err)
			res, err := grpc_service.SetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	mock_srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserDetailsServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *detailspb.GetUserDetailsRequest
		res       *detailspb.GetUserDetailsResponse
		srv_res   entities.UserDetails
		srv_err   error
		err       error
	}{
		{
			test_name: "get details success",
			data: &detailspb.GetUserDetailsRequest{
				UserId: 0,
			},
			srv_res: entities.UserDetails{
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			srv_err: nil,
		},
		{
			test_name: "get details which does not exist error",
			data: &detailspb.GetUserDetailsRequest{
				UserId: 1,
			},
			srv_err: errors.NewUserNotFoundError(),
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srv_err != nil {
				tc.res = nil
				e, _ := tc.srv_err.(errors.ErrorResolver)
				tc.err = status.Error(e.GrpcCode(), tc.srv_err.Error())
			} else {
				tc.res = &detailspb.GetUserDetailsResponse{Country: tc.srv_res.Country, City: tc.srv_res.City, MobileNumber: tc.srv_res.MobileNumber,
					Married: tc.srv_res.Married, Height: tc.srv_res.Height, Weight: tc.srv_res.Weight}
			}

			// act
			mock_srv.On("GetUserDetails", ctx, int(tc.data.GetUserId())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.GetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUserDetails(t *testing.T) {
	mock_srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcUserDetailsServiceEndpoints(mock_srv)
	grpc_service := transport.NewGrpcUserDetailsServer(endpoints)

	test_cases := []struct {
		test_name string
		data      *detailspb.DeleteUserDetailsRequest
		res       *detailspb.DeleteUserDetailsResponse
		srv_res   bool
		srv_err   error
		err       error
	}{
		{
			test_name: "delete details success",
			data: &detailspb.DeleteUserDetailsRequest{
				UserId: 0,
			},
			srv_res: true,
			srv_err: nil,
		},
		{
			test_name: "delete details which does not exist error",
			data: &detailspb.DeleteUserDetailsRequest{
				UserId: 1,
			},
			srv_err: errors.NewUserNotFoundError(),
		},
	}
	for _, tc := range test_cases {
		t.Run(tc.test_name, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.srv_err != nil {
				tc.res = nil
				c, _ := tc.srv_err.(errors.ErrorResolver)
				tc.err = status.Error(c.GrpcCode(), tc.srv_err.Error())
			} else {
				tc.res = &detailspb.DeleteUserDetailsResponse{Success: tc.srv_res}
			}

			// act
			mock_srv.On("DeleteUserDetails", ctx, int(tc.data.GetUserId())).Return(tc.srv_res, tc.err)
			res, err := grpc_service.DeleteUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
