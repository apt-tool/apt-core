package handler

import (
	"github.com/ptaas-tool/base-api/internal/core/worker"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	WorkerPool *worker.Pool
	Secret     string
}

// secure checks that the connection is from api
func (h Handler) secure(ctx *fiber.Ctx) error {
	cypher := ctx.Get("x-secure", "")

	if cypher != crypto.GetMD5Hash(h.Secret) {
		return ctx.Status(fiber.StatusForbidden).SendString("cannot access core")
	}

	return ctx.Next()
}

// process will perform the operation
func (h Handler) process(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("project_id", 0)

	if !h.WorkerPool.Do(id, false) {
		return ctx.SendStatus(fiber.StatusServiceUnavailable)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// rerun will perform operation of single document
func (h Handler) rerun(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("document_id", 0)

	if !h.WorkerPool.Do(id, true) {
		return ctx.SendStatus(fiber.StatusServiceUnavailable)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Register core apis
func (h Handler) Register(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Get("/api/:project_id", h.secure, h.process)
	app.Get("/api/rerun/:document_id", h.secure, h.rerun)
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
}
