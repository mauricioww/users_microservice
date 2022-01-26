package transport

// SetUserDetailsResponse stores the data that gRPC GetUserDetails method will return
type SetUserDetailsResponse struct {
	Success bool
}

// GetUserDetailsResponse stores the data that gRPC GetUserDetails method will return
type GetUserDetailsResponse struct {
	Country      string
	City         string
	MobileNumber string
	Married      bool
	Height       float32
	Weight       float32
}

// DeleteUserDetailsResponse stores the data sent that gRPC DeleteUserDetails method will return
type DeleteUserDetailsResponse struct {
	Success bool
}
