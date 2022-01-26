package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/http_srv/entities"
	"github.com/mauricioww/user_microsrv/http_srv/repository"
)

// HTTPServicer describes the logic business of the services
type HTTPServicer interface {
	CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error)
	Authenticate(ctx context.Context, email string, pwd string) (bool, error)
	UpdateUser(ctx context.Context, userID int, email string, pwd string, age int, details entities.Details) (bool, error)
	GetUser(ctx context.Context, userID int) (entities.User, error)
	DeleteUser(ctx context.Context, userID int) (bool, error)
}

// HTTPService type implement the HTTPServicer interface
type HTTPService struct {
	repository repository.HTTPRepositorier
	logger     log.Logger
}

// NewHTTPService returns a HTTPService pointer type
func NewHTTPService(r repository.HTTPRepositorier, l log.Logger) *HTTPService {
	return &HTTPService{
		logger:     l,
		repository: r,
	}
}

// CreateUser receives data for a new user and send it to the repository
func (s *HTTPService) CreateUser(ctx context.Context, email string, pwd string, age int, details entities.Details) (int, error) {
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
		return -1, err
	}

	logger.Log("action", "success")
	return res, nil
}

// Authenticate receives data of a user to do a login and send it to repository
func (s *HTTPService) Authenticate(ctx context.Context, email string, pwd string) (bool, error) {
	logger := log.With(s.logger, "method", "authenticate")

	session := entities.Session{
		Email:    email,
		Password: pwd,
	}

	res, err := s.repository.Authenticate(ctx, session)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}

// UpdateUser receives new data to replace the old data of a user and send it to repository
func (s *HTTPService) UpdateUser(ctx context.Context, userID int, email string, pwd string, age int, details entities.Details) (bool, error) {
	logger := log.With(s.logger, "method", "update_user")
	info := entities.UserUpdate{
		UserID: userID,
		User: entities.User{
			Email:    email,
			Password: pwd,
			Age:      age,
			Details:  details,
		},
	}

	res, err := s.repository.UpdateUser(ctx, info)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}

// GetUser receives one ID and send it to repository
func (s *HTTPService) GetUser(ctx context.Context, userID int) (entities.User, error) {
	logger := log.With(s.logger, "method", "get_user")

	res, err := s.repository.GetUser(ctx, userID)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return entities.User{}, err
	}

	logger.Log("action", "success")
	return res, nil
}

// DeleteUser receives one ID and send it to repository
func (s *HTTPService) DeleteUser(ctx context.Context, userID int) (bool, error) {
	logger := log.With(s.logger, "method", "delete_user")

	res, err := s.repository.DeleteUser(ctx, userID)

	if err != nil {
		level.Error(logger).Log("ERROR: ", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}
