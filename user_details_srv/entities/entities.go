package entities

type (
	UserDetails struct {
		UserId       int     `bson:"_id"`
		Country      string  `bson:"country"`
		City         string  `bson:"city"`
		MobileNumber string  `bson:"mobile_number"`
		Married      bool    `bson:"married"`
		Height       float32 `bson:"height"`
		Weight       float32 `bson:"weight"`
	}
)
