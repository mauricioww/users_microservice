package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/http_srv/service"
)

type HttpEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

func MakeHttpEndpoints(http_srv service.HttpService) HttpEndpoints {
	return HttpEndpoints{
		CreateUser:   makeCreateUserEndpoint(http_srv),
		Authenticate: makeAuthenticateEndpoint(http_srv),
		UpdateUser:   makeUpdateUserEndpoint(http_srv),
		GetUser:      makeGetUserEndpoint(http_srv),
		DeleteUser:   makeDeleteUserEndpont(http_srv),
	}
}

func makeCreateUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		res, err := http_srv.CreateUser(ctx, req.Email, req.Password, req.Age, req.Details)
		return CreateUserResponse{Id: res, Email: req.Email, Password: req.Password, Age: req.Age, Details: req.Details}, err
	}
}

func makeAuthenticateEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		res, err := http_srv.Authenticate(ctx, req.Email, req.Password)
		return AuthenticateResponse{Token: res}, err
	}
}

func makeUpdateUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)
		res, err := http_srv.UpdateUser(ctx, req.UserId, req.Email, req.Password, req.Age, req.Details)
		return UpdateUserResponse{Success: res}, err
	}
}

func makeGetUserEndpoint(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		res, err := http_srv.GetUser(ctx, req.UserId)
		return GetUserResponse{Id: req.UserId, Email: res.Email, Password: res.Password, Age: res.Age, Details: res.Details}, err
	}
}

func makeDeleteUserEndpont(http_srv service.HttpService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)
		res, err := http_srv.DeleteUser(ctx, req.UserId)
		return DeleteUserResponse{Success: res}, err
	}
}
