package model

type CustomError struct {
	Message  string
	Code     string
	HttpCode int
}

func (e CustomError) Error() string {
	return e.Message
}

var (
	NotiUserNotFoundError error = CustomError{Message: "noti user not found", Code: "NOTI.NOTI_USER_NOT_FOUND.0", HttpCode: 404}
)
