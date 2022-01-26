package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
	"github.com/mauricioww/user_microsrv/user_srv/repository"
)

// GrpcUserServicer describe the business logic used to do validations and operations
type GrpcUserServicer interface {
	CreateUser(ctx context.Context, email string, pwd string, age int) (int, error)
	Authenticate(ctx context.Context, email string, pwd string) (bool, error)
	UpdateUser(ctx context.Context, id int, email string, pwd string, age int) (bool, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

// GrpcUserService implements the GrpcUserServicer interface
type GrpcUserService struct {
	repository repository.UserRepositorier
	logger     log.Logger
}

// NewGrpcUserService returns a GrpcUserService pointer type
func NewGrpcUserService(r repository.UserRepositorier, l log.Logger) *GrpcUserService {
	return &GrpcUserService{
		repository: r,
		logger:     l,
	}
}

// CreateUser does the both email and password validations and send the data to repository layer
func (g *GrpcUserService) CreateUser(ctx context.Context, email string, pwd string, age int) (int, error) {
	logger := log.With(g.logger, "method", "create_user")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	}

	if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return -1, e
	}

	cipheredPwd := helpers.Cipher(pwd)

	user := entities.User{
		Email:    email,
		Password: cipheredPwd,
		Age:      age,
	}

	res, err := g.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return -1, err
	}

	logger.Log("action", "success")
	return res, nil
}

// Authenticate does the both email and password validations and send the data to repository layer
func (g *GrpcUserService) Authenticate(ctx context.Context, email string, pwd string) (bool, error) {
	logger := log.With(g.logger, "method", "authenticate")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	}

	if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	}

	auth := entities.Session{
		Email:    email,
		Password: pwd,
	}

	hashedPwd, err := g.repository.Authenticate(ctx, &auth)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return false, err
	}

	if cipheredPwd := helpers.Cipher(auth.Password); cipheredPwd != hashedPwd {
		e := errors.NewUnauthenticatedError()
		level.Error(logger).Log("ERROR", e)
		return false, e
	}

	logger.Log("action", "success")
	return true, nil
}

// UpdateUser does the both email and password validations and send the data to repository layer
func (g *GrpcUserService) UpdateUser(ctx context.Context, id int, email string, pwd string, age int) (bool, error) {
	logger := log.With(g.logger, "method", "update_user")

	if email == "" {
		e := errors.NewBadRequestEmailError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	}

	if pwd == "" {
		e := errors.NewBadRequestPasswordError()
		level.Error(logger).Log("validation: ", e)
		return false, e
	}

	cipheredPwd := helpers.Cipher(pwd)

	updateInfo := entities.Update{
		UserID: id,
		User: entities.User{
			Email:    email,
			Password: cipheredPwd,
			Age:      age,
		},
	}

	u, err := g.repository.UpdateUser(ctx, updateInfo)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return false, err
	}

	logger.Log("action", "success")
	return u.Email == email && helpers.Decipher(u.Password) == pwd && u.Age == age, err
}

// GetUser receives one ID and send it to repository layer
func (g *GrpcUserService) GetUser(ctx context.Context, id int) (entities.User, error) {
	logger := log.With(g.logger, "method", "get_user")

	res, err := g.repository.GetUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return entities.User{}, err
	}

	logger.Log("action", "success")
	res.Password = helpers.Decipher(res.Password)
	return res, err
}

// DeleteUser receives one ID and send it to repository layer
func (g *GrpcUserService) DeleteUser(ctx context.Context, id int) (bool, error) {
	logger := log.With(g.logger, "method", "delete_user")

	res, err := g.repository.DeleteUser(ctx, id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, err
}
