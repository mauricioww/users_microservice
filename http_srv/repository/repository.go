package repository

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

// HTTPRepositorier describes the necessary methods to send requests to both gRPC servers
type HTTPRepositorier interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	Authenticate(ctx context.Context, session entities.Session) (bool, error)
	UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

// HTTPRepository type implement the HTTPRepositorier interface
type HTTPRepository struct {
	userClient    userpb.UserServiceClient
	detailsClient detailspb.UserDetailsServiceClient
	logger        log.Logger
}

// NewHTTPRepository returns an HTTPRepository tye pointer
func NewHTTPRepository(conn1 *grpc.ClientConn, conn2 *grpc.ClientConn, logger log.Logger) *HTTPRepository {
	return &HTTPRepository{
		userClient:    userpb.NewUserServiceClient(conn1),
		detailsClient: detailspb.NewUserDetailsServiceClient(conn2),
		logger:        log.With(logger, "http_service", "repository"),
	}
}

// CreateUser sends the original data to both gRPC servers
func (r *HTTPRepository) CreateUser(ctx context.Context, user entities.User) (int, error) {
	logger := log.With(r.logger, "method", "create_users")

	userReq := userpb.CreateUserRequest{
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}

	userRes, err := r.userClient.CreateUser(ctx, &userReq)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return -1, err
	}

	detailsReq := detailspb.SetUserDetailsRequest{
		UserId:       uint32(userRes.GetId()),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weight:       user.Details.Weight,
	}

	detailsRes, err := r.detailsClient.SetUserDetails(ctx, &detailsReq)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return -1, err
	}

	if detailsRes.GetSuccess() {
		return int(userRes.GetId()), nil
	}

	return -1, nil
}

// Authenticate sends the credentials to gRPC server in order to do a login
func (r *HTTPRepository) Authenticate(ctx context.Context, session entities.Session) (bool, error) {
	logger := log.With(r.logger, "method", "authenticate_user")

	userReq := userpb.AuthenticateRequest{
		Email:    session.Email,
		Password: session.Password,
	}
	userRes, err := r.userClient.Authenticate(ctx, &userReq)

	if err != nil {
		level.Error(logger).Log("err_user", err)
	}

	return userRes.GetSuccess(), err
}

// UpdateUser sends the original request to both gRPC servers
func (r *HTTPRepository) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	logger := log.With(r.logger, "method", "update_user")

	userReq := userpb.UpdateUserRequest{
		Id:       uint32(user.UserID),
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}
	detailsReq := detailspb.SetUserDetailsRequest{
		UserId:       uint32(user.UserID),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weight:       user.Details.Weight,
	}

	userRes, err := r.userClient.UpdateUser(ctx, &userReq)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return false, err
	}

	detailsRes, err := r.detailsClient.SetUserDetails(ctx, &detailsReq)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return false, err
	}

	return detailsRes.GetSuccess() && userRes.GetSuccess(), nil
}

// GetUser fetchs the information of a user from both gRPC servers
func (r *HTTPRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(r.logger, "method", "get_user")

	userReq := userpb.GetUserRequest{
		Id: uint32(id),
	}
	detailsReq := detailspb.GetUserDetailsRequest{
		UserId: uint32(id),
	}

	userRes, err := r.userClient.GetUser(ctx, &userReq)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return entities.User{}, err
	}

	detailsRes, err := r.detailsClient.GetUserDetails(ctx, &detailsReq)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return entities.User{}, err
	}

	res := entities.User{
		Email:    userRes.GetEmail(),
		Password: userRes.GetPassword(),
		Age:      int(userRes.GetAge()),
		Details: entities.Details{
			Country:      detailsRes.GetCountry(),
			City:         detailsRes.GetCity(),
			MobileNumber: detailsRes.GetMobileNumber(),
			Married:      detailsRes.GetMarried(),
			Height:       detailsRes.GetHeight(),
			Weight:       detailsRes.GetWeight(),
		},
	}

	return res, nil
}

// DeleteUser does a soft delete in both databases inside each gRPC server
func (r *HTTPRepository) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(r.logger, "method", "delete_user")

	userReq := userpb.DeleteUserRequest{
		Id: uint32(id),
	}
	detailsReq := detailspb.DeleteUserDetailsRequest{
		UserId: uint32(id),
	}

	userRes, err := r.userClient.DeleteUser(ctx, &userReq)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return false, err
	}

	detailsRes, err := r.detailsClient.DeleteUserDetails(ctx, &detailsReq)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return false, err
	}

	return detailsRes.GetSuccess() && userRes.GetSuccess(), nil
}
