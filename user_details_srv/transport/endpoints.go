package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
)

// GrpcEndpoints stores the endpoints for the current service
type GrpcEndpoints struct {
	SetUserDetails    endpoint.Endpoint
	GetUserDetails    endpoint.Endpoint
	DeleteUserDetails endpoint.Endpoint
}

// MakeGrpcEndpoints returns a struct that stores the endpoints of the current service
func MakeGrpcEndpoints(srv service.GrpcUserDetailsServicer) GrpcEndpoints {
	return GrpcEndpoints{
		SetUserDetails:    makeSetUserDetailsEndpoint(srv),
		GetUserDetails:    makeGetUserDetailsEndpoint(srv),
		DeleteUserDetails: makeDeleteUserDetailsEndpoint(srv),
	}
}

func makeSetUserDetailsEndpoint(srv service.GrpcUserDetailsServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SetUserDetailsRequest)
		res, err := srv.SetUserDetails(ctx, req.UserID, req.Country, req.City, req.MobileNumber, req.Married, req.Height, req.Weigth)
		return SetUserDetailsResponse{Success: res}, err
	}
}

func makeGetUserDetailsEndpoint(srv service.GrpcUserDetailsServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserDetailsRequest)
		res, err := srv.GetUserDetails(ctx, req.UserID)
		return GetUserDetailsResponse{Country: res.Country, City: res.City, MobileNumber: res.MobileNumber, Married: res.Married, Height: res.Height, Weight: res.Weight}, err
	}
}

func makeDeleteUserDetailsEndpoint(srv service.GrpcUserDetailsServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserDetailsRequest)
		res, err := srv.DeleteUserDetails(ctx, req.UserID)
		return DeleteUserDetailsResponse{Success: res}, err
	}
}
