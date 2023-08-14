package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type connector struct {
	connection *redis.Client
}

func (c *connector) Get(key string) (string, error) {
	ctx := context.Background()

	value, err := c.connection.Get(ctx, key).Result()
	if err != nil {
		return "", ErrSetNotFound
	}

	return value, nil
}

func (c *connector) Set(key, value string, expire time.Duration) error {
	ctx := context.Background()

	if err := c.connection.Set(ctx, key, value, expire).Err(); err != nil {
		return ErrSaveSet
	}

	return nil
}
