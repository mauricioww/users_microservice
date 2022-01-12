package repository

import (
	"context"
	l "log"
	"net"
	"os"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type GrpcUserMock struct {
	mock.Mock
	userpb.UnimplementedUserServiceServer
}

type GrpcDetailsMock struct {
	mock.Mock
	detailspb.UnimplementedUserDetailsServiceServer
}

func InitRepoMock(mock1 *GrpcUserMock, mock2 *GrpcDetailsMock) (*grpc.ClientConn, *grpc.ClientConn, HttpRepository) {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(
			logger,
			"service",
			"account",
			"time",
			log.DefaultTimestampUTC,
			"caller",
			log.DefaultCaller,
		)
	}

	conn1, _ := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(Dialer1(mock1)))
	conn2, _ := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(Dialer2(mock2)))

	r := NewHttpRepository(conn1, conn2, logger)
	return conn1, conn2, r
}

func Dialer1(m *GrpcUserMock) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, m)

	go func() {
		if err := server.Serve(listener); err != nil {
			l.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func Dialer2(m *GrpcDetailsMock) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	detailspb.RegisterUserDetailsServiceServer(server, m)

	go func() {
		if err := server.Serve(listener); err != nil {
			l.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func (m *GrpcUserMock) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.CreateUserResponse), args.Error(1)
}

func (m *GrpcUserMock) Authenticate(ctx context.Context, req *userpb.AuthenticateRequest) (*userpb.AuthenticateResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.AuthenticateResponse), args.Error(1)
}

func (m *GrpcUserMock) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.UpdateUserResponse), args.Error(1)
}

func (m *GrpcUserMock) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.GetUserResponse), args.Error(1)
}

func (m *GrpcUserMock) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*userpb.DeleteUserResponse), args.Error(1)
}

func (m *GrpcDetailsMock) SetUserDetails(ctx context.Context, req *detailspb.SetUserDetailsRequest) (*detailspb.SetUserDetailsResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*detailspb.SetUserDetailsResponse), args.Error(1)
}

func (m *GrpcDetailsMock) GetUserDetails(ctx context.Context, req *detailspb.GetUserDetailsRequest) (*detailspb.GetUserDetailsResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*detailspb.GetUserDetailsResponse), args.Error(1)
}

func (m *GrpcDetailsMock) DeleteUserDetails(ctx context.Context, req *detailspb.DeleteUserDetailsRequest) (*detailspb.DeleteUserDetailsResponse, error) {
	args := m.Called(ctx, req)

	return args.Get(0).(*detailspb.DeleteUserDetailsResponse), args.Error(1)
}

func GenereateDetails() entities.Details {
	return entities.Details{
		Country:      "Mexico",
		City:         "CDMX",
		MobileNumber: "11223344",
		Married:      false,
		Height:       1.75,
		Weight:       76.0,
	}
}

func TestErrors(err1 error, err2 error) bool {
	e1 := status.Convert(err1)
	e2 := status.Convert(err2)
	return (e1.Code() == e2.Code()) && (e1.Message() == e2.Message())
}
