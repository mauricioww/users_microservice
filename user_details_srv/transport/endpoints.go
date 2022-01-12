package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_details_srv/service"
)

type GrpcUserDetailsServiceEndpoints struct {
	SetUserDetails    endpoint.Endpoint
	GetUserDetails    endpoint.Endpoint
	DeleteUserDetails endpoint.Endpoint
}

func MakeGrpcUserDetailsServiceEndpoints(grpc_srv service.GrpcUserDetailsService) GrpcUserDetailsServiceEndpoints {
	return GrpcUserDetailsServiceEndpoints{
		SetUserDetails:    makeSetUserDetailsEndpoint(grpc_srv),
		GetUserDetails:    makeGetUserDetailsEndpoint(grpc_srv),
		DeleteUserDetails: makeDeleteUserDetailsEndpoint(grpc_srv),
	}
}

func makeSetUserDetailsEndpoint(srv service.GrpcUserDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(SetUserDetailsRequest)
		res, err := srv.SetUserDetails(ctx, req.UserId, req.Country, req.City, req.MobileNumber, req.Married, req.Height, req.Weigth)
		return SetUserDetailsResponse{Success: res}, err
	}
}

func makeGetUserDetailsEndpoint(srv service.GrpcUserDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(GetUserDetailsRequest)
		res, err := srv.GetUserDetails(ctx, req.UserId)
		return GetUserDetailsResponse{Country: res.Country, City: res.City, MobileNumber: res.MobileNumber, Married: res.Married, Height: res.Height, Weight: res.Weight}, err
	}
}

func makeDeleteUserDetailsEndpoint(srv service.GrpcUserDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(DeleteUserDetailsRequest)
		res, err := srv.DeleteUserDetails(ctx, req.UserId)
		return DeleteUserDetailsResponse{Success: res}, err
	}
}
