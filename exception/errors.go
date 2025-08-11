package exception

import "errors"

var (
	ErrConflict             = errors.New("username already exist")
	ErrUnauthorized         = errors.New("invalid or missing token")
	ErrUnauthorizedLogin    = errors.New("invalid username or password")
	ErrBadRequestTimeFormat = errors.New("invalid date format, please use rfc 3339")
	ErrNotFoundTask         = errors.New("task is not found")
)
