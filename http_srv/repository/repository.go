package repository

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"github.com/mauricioww/user_microsrv/user_srv/userpb"
	"google.golang.org/grpc"
)

type HttpRepository interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	Authenticate(ctx context.Context, session entities.Session) (int, error)
	UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type httpRepository struct {
	user_client    userpb.UserServiceClient
	details_client detailspb.UserDetailsServiceClient
	logger         log.Logger
}

func NewHttpRepository(conn1 *grpc.ClientConn, conn2 *grpc.ClientConn, logger log.Logger) HttpRepository {
	return &httpRepository{
		user_client:    userpb.NewUserServiceClient(conn1),
		details_client: detailspb.NewUserDetailsServiceClient(conn2),
		logger:         log.With(logger, "http_service", "repository"),
	}
}

func (r *httpRepository) CreateUser(ctx context.Context, user entities.User) (int, error) {
	logger := log.With(r.logger, "method", "create_users")
	res := -1

	userpb_req := userpb.CreateUserRequest{
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}

	user_res, err := r.user_client.CreateUser(ctx, &userpb_req)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return res, err
	}

	details_req := detailspb.SetUserDetailsRequest{
		UserId:       uint32(user_res.GetId()),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weight:       user.Details.Weight,
	}

	details_res, err := r.details_client.SetUserDetails(ctx, &details_req)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return res, err
	}

	if details_res.GetSuccess() {
		res = int(user_res.GetId())
	}

	return res, nil
}

func (r *httpRepository) Authenticate(ctx context.Context, session entities.Session) (int, error) {
	logger := log.With(r.logger, "method", "authenticate_user")
	res := -1

	auth_req := userpb.AuthenticateRequest{
		Email:    session.Email,
		Password: session.Password,
	}
	auth_res, err := r.user_client.Authenticate(ctx, &auth_req)

	if err != nil {
		level.Error(logger).Log("err_user", err)
	} else {
		res = int(auth_res.GetUserId())
	}

	return res, err
}

func (r *httpRepository) UpdateUser(ctx context.Context, user entities.UserUpdate) (bool, error) {
	logger := log.With(r.logger, "method", "update_user")
	var res bool

	user_req := userpb.UpdateUserRequest{
		Id:       uint32(user.UserId),
		Email:    user.Email,
		Password: user.Password,
		Age:      uint32(user.Age),
	}
	details_req := detailspb.SetUserDetailsRequest{
		UserId:       uint32(user.UserId),
		Country:      user.Details.Country,
		City:         user.Details.City,
		MobileNumber: user.Details.MobileNumber,
		Married:      user.Details.Married,
		Height:       user.Details.Height,
		Weight:       user.Details.Weight,
	}

	user_res, err := r.user_client.UpdateUser(ctx, &user_req)
	if err != nil {
		level.Error(logger).Log("err_user", err)
		return res, err
	}

	details_res, err := r.details_client.SetUserDetails(ctx, &details_req)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return res, err
	}

	res = details_res.GetSuccess() && user_res.GetSuccess()

	return res, nil
}

func (r *httpRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(r.logger, "method", "get_user")
	var res entities.User

	user_req := userpb.GetUserRequest{
		Id: uint32(id),
	}
	details_req := detailspb.GetUserDetailsRequest{
		UserId: uint32(id),
	}

	user_res, err := r.user_client.GetUser(ctx, &user_req)

	if err != nil {
		level.Error(logger).Log("err_user", err)
		return res, err
	}

	details_res, err := r.details_client.GetUserDetails(ctx, &details_req)
	if err != nil {
		level.Error(logger).Log("err_details", err)
		return res, err
	}

	res = entities.User{
		Email:    user_res.GetEmail(),
		Password: user_res.GetPassword(),
		Age:      int(user_res.GetAge()),
		Details: entities.Details{
			Country:      details_res.GetCountry(),
			City:         details_res.GetCity(),
			MobileNumber: details_res.GetMobileNumber(),
			Married:      details_res.GetMarried(),
			Height:       details_res.GetHeight(),
			Weight:       details_res.GetWeight(),
		},
	}

	return res, nil
}

func (r *httpRepository) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(r.logger, "method", "delete_user")
	var res bool

	user_req := userpb.DeleteUserRequest{
		Id: uint32(id),
	}
	user_res, err := r.user_client.DeleteUser(ctx, &user_req)

	if err != nil {
		level.Error(logger).Log("err_user", err)
		return res, err
	}

	details_req := detailspb.DeleteUserDetailsRequest{
		UserId: uint32(id),
	}
	details_res, err := r.details_client.DeleteUserDetails(ctx, &details_req)

	if err != nil {
		level.Error(logger).Log("err_details", err)
		return res, err
	}

	res = details_res.GetSuccess() && user_res.GetSuccess()

	return res, nil
}
