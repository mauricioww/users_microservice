package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

type HttpService interface {
	CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error)
	Authenticate(ctx context.Context, email string, pwd string) (bool, error)
	UpdateUser(ctx context.Context, user_id int, email string, pwd string, age int, details entities.Details) (bool, error)
	GetUser(ctx context.Context, user_id int) (entities.User, error)
	DeleteUser(ctx context.Context, user_id int) (bool, error)
}

type httpService struct {
	repository repository.HttpRepository
	logger     log.Logger
}

func NewHttpService(r repository.HttpRepository, l log.Logger) HttpService {
	return &httpService{
		logger:     l,
		repository: r,
	}
}

func (s *httpService) CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error) {
	logger := log.With(s.logger, "method", "create_user")

	user := entities.User{
		Email:    email,
		Password: pwd,
		Age:      age,
		Details:  details,
	}

	res, err := s.repository.CreateUser(ctx, user)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (s *httpService) Authenticate(ctx context.Context, email string, pwd string) (bool, error) {
	logger := log.With(s.logger, "method", "authenticate")

	session := entities.Session{
		Email:    email,
		Password: pwd,
	}

	res, err := s.repository.Authenticate(ctx, session)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (s *httpService) UpdateUser(ctx context.Context, user_id int, email string, pwd string, age int, details entities.Details) (bool, error) {
	logger := log.With(s.logger, "method", "update_user")
	info_update := entities.UserUpdate{
		UserId: user_id,
		User: entities.User{
			Email:    email,
			Password: pwd,
			Age:      age,
			Details:  details,
		},
	}

	res, err := s.repository.UpdateUser(ctx, info_update)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (s *httpService) GetUser(ctx context.Context, user_id int) (entities.User, error) {
	logger := log.With(s.logger, "method", "get_user")

	res, err := s.repository.GetUser(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (s *httpService) DeleteUser(ctx context.Context, user_id int) (bool, error) {
	logger := log.With(s.logger, "method", "delete_user")

	res, err := s.repository.DeleteUser(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}
