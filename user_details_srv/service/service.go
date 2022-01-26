package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/repository"
)

// GrpcUserDetailsServicer describe the business logic used to do validations and operations
type GrpcUserDetailsServicer interface {
	SetUserDetails(ctx context.Context, UserID int, country string, city string, number string, married bool, height float32, weigth float32) (bool, error)
	GetUserDetails(ctx context.Context, UserID int) (entities.UserDetails, error)
	DeleteUserDetails(ctx context.Context, UserID int) (bool, error)
}

// GrpcUserDetailsService implements the GrpcUserDetailsServicer interface
type GrpcUserDetailsService struct {
	repository repository.UserDetailsRepositorier
	logger     log.Logger
}

// NewGrpcUserDetailsService returns a GrpcUserDetailsService pointer type
func NewGrpcUserDetailsService(r repository.UserDetailsRepositorier, l log.Logger) *GrpcUserDetailsService {
	return &GrpcUserDetailsService{
		repository: r,
		logger:     l,
	}
}

// SetUserDetails receives data for user and send it to the repository
func (g *GrpcUserDetailsService) SetUserDetails(ctx context.Context, UserID int, country string, city string, number string, married bool, height float32, weight float32) (bool, error) {
	logger := log.With(g.logger, "method", "set_user_details")

	information := entities.UserDetails{
		UserID:       UserID,
		Country:      country,
		City:         city,
		MobileNumber: number,
		Married:      married,
		Height:       height,
		Weight:       weight,
	}

	res, err := g.repository.SetUserDetails(ctx, information)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}

// GetUserDetails receives one ID and send it to the repository
func (g *GrpcUserDetailsService) GetUserDetails(ctx context.Context, UserID int) (entities.UserDetails, error) {
	logger := log.With(g.logger, "method", "get_user_details")
	res, err := g.repository.GetUserDetails(ctx, UserID)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return entities.UserDetails{}, err
	}

	logger.Log("action", "success")
	return res, nil
}

// DeleteUserDetails receives one ID and send it to the repository
func (g *GrpcUserDetailsService) DeleteUserDetails(ctx context.Context, UserID int) (bool, error) {
	logger := log.With(g.logger, "method", "delete_user_details")
	res, err := g.repository.DeleteUserDetails(ctx, UserID)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
		return false, err
	}

	logger.Log("action", "success")
	return res, nil
}
