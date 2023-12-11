package errors

import "errors"

const (
	ErrSomethingWentWrong = "something went wrong"
	ErrUnauthorized       = "unauthorized"
	ErrInvalidJson        = "Invalid JSON"
	ErrMissingBody        = "missing body request"
)

var ErrInvalidID = errors.New("id is not in its proper form")
