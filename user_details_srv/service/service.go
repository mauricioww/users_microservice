package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"github.com/mauricioww/user_microsrv/user_details_srv/repository"
)

type GrpcUserDetailsService interface {
	SetUserDetails(ctx context.Context, user_id int, country string, city string, mobile_number string, married bool, height float32, weigth float32) (bool, error)
	GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error)
	DeleteUserDetails(ctx context.Context, user_id int) (bool, error)
}

type grpcUserDetailsService struct {
	repository repository.UserDetailsRepository
	logger     log.Logger
}

func NewGrpcUserDetailsService(r repository.UserDetailsRepository, l log.Logger) GrpcUserDetailsService {
	return &grpcUserDetailsService{
		repository: r,
		logger:     l,
	}
}

func (g *grpcUserDetailsService) SetUserDetails(ctx context.Context, user_id int, country string, city string, mobile_number string, married bool, height float32, weight float32) (bool, error) {
	logger := log.With(g.logger, "method", "set_user_details")

	information := entities.UserDetails{
		UserId:       user_id,
		Country:      country,
		City:         city,
		MobileNumber: mobile_number,
		Married:      married,
		Height:       height,
		Weight:       weight,
	}

	res, err := g.repository.SetUserDetails(ctx, information)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (g *grpcUserDetailsService) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	logger := log.With(g.logger, "method", "get_user_details")
	res, err := g.repository.GetUserDetails(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}

func (g *grpcUserDetailsService) DeleteUserDetails(ctx context.Context, user_id int) (bool, error) {
	logger := log.With(g.logger, "method", "delete_user_details")
	res, err := g.repository.DeleteUserDetails(ctx, user_id)

	if err != nil {
		level.Error(logger).Log("ERROR", err)
	} else {
		logger.Log("action", "success")
	}

	return res, err
}
