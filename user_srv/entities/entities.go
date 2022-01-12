package entities

type (
	User struct {
		Email    string
		Password string
		Age      int
	}

	Session struct {
		Email    string
		Password string
		Id       int
	}

	Update struct {
		UserId int
		User
	}
)
