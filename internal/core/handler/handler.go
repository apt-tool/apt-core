package handler

import (
	"github.com/ptaas-tool/base-api/internal/core/worker"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"
	"github.com/ptaas-tool/base-api/pkg/enum"
	"github.com/ptaas-tool/base-api/pkg/models"
	"github.com/ptaas-tool/base-api/pkg/models/track"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	WorkerPool *worker.Pool
	Secret     string
	DB         *models.Interface
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

	_ = h.DB.Tracks.Create(&track.Track{
		ProjectID:   uint(id),
		Service:     "base-api",
		Description: "Got execute request",
		Type:        enum.TrackSuccess,
	})

	if !h.WorkerPool.Do(id, false) {
		_ = h.DB.Tracks.Create(&track.Track{
			ProjectID:   uint(id),
			Service:     "base-api",
			Description: "Base api worker pool is empty!",
			Type:        enum.TrackError,
		})

		return ctx.SendStatus(fiber.StatusServiceUnavailable)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// rerun will perform operation of single document
func (h Handler) rerun(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("document_id", 0)

	_ = h.DB.Tracks.Create(&track.Track{
		DocumentID:  uint(id),
		Service:     "base-api",
		Description: "Got rerun request",
		Type:        enum.TrackSuccess,
	})

	if !h.WorkerPool.Do(id, true) {
		_ = h.DB.Tracks.Create(&track.Track{
			DocumentID:  uint(id),
			Service:     "base-api",
			Description: "Base api worker pool is empty!",
			Type:        enum.TrackError,
		})

		return ctx.SendStatus(fiber.StatusServiceUnavailable)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Register core apis
func (h Handler) Register(app *fiber.App) {
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Get("/readyz", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(h.secure)

	app.Get("/api/:project_id", h.process)
	app.Get("/api/rerun/:document_id", h.rerun)
}
