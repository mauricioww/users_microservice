package transport

// SetUserDetailsRequest stores the data sent to gRPC SetUserDetails method
type SetUserDetailsRequest struct {
	UserID       int
	Country      string
	City         string
	MobileNumber string
	Married      bool
	Height       float32
	Weigth       float32
}

// GetUserDetailsRequest stores the data sent to gRPC GetUserDetails method
type GetUserDetailsRequest struct {
	UserID int
}

// DeleteUserDetailsRequest stores the data sent to gRPC DeleteUserDetails method
type DeleteUserDetailsRequest struct {
	UserID int
}
