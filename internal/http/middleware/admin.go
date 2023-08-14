package middleware

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/automated-pen-testing/api/pkg/enum"

	"github.com/gofiber/fiber/v2"
)

// Admin middleware checks the user admin role.
func (m Middleware) Admin(ctx *fiber.Ctx) error {
	tmp, err := m.RedisConnector.Get(ctx.Locals("name").(string))
	if err != nil {
		return m.ErrHandler.ErrUnauthorized(ctx, fmt.Errorf("token expired error=%v", err))
	}

	role, _ := strconv.Atoi(tmp)

	if role == int(enum.RoleAdmin) {
		return ctx.Next()
	}

	return m.ErrHandler.ErrAccess(ctx, errors.New("user cannot access this endpoint"))
}
