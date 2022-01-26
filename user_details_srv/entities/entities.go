package entities

// UserDetails stores the user's information
type UserDetails struct {
	UserID       int     `bson:"_id"`
	Country      string  `bson:"country"`
	City         string  `bson:"city"`
	MobileNumber string  `bson:"mobile_number"`
	Married      bool    `bson:"married"`
	Height       float32 `bson:"height"`
	Weight       float32 `bson:"weight"`
	Active       bool    `bson:"active"`
}
