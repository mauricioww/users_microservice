package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_srv/service"
)

// GrpcEndpoints stores the endpoinst for the current service
type GrpcEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

// MakeGrpcEndpoints returns a truct that stores the endpoints of the current service
func MakeGrpcEndpoints(srv service.GrpcUserServicer) GrpcEndpoints {
	return GrpcEndpoints{
		CreateUser:   makeCreateUserEndpoint(srv),
		Authenticate: makeAuthenticateEndpoint(srv),
		UpdateUser:   makeUpdateUserEndpoint(srv),
		GetUser:      makeGetUserEndpoint(srv),
		DeleteUser:   makeDeleteUserEndpoint(srv),
	}
}

func makeCreateUserEndpoint(srv service.GrpcUserServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		res, err := srv.CreateUser(ctx, req.Email, req.Password, req.Age)
		return CreateUserResponse{UserID: res}, err
	}
}

func makeAuthenticateEndpoint(srv service.GrpcUserServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		res, err := srv.Authenticate(ctx, req.Email, req.Password)
		return AuthenticateResponse{Success: res}, err
	}
}

func makeUpdateUserEndpoint(srv service.GrpcUserServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		res, err := srv.UpdateUser(ctx, req.UserID, req.Email, req.Password, req.Age)
		return UpdateUserResponse{Success: res}, err
	}
}

func makeGetUserEndpoint(srv service.GrpcUserServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		res, err := srv.GetUser(ctx, req.UserID)
		return GetUserResponse{Email: res.Email, Password: res.Password, Age: res.Age}, err
	}
}

func makeDeleteUserEndpoint(srv service.GrpcUserServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)
		res, err := srv.DeleteUser(ctx, req.UserID)
		return DeleteUserResponse{Success: res}, err
	}
}
