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
	srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srv)
	service := transport.NewGrpcUserDetailsServer(endpoints)

	testCases := []struct {
		testName string
		data     *detailspb.SetUserDetailsRequest
		res      *detailspb.SetUserDetailsResponse
		srvRes   bool
		err      error
	}{
		{
			testName: "set details success",
			data: &detailspb.SetUserDetailsRequest{
				UserId:       1,
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			srvRes: true,
			err:    nil,
		},
		{
			testName: "update details success",
			data: &detailspb.SetUserDetailsRequest{
				UserId:       1,
				MobileNumber: "12345789",
				Married:      false,
				Height:       1.75,
			},
			srvRes: true,
			err:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			// prepare
			assert := assert.New(t)
			ctx := context.Background()
			if tc.err != nil {
				tc.res = nil
			} else {
				tc.res = &detailspb.SetUserDetailsResponse{
					Success: tc.srvRes,
				}
			}

			// act
			srv.On("SetUserDetails", ctx, int(tc.data.GetUserId()), tc.data.GetCountry(), tc.data.GetCity(),
				tc.data.GetMobileNumber(), tc.data.GetMarried(), tc.data.GetHeight(), tc.data.GetWeight()).Return(tc.srvRes, tc.err)
			res, err := service.SetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestGetUserDetails(t *testing.T) {
	srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srv)
	service := transport.NewGrpcUserDetailsServer(endpoints)

	testCases := []struct {
		testName string
		data     *detailspb.GetUserDetailsRequest
		res      *detailspb.GetUserDetailsResponse
		srvRes   entities.UserDetails
		srvErr   error
		err      error
	}{
		{
			testName: "get details success",
			data: &detailspb.GetUserDetailsRequest{
				UserId: 0,
			},
			srvRes: entities.UserDetails{
				Country:      "Mexico",
				City:         "CDMX",
				MobileNumber: "11223344",
				Married:      false,
				Height:       1.75,
				Weight:       76.0,
			},
			srvErr: nil,
		},
		{
			testName: "get details which does not exist error",
			data: &detailspb.GetUserDetailsRequest{
				UserId: 1,
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
				tc.res = &detailspb.GetUserDetailsResponse{Country: tc.srvRes.Country, City: tc.srvRes.City, MobileNumber: tc.srvRes.MobileNumber,
					Married: tc.srvRes.Married, Height: tc.srvRes.Height, Weight: tc.srvRes.Weight}
			}

			// act
			srv.On("GetUserDetails", ctx, int(tc.data.GetUserId())).Return(tc.srvRes, tc.srvErr)
			res, err := service.GetUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}

func TestDeleteUserDetails(t *testing.T) {
	srv := new(transport.GrpcUserDetailsSrvMock)
	endpoints := transport.MakeGrpcEndpoints(srv)
	service := transport.NewGrpcUserDetailsServer(endpoints)

	testCases := []struct {
		testName string
		data     *detailspb.DeleteUserDetailsRequest
		res      *detailspb.DeleteUserDetailsResponse
		srvRes   bool
		srvErr   error
		err      error
	}{
		{
			testName: "delete details success",
			data: &detailspb.DeleteUserDetailsRequest{
				UserId: 0,
			},
			srvRes: true,
			srvErr: nil,
		},
		{
			testName: "delete details which does not exist error",
			data: &detailspb.DeleteUserDetailsRequest{
				UserId: 1,
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
				c, _ := tc.srvErr.(errors.ErrorResolver)
				tc.err = status.Error(c.GrpcCode(), tc.srvErr.Error())
			} else {
				tc.res = &detailspb.DeleteUserDetailsResponse{Success: tc.srvRes}
			}

			// act
			srv.On("DeleteUserDetails", ctx, int(tc.data.GetUserId())).Return(tc.srvRes, tc.srvErr)
			res, err := service.DeleteUserDetails(ctx, tc.data)

			// assert
			assert.Equal(tc.res, res)
			assert.Equal(tc.err, err)
		})
	}
}
