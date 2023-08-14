package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
)
