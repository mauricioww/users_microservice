package transport

import (
	"context"
	"errors"

	grpcGokit "github.com/go-kit/kit/transport/grpc"
	grpcError "github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"google.golang.org/grpc/status"
)

type gRPCServer struct {
	setUserDetails    grpcGokit.Handler
	getUserDetails    grpcGokit.Handler
	deleteUserDetails grpcGokit.Handler

	detailspb.UnimplementedUserDetailsServiceServer
}

// NewGrpcUserDetailsServer returns the server with the endpoints and the specifications for each one
func NewGrpcUserDetailsServer(endpoints GrpcEndpoints) detailspb.UserDetailsServiceServer {
	return &gRPCServer{
		setUserDetails: grpcGokit.NewServer(
			endpoints.SetUserDetails,
			decodeSetUserDetailsRequest,
			encodeSetUserDetailsResponse,
		),

		getUserDetails: grpcGokit.NewServer(
			endpoints.GetUserDetails,
			decodeGetUserDetailsRequest,
			encodeGetUserDetailsResponse,
		),

		deleteUserDetails: grpcGokit.NewServer(
			endpoints.DeleteUserDetails,
			decodeDeleteUserDetails,
			encodeDeleteUserDetails,
		),
	}
}

func decodeSetUserDetailsRequest(_ context.Context, request interface{}) (interface{}, error) {
	setDetails, ok := request.(*detailspb.SetUserDetailsRequest)

	if !ok {
		return nil, errors.New("no proto message 'SetUserDetailsRequest'")
	}

	req := SetUserDetailsRequest{
		UserID:       int(setDetails.GetUserId()),
		Country:      setDetails.GetCountry(),
		City:         setDetails.GetCity(),
		MobileNumber: setDetails.GetMobileNumber(),
		Married:      setDetails.GetMarried(),
		Height:       setDetails.GetHeight(),
		Weigth:       setDetails.GetWeight(),
	}

	return req, nil
}

func encodeSetUserDetailsResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(SetUserDetailsResponse)
	return &detailspb.SetUserDetailsResponse{Success: res.Success}, nil
}

func decodeGetUserDetailsRequest(_ context.Context, request interface{}) (interface{}, error) {
	getDetails, ok := request.(*detailspb.GetUserDetailsRequest)

	if !ok {
		return nil, errors.New("no proto message 'GetUserDetailsRequest'")
	}

	req := GetUserDetailsRequest{
		UserID: int(getDetails.GetUserId()),
	}

	return req, nil
}

func encodeGetUserDetailsResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(GetUserDetailsResponse)
	return &detailspb.GetUserDetailsResponse{Country: res.Country, City: res.City, MobileNumber: res.MobileNumber,
		Married: res.Married, Height: res.Height, Weight: res.Weight}, nil
}

func decodeDeleteUserDetails(_ context.Context, request interface{}) (interface{}, error) {
	deleteDetails, ok := request.(*detailspb.DeleteUserDetailsRequest)

	if !ok {
		return nil, errors.New("no proto message 'DeleteUserDetailsRequest'")
	}

	req := DeleteUserDetailsRequest{
		UserID: int(deleteDetails.GetUserId()),
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
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.SetUserDetailsResponse), nil
}

func (g *gRPCServer) GetUserDetails(ctx context.Context, req *detailspb.GetUserDetailsRequest) (*detailspb.GetUserDetailsResponse, error) {
	_, res, err := g.getUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.GetUserDetailsResponse), nil
}

func (g *gRPCServer) DeleteUserDetails(ctx context.Context, req *detailspb.DeleteUserDetailsRequest) (*detailspb.DeleteUserDetailsResponse, error) {
	_, res, err := g.deleteUserDetails.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*detailspb.DeleteUserDetailsResponse), nil
}
