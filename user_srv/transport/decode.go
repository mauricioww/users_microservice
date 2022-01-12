package transport

type (
	CreateUserRequest struct {
		Email    string
		Password string
		Age      int
	}

	AuthenticateRequest struct {
		Email    string
		Password string
	}

	UpdateUserRequest struct {
		Id       int
		Email    string
		Password string
		Age      int
	}

	GetUserRequest struct {
		UserId int
	}

	DeleteUserRequest struct {
		UserId int
	}
)
