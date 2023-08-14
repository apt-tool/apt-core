package middleware

import (
	"github.com/automated-pen-testing/api/internal/http/controller/handler"
	"github.com/automated-pen-testing/api/internal/storage/redis"
	"github.com/automated-pen-testing/api/internal/utils/jwt"
	"github.com/automated-pen-testing/api/pkg/models"
)

type Middleware struct {
	JWTAuthenticator jwt.Authenticator
	Models           *models.Interface
	RedisConnector   redis.Connector
	ErrHandler       handler.ErrorHandler
}
