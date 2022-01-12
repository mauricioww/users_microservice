package transport

import "github.com/mauricioww/user_microsrv/http_srv/entities"

type (
	CreateUserRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		Age              int    `json:"age"`
		entities.Details `json:"information"`
	}

	AuthenticateRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserRequest struct {
		UserId           int
		Email            string `json:"email"`
		Password         string `json:"password"`
		Age              int    `json:"age"`
		entities.Details `json:"information"`
	}

	GetUserRequest struct {
		UserId int
	}

	DeleteUserRequest struct {
		UserId int
	}
)
