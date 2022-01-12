package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
)

type GrpcUserService interface {
	CreateUser(ctx context.Context, email string, pwd string, age int) (int, error)
	Authenticate(ctx context.Context, email string, pwd string) (int, error)
	UpdateUser(ctx context.Context, id int, email string, pwd string, age int) (bool, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type grpcUserService struct {
	repository repository.UserRepository
	logger     log.Logger
}

func NewGrpcUserService(r repository.UserRepository, l log.Logger) GrpcUserService {
	return &grpcUserService{
		repository: r,
		logger:     l,
	}
}

func (g *grpcUserService) CreateUser(ctx context.Context, email string, pwd string, age int) (int, error) {
	logger := log.With(g.logger, "method", "create_user")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	} else if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	}

	ciphered_pwd := helpers.Cipher(pwd)

	user := entities.User{
		Email:    email,
		Password: ciphered_pwd,
		Age:      age,
	}

	res, err := g.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (g *grpcUserService) Authenticate(ctx context.Context, email string, pwd string) (int, error) {
	logger := log.With(g.logger, "method", "authenticate")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	} else if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	}

	auth := entities.Session{
		Email:    email,
		Password: pwd,
	}

	hashed_pwd, err := g.repository.Authenticate(ctx, &auth)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return -1, err
	} else if ciphered_pwd := helpers.Cipher(auth.Password); ciphered_pwd != hashed_pwd {
		e := errors.NewUnauthenticatedError()
		level.Error(logger).Log("ERROR", e)
		return -1, e
	}

	logger.Log("action", "success")
	return auth.Id, nil
}

func (g *grpcUserService) UpdateUser(ctx context.Context, id int, email string, pwd string, age int) (bool, error) {
	logger := log.With(g.logger, "method", "update_user")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	} else if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	}

	ciphered_pwd := helpers.Cipher(pwd)
	var res bool

	update_info := entities.Update{
		UserId: id,
		User: entities.User{
			Email:    email,
			Password: ciphered_pwd,
			Age:      age,
		},
	}

	u, err := g.repository.UpdateUser(ctx, update_info)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else if u.Email == email && helpers.Decipher(u.Password) == pwd && u.Age == age {
		logger.Log("action", "success")
		res = true
	}

	return res, err
}

func (g *grpcUserService) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(g.logger, "method", "get_user")

	res, err := g.repository.GetUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		res.Password = helpers.Decipher(res.Password)
		logger.Log("action", "success")
	}

	return res, err
}

func (g *grpcUserService) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(g.logger, "method", "delete_user")

	res, err := g.repository.DeleteUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}
