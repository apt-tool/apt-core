package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Connector interface {
	Get(key string) (string, error)
	Set(key, value string, expire time.Duration) error
}

func New(cfg Config) (Connector, error) {
	ctx := context.Background()

	conn := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Pass,
		DB:       0,
	})

	if err := conn.Ping(ctx).Err(); err != nil {
		return nil, ErrBadConnection
	}

	return &connector{
		connection: conn,
	}, nil
}
