package transport

type (
	SetUserDetailsRequest struct {
		UserId       int
		Country      string
		City         string
		MobileNumber string
		Married      bool
		Height       float32
		Weigth       float32
	}

	GetUserDetailsRequest struct {
		UserId int
	}

	DeleteUserDetailsRequest struct {
		UserId int
	}
)
