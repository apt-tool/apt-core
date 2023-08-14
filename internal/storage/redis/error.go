package redis

import "errors"

var (
	ErrBadConnection = errors.New("failed to connect to redis cluster")
	ErrSaveSet       = errors.New("failed to save the set")
	ErrSetNotFound   = errors.New("set not found")
)
