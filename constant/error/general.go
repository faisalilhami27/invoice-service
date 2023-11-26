package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrSQLError                = errors.New("database server failed to execute, please try again")
	ErrOrderDate               = errors.New("order date must be greater than now")
	ErrStatus                  = errors.New("invalid status")
	ErrTooManyRequest          = errors.New("too many request, please try again later")
)
