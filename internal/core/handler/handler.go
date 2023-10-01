package handler

import (
	"github.com/ptaas-tool/base-api/internal/core/worker"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	WorkerPool *worker.Pool
	Secret     string
}

// secure checks that the connection is from api
func (h Handler) secure(ctx *fiber.Ctx) error {
	cypher := ctx.Get("x-secure", "")

	if crypto.GetMD5Hash(cypher) != h.Secret {
		return ctx.Status(fiber.StatusForbidden).SendString("cannot access core")
	}

	return ctx.Next()
}

// process will perform the operation
func (h Handler) process(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("project_id", 0)

	if !h.WorkerPool.Do(id) {
		return ctx.SendStatus(fiber.StatusServiceUnavailable)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Register core apis
func (h Handler) Register(app *fiber.App) {
	app.Get("/api/:project_id", h.process)
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
