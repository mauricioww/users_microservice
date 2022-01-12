package transport

import (
	"context"
	"errors"

	grpc_gokit "github.com/go-kit/kit/transport/grpc"
	grpc_err "github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	setUserDetails    grpc_gokit.Handler
	getUserDetails    grpc_gokit.Handler
	deleteUserDetails grpc_gokit.Handler

	detailspb.UnimplementedUserDetailsServiceServer
}

func NewGrpcUserDetailsServer(endpoints GrpcUserDetailsServiceEndpoints) detailspb.UserDetailsServiceServer {
	return &gRPCServer{
		setUserDetails: grpc_gokit.NewServer(
			endpoints.SetUserDetails,
			decodeSetUserDetailsRequest,
			encodeSetUserDetailsResponse,
		),

		getUserDetails: grpc_gokit.NewServer(
			endpoints.GetUserDetails,
			decodeGetUserDetailsRequest,
			encodeGetUserDetailsResponse,
		),

		deleteUserDetails: grpc_gokit.NewServer(
			endpoints.DeleteUserDetails,
			decodeDeleteUserDetails,
			encodeDeleteUserDetails,
		),
	}
}

func decodeSetUserDetailsRequest(_ context.Context, request interface{}) (interface{}, error) {
	set_details, ok := request.(*detailspb.SetUserDetailsRequest)

	if !ok {
		return nil, errors.New("No proto message 'SetUserDetailsRequest'")
	}

	req := SetUserDetailsRequest{
		UserId:       int(set_details.GetUserId()),
		Country:      set_details.GetCountry(),
		City:         set_details.GetCity(),
		MobileNumber: set_details.GetMobileNumber(),
		Married:      set_details.GetMarried(),
		Height:       set_details.GetHeight(),
		Weigth:       set_details.GetWeight(),
	}

	return req, nil
}

func encodeSetUserDetailsResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(SetUserDetailsResponse)
	return &detailspb.SetUserDetailsResponse{Success: res.Success}, nil
}

func decodeGetUserDetailsRequest(_ context.Context, request interface{}) (interface{}, error) {
	get_details, ok := request.(*detailspb.GetUserDetailsRequest)

	if !ok {
		return nil, errors.New("No proto message 'GetUserDetailsRequest'")
	}

	req := GetUserDetailsRequest{
		UserId: int(get_details.GetUserId()),
	}

	return req, nil
}

func encodeGetUserDetailsResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(GetUserDetailsResponse)

	return &detailspb.GetUserDetailsResponse{Country: res.Country, City: res.City, MobileNumber: res.MobileNumber,
		Married: res.Married, Height: res.Height, Weight: res.Weight}, nil
}

func decodeDeleteUserDetails(_ context.Context, request interface{}) (interface{}, error) {
	delete_details, ok := request.(*detailspb.DeleteUserDetailsRequest)

	if !ok {
		return nil, errors.New("No proto message 'DeleteUserDetailsRequest'")
	}

	req := DeleteUserDetailsRequest{
		UserId: int(delete_details.GetUserId()),
	}

	return req, nil
}

func encodeDeleteUserDetails(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(DeleteUserDetailsResponse)
	return &detailspb.DeleteUserDetailsResponse{Success: res.Success}, nil
}

func (g *gRPCServer) SetUserDetails(ctx context.Context, req *detailspb.SetUserDetailsRequest) (*detailspb.SetUserDetailsResponse, error) {
	_, res, err := g.setUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		e, _ := err.(grpc_err.ErrorResolver)
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.SetUserDetailsResponse), nil
}

func (g *gRPCServer) GetUserDetails(ctx context.Context, req *detailspb.GetUserDetailsRequest) (*detailspb.GetUserDetailsResponse, error) {
	_, res, err := g.getUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		e, _ := err.(grpc_err.ErrorResolver)
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.GetUserDetailsResponse), nil
}

func (g *gRPCServer) DeleteUserDetails(ctx context.Context, req *detailspb.DeleteUserDetailsRequest) (*detailspb.DeleteUserDetailsResponse, error) {
	_, res, err := g.deleteUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		e, _ := err.(grpc_err.ErrorResolver)
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.DeleteUserDetailsResponse), nil
}
