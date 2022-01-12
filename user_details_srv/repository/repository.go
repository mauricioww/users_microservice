package repository

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDetailsRepository interface {
	SetUserDetails(ctx context.Context, info entities.UserDetails) (bool, error)
	GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error)
	DeleteUserDetails(ctx context.Context, user_id int) (bool, error)
}

type userDetailsRepository struct {
	db     *mongo.Database
	logger log.Logger
}

func NewUserDetailsRepository(mongo_db *mongo.Database, l log.Logger) UserDetailsRepository {
	return &userDetailsRepository{
		db:     mongo_db,
		logger: log.With(l, "repository", "mongodb"),
	}
}

func (r *userDetailsRepository) SetUserDetails(ctx context.Context, details entities.UserDetails) (bool, error) {
	collection := r.db.Collection("information")
	var err error

	if helpers.NoExists(collection, ctx, details.UserId) {
		_, err = collection.InsertOne(ctx, helpers.BuildInsertBson(details))
	} else {
		_, err = collection.UpdateByID(ctx, details.UserId, helpers.BuildUpdateBson(details))
	}

	if err != nil {
		return false, errors.NewInternalError()
	}

	return true, nil
}

func (r *userDetailsRepository) GetUserDetails(ctx context.Context, user_id int) (entities.UserDetails, error) {
	collection := r.db.Collection("information")
	var res entities.UserDetails

	if helpers.NoExists(collection, ctx, user_id) {
		return res, errors.NewUserNotFoundError()
	} else if err := collection.FindOne(ctx, bson.D{{"_id", user_id}}).Decode(&res); err != nil {
		return res, errors.NewInternalError()
	}

	return res, nil
}

func (r *userDetailsRepository) DeleteUserDetails(ctx context.Context, user_id int) (bool, error) {
	collection := r.db.Collection("information")

	if helpers.NoExists(collection, ctx, user_id) {
		return false, errors.NewUserNotFoundError()
	} else if _, err := collection.DeleteOne(ctx, bson.D{{"_id", user_id}}); err != nil {
		return false, errors.NewInternalError()
	}

	return true, nil
}
