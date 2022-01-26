package transport

import "github.com/mauricioww/user_microsrv/http_srv/entities"

// CreateUserResponse struct stores the data that users endpoint, with POST action, will return
type CreateUserResponse struct {
	UserID           int    `json:"user_id"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Age              int    `json:"age"`
	entities.Details `json:"information"`
}

// AuthenticateResponse struct stores the data that auth endpoint, with POST action, will return
type AuthenticateResponse struct {
	Success bool `json:"success"`
}

// UpdateUserResponse struct stores the data that users endpoint, with PUT action, will return
type UpdateUserResponse struct {
	Success bool `json:"success"`
}

// GetUserResponse struct stores the data that users endpoint, with GET action, will return
type GetUserResponse struct {
	UserID           int    `json:"user_id"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Age              int    `json:"age"`
	entities.Details `json:"information"`
}

// DeleteUserResponse struct stores the data that users endpoint, with DELETE action, will return
type DeleteUserResponse struct {
	Success bool `json:"success"`
}
