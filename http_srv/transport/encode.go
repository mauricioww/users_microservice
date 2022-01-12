package transport

import "github.com/mauricioww/user_microsrv/http_srv/entities"

type (
	CreateUserResponse struct {
		Id               int    `json:"user_id"`
		Email            string `json:"email"`
		Password         string `json:"password"`
		Age              int    `json:"age"`
		entities.Details `json:"information"`
	}

	AuthenticateResponse struct {
		Token string `json:"token"`
	}

	UpdateUserResponse struct {
		Success bool `json:"success"`
	}

	GetUserResponse struct {
		Id               int    `json:"user_id"`
		Email            string `json:"email"`
		Password         string `json:"password"`
		Age              int    `json:"age"`
		entities.Details `json:"information"`
	}

	DeleteUserResponse struct {
		Success bool `json:"success"`
	}
)
