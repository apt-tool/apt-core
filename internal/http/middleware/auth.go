package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// Auth middleware check user authentication.
func (m Middleware) Auth(ctx *fiber.Ctx) error {
	if token := ctx.Get("x-token", ""); token != "" {
		if name, err := m.JWTAuthenticator.ParseToken(token); err == nil {
			ctx.Locals("name", name)

			return ctx.Next()
		} else {
			return m.ErrHandler.ErrUnauthorized(ctx, err)
		}
	} else {
		return m.ErrHandler.ErrUnauthorized(ctx, errors.New("token not found"))
	}
}
