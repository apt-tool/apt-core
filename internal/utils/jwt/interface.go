package jwt

import (
	"time"

	"github.com/automated-pen-testing/api/pkg/enum"
)

type Authenticator interface {
	GenerateToken(name string, role enum.Role) (string, time.Time, error)
	ParseToken(token string) (string, error)
}

func New(cfg Config) Authenticator {
	return &authenticator{
		key:    cfg.PrivateKey,
		expire: cfg.ExpireTime,
	}
}
