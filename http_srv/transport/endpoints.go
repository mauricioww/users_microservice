package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mauricioww/user_microsrv/http_srv/service"
)

// HTTPEndpoints stores the endpoins of the current service
type HTTPEndpoints struct {
	CreateUser   endpoint.Endpoint
	Authenticate endpoint.Endpoint
	UpdateUser   endpoint.Endpoint
	GetUser      endpoint.Endpoint
	DeleteUser   endpoint.Endpoint
}

// MakeHTTPEndpoints build the custom endpoints for the http service
func MakeHTTPEndpoints(httpSrv service.HTTPServicer) HTTPEndpoints {
	return HTTPEndpoints{
		CreateUser:   makeCreateUserEndpoint(httpSrv),
		Authenticate: makeAuthenticateEndpoint(httpSrv),
		UpdateUser:   makeUpdateUserEndpoint(httpSrv),
		GetUser:      makeGetUserEndpoint(httpSrv),
		DeleteUser:   makeDeleteUserEndpont(httpSrv),
	}
}

func makeCreateUserEndpoint(httpSrv service.HTTPServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		res, err := httpSrv.CreateUser(ctx, req.Email, req.Password, req.Age, req.Details)
		return CreateUserResponse{UserID: res, Email: req.Email, Password: req.Password, Age: req.Age, Details: req.Details}, err
	}
}

func makeAuthenticateEndpoint(httpSrv service.HTTPServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthenticateRequest)
		res, err := httpSrv.Authenticate(ctx, req.Email, req.Password)
		return AuthenticateResponse{Success: res}, err
	}
}

func makeUpdateUserEndpoint(httpSrv service.HTTPServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUserRequest)
		res, err := httpSrv.UpdateUser(ctx, req.UserID, req.Email, req.Password, req.Age, req.Details)
		return UpdateUserResponse{Success: res}, err
	}
}

func makeGetUserEndpoint(httpSrv service.HTTPServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetUserRequest)
		res, err := httpSrv.GetUser(ctx, req.UserID)
		return GetUserResponse{UserID: req.UserID, Email: res.Email, Password: res.Password, Age: res.Age, Details: res.Details}, err
	}
}

func makeDeleteUserEndpont(httpSrv service.HTTPServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteUserRequest)
		res, err := httpSrv.DeleteUser(ctx, req.UserID)
		return DeleteUserResponse{Success: res}, err
	}
}
