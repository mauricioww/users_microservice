package entities

// Details struct stores the user's extra information
type Details struct {
	Country      string  `json:"country"`
	City         string  `json:"city"`
	MobileNumber string  `json:"mobile_number"`
	Married      bool    `json:"married"`
	Height       float32 `json:"height_m"`
	Weight       float32 `json:"weight_kg"`
}

// User struct stores the user's basic information
type User struct {
	Email    string
	Password string
	Age      int
	Details
}

// Session struct stores the credentials of the user which attempt to login
type Session struct {
	Email    string
	Password string
}

// UserUpdate struct stores new information to update the old one
type UserUpdate struct {
	UserID int
	User
}
