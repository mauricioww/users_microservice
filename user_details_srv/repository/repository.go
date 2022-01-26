package repository

import (
	"context"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/helpers"
	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserDetailsRepositorier describes the methods used to do DB operations
type UserDetailsRepositorier interface {
	SetUserDetails(ctx context.Context, info entities.UserDetails) (bool, error)
	GetUserDetails(ctx context.Context, UserID int) (entities.UserDetails, error)
	DeleteUserDetails(ctx context.Context, UserID int) (bool, error)
}

// UserDetailsRepository implements the UserDetailsRepositorier interface
type UserDetailsRepository struct {
	db     *mongo.Database
	logger log.Logger
}

// NewUserDetailsRepository returns a UserDetailsRepository pointer type
func NewUserDetailsRepository(mongoDb *mongo.Database, l log.Logger) *UserDetailsRepository {
	return &UserDetailsRepository{
		db:     mongoDb,
		logger: log.With(l, "repository", "mongodb"),
	}
}

// SetUserDetails does the DB operation to insert or update information for a specific user
func (r *UserDetailsRepository) SetUserDetails(ctx context.Context, details entities.UserDetails) (bool, error) {
	collection := r.db.Collection("information")
	details.Active = true
	var err error

	if helpers.NoExists(ctx, collection, details.UserID) {
		_, err = collection.InsertOne(ctx, helpers.BuildInsertBson(details))
	} else {
		_, err = collection.UpdateByID(ctx, details.UserID, helpers.BuildUpdateBson(details))
	}

	if err != nil {
		return false, errors.NewInternalError()
	}

	return true, nil
}

// GetUserDetails fetchs the information within the DB about a specific user
func (r *UserDetailsRepository) GetUserDetails(ctx context.Context, UserID int) (entities.UserDetails, error) {
	collection := r.db.Collection("information")
	var res entities.UserDetails

	if helpers.NoExists(ctx, collection, UserID) {
		return res, errors.NewUserNotFoundError()
	}

	if err := collection.FindOne(ctx, bson.D{{"_id", UserID}}).Decode(&res); err != nil {
		return res, errors.NewInternalError()
	}

	return res, nil
}

// DeleteUserDetails does a soft delete operation over a specific user
func (r *UserDetailsRepository) DeleteUserDetails(ctx context.Context, UserID int) (bool, error) {
	collection := r.db.Collection("information")
	var data entities.UserDetails

	if helpers.NoExists(ctx, collection, UserID) {
		return false, errors.NewUserNotFoundError()
	}

	if err := collection.FindOne(ctx, bson.D{{"_id", UserID}}).Decode(&data); err != nil {
		return false, errors.NewInternalError()
	}

	data.Active = false

	if _, err := collection.UpdateByID(ctx, UserID, helpers.BuildUpdateBson(data)); err != nil {
		return false, errors.NewInternalError()
	}

	return true, nil
}
