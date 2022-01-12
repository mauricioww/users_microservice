package transport

type (
	SetUserDetailsResponse struct {
		Success bool
	}

	GetUserDetailsResponse struct {
		Country      string
		City         string
		MobileNumber string
		Married      bool
		Height       float32
		Weight       float32
	}

	DeleteUserDetailsResponse struct {
		Success bool
	}
)
