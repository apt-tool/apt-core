package jwt

import "errors"

var (
	ErrInvalidToken  = errors.New("token is invalid")
	ErrSigningMethod = errors.New("signing method is incorrect")
)
