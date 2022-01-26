package entities

// User struct stores the basic information
type User struct {
	Email    string
	Password string
	Age      int
}

// Session struct stores credentials to do a login
type Session struct {
	Email    string
	Password string
}

// Update struct stores new information to replace the old one
type Update struct {
	UserID int
	User
}
