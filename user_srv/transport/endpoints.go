package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/user_srv/service"
)

type GrpcUserServiceEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

func MakeGrpcUserServiceEndpoints(grpc_user_srv service.GrpcUserService) GrpcUserServiceEndpoints {
	return GrpcUserServiceEndpoints{
		CreateUser:   makeCreateUserEndpoint(grpc_user_srv),
		Authenticate: makeAuthenticateEndpoint(grpc_user_srv),
		UpdateUser:   makeUpdateUserEndpoint(grpc_user_srv),
		GetUser:      makeGetUserEndpoint(grpc_user_srv),
		DeleteUser:   makeDeleteUserEndpoint(grpc_user_srv),
	}
}

func makeCreateUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(CreateUserRequest)
		res, err := srv.CreateUser(ctx, req.Email, req.Password, req.Age)
		return CreateUserResponse{Id: res}, err
	}
}

func makeAuthenticateEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(AuthenticateRequest)
		res, err := srv.Authenticate(ctx, req.Email, req.Password)
		return AuthenticateResponse{Id: res}, err
	}
}

func makeUpdateUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, _ := request.(UpdateUserRequest)
		res, err := srv.UpdateUser(ctx, req.Id, req.Email, req.Password, req.Age)
		return UpdateUserResponse{Success: res}, err
	}
}

func makeGetUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(GetUserRequest)
		res, err := srv.GetUser(ctx, req.UserId)
		return GetUserResponse{Email: res.Email, Password: res.Password, Age: res.Age}, err
	}
}

func makeDeleteUserEndpoint(srv service.GrpcUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(DeleteUserRequest)
		res, err := srv.DeleteUser(ctx, req.UserId)
		return DeleteUserResponse{Success: res}, err
	}
}
