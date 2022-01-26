package helpers

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_details_srv/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BuildUpdateBson returns a bson with the new information to update
func BuildUpdateBson(d entities.UserDetails) bson.D {
	return bson.D{
		{"$set", injectFields(d)},
	}
}

// BuildInsertBson returns a simple bson to create a record within the database
func BuildInsertBson(d entities.UserDetails) bson.M {
	b := injectFields(d)
	b["_id"] = d.UserID
	return b
}

// NoExists returns true if the user is not into the database otherwise returns false
func NoExists(ctx context.Context, coll *mongo.Collection, id int) bool {
	var results entities.UserDetails
	if err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&results); err == mongo.ErrNoDocuments {
		return true
	}
	return !results.Active
}

func injectFields(d entities.UserDetails) bson.M {
	return bson.M{
		"country":       d.Country,
		"city":          d.City,
		"mobile_number": d.MobileNumber,
		"married":       d.Married,
		"height":        d.Height,
		"weight":        d.Weight,
		"active":        d.Active,
	}
}
