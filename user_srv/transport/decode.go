package transport

// CreateUserRequest stores the data sent to gRPC CreateUser method
type CreateUserRequest struct {
	Email    string
	Password string
	Age      int
}

// AuthenticateRequest stores the data sent to gRPC Authenticate method
type AuthenticateRequest struct {
	Email    string
	Password string
}

// UpdateUserRequest stores the data sent to gRPC UpdateUser method
type UpdateUserRequest struct {
	UserID   int
	Email    string
	Password string
	Age      int
}

// GetUserRequest stores the data sent to gRPC GetUser method
type GetUserRequest struct {
	UserID int
}

// DeleteUserRequest stores the data sent to gRPC DeleteUser method
type DeleteUserRequest struct {
	UserID int
}
