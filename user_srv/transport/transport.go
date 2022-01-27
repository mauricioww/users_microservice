package transport

import (
	"context"
	"errors"

	grpcGokit "github.com/go-kit/kit/transport/grpc"
	grpcError "github.com/mauricioww/user_microsrv/errors"
	"google.golang.org/grpc/status"

	"github.com/mauricioww/user_microsrv/user_srv/userpb"
)

type gRPCServer struct {
	createUser   grpcGokit.Handler
	authenticate grpcGokit.Handler
	updateUser   grpcGokit.Handler
	getUser      grpcGokit.Handler
	deleteUser   grpcGokit.Handler

	userpb.UnimplementedUserServiceServer
}

// NewGrpcUserServer returns the server with the endpoints and the specifications for each one
func NewGrpcUserServer(endpoints GrpcEndpoints) userpb.UserServiceServer {
	return &gRPCServer{
		createUser: grpcGokit.NewServer(
			endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),

		authenticate: grpcGokit.NewServer(
			endpoints.Authenticate,
			decodeAuthenticateRequest,
			encodeAuthenticateResponse,
		),

		updateUser: grpcGokit.NewServer(
			endpoints.UpdateUser,
			decodeUpdateUserRequest,
			encondeUpdateUserResponse,
		),

		getUser: grpcGokit.NewServer(
			endpoints.GetUser,
			decodeGetUserRequest,
			encodeGetUserResponse,
		),

		deleteUser: grpcGokit.NewServer(
			endpoints.DeleteUser,
			decodeDeleteUserRequest,
			encondeDeleteUserResponse,
		),
	}
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	createPb, ok := request.(*userpb.CreateUserRequest)

	if !ok {
		return nil, errors.New("no proto message 'CreateUserRequest'")
	}

	req := CreateUserRequest{
		Email:    createPb.GetEmail(),
		Password: createPb.GetPassword(),
		Age:      int(createPb.GetAge()),
	}

	return req, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(CreateUserResponse)
	return &userpb.CreateUserResponse{Id: int32(res.UserID), Email: res.Email, Password: res.Password, Age: uint32(res.Age)}, nil
}

func decodeAuthenticateRequest(_ context.Context, request interface{}) (interface{}, error) {
	authPb, ok := request.(*userpb.AuthenticateRequest)

	if !ok {
		return nil, errors.New("no proto message 'AuthenticateRequest'")
	}

	req := AuthenticateRequest{
		Email:    authPb.GetEmail(),
		Password: authPb.GetPassword(),
	}

	return req, nil
}

func encodeAuthenticateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(AuthenticateResponse)
	return &userpb.AuthenticateResponse{Success: res.Success}, nil
}

func decodeUpdateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	updatePb, ok := request.(*userpb.UpdateUserRequest)

	if !ok {
		return nil, errors.New("no proto message 'UpdateUserRequest'")
	}

	req := UpdateUserRequest{
		UserID:   int(updatePb.GetId()),
		Email:    updatePb.GetEmail(),
		Password: updatePb.GetPassword(),
		Age:      int(updatePb.GetAge()),
	}

	return req, nil
}

func encondeUpdateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(UpdateUserResponse)
	return &userpb.UpdateUserResponse{Success: res.Success}, nil
}

func decodeGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	getPb, ok := request.(*userpb.GetUserRequest)

	if !ok {
		return nil, errors.New("no proto message 'GetUserRequest'")
	}

	req := GetUserRequest{
		UserID: int(getPb.GetId()),
	}

	return req, nil
}

func encodeGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(GetUserResponse)
	return &userpb.GetUserResponse{Id: uint32(res.UserID), Email: res.Email, Password: res.Password, Age: uint32(res.Age)}, nil
}

func decodeDeleteUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	deletePb, ok := request.(*userpb.DeleteUserRequest)

	if !ok {
		return nil, errors.New("no proto message 'DeleteUserRequest'")
	}

	req := DeleteUserRequest{
		UserID: int(deletePb.GetId()),
	}

	return req, nil
}

func encondeDeleteUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(DeleteUserResponse)
	return &userpb.DeleteUserResponse{Success: res.Success}, nil
}

func (g *gRPCServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	_, res, err := g.createUser.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*userpb.CreateUserResponse), err
}

func (g *gRPCServer) Authenticate(ctx context.Context, req *userpb.AuthenticateRequest) (*userpb.AuthenticateResponse, error) {
	_, res, err := g.authenticate.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*userpb.AuthenticateResponse), nil
}

func (g *gRPCServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	_, res, err := g.updateUser.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*userpb.UpdateUserResponse), nil
}

func (g *gRPCServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	_, res, err := g.getUser.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*userpb.GetUserResponse), nil
}

func (g *gRPCServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	_, res, err := g.deleteUser.ServeGRPC(ctx, req)

	if err != nil {
		e, ok := err.(grpcError.ErrorResolver)
		if !ok {
			u := grpcError.NewUnknownError()
			return nil, status.Error(u.GrpcCode(), u.Error())
		}
		return nil, status.Error(e.GrpcCode(), err.Error())
	}

	return res.(*userpb.DeleteUserResponse), nil
}
