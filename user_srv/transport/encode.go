package transport

type (
	CreateUserResponse struct {
		Id int
	}

	AuthenticateResponse struct {
		Id int
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
