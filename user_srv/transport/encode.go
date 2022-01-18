package transport

type (
	CreateUserResponse struct {
		Id       int
		Email    string
		Password string
		Age      int
	}

	AuthenticateResponse struct {
		Success bool
	}

	UpdateUserResponse struct {
		Success bool
	}

	GetUserResponse struct {
		Email    string
		Password string
		Age      int
	}

	DeleteUserResponse struct {
		Success bool
	}
)
