package transport

// CreateUserResponse stores the data that gRPC CreateUser method will return
type CreateUserResponse struct {
	UserID   int
	Email    string
	Password string
	Age      int
}

// AuthenticateResponse stores the data that gRPC Authenticate method will return
type AuthenticateResponse struct {
	Success bool
}

// UpdateUserResponse stores the data that gRPC Authenticate method will return
type UpdateUserResponse struct {
	Success bool
}

// GetUserResponse stores the data that gRPC GetUser method will return
type GetUserResponse struct {
	UserID   int
	Email    string
	Password string
	Age      int
}

// DeleteUserResponse stores the data that gRPC DeleteUser method will return
type DeleteUserResponse struct {
	Success bool
}
