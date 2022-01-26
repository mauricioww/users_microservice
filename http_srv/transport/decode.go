package transport

import "github.com/mauricioww/user_microsrv/http_srv/entities"

// CreateUserRequest struct stores the data sent to users endpoint with POST action
type CreateUserRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	Age              int    `json:"age"`
	entities.Details `json:"information"`
}

// AuthenticateRequest struct stores the data sent to auth endpoint with POST action
type AuthenticateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest struct stores the data sent to users endpoint with PUT action
type UpdateUserRequest struct {
	UserID           int
	Email            string `json:"email"`
	Password         string `json:"password"`
	Age              int    `json:"age"`
	entities.Details `json:"information"`
}

// GetUserRequest struct stores the data sent to users endpoint with GET action
type GetUserRequest struct {
	UserID int
}

// DeleteUserRequest struct stores the data sent to users endpoint with DELETE action
type DeleteUserRequest struct {
	UserID int
}
