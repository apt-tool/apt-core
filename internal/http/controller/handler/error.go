package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct {
	DevMode bool
}

func (e ErrorHandler) ErrBodyParser(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.BodyParser] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
}

func (e ErrorHandler) ErrQueryParser(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.QueryParser] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
}

func (e ErrorHandler) ErrRecordNotFound(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.RecordNotFound] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
}

func (e ErrorHandler) ErrValidation(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.Validation] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
}

func (e ErrorHandler) ErrDatabase(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.Database] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
}

func (e ErrorHandler) ErrLogical(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.Logical] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
}

func (e ErrorHandler) ErrUnauthorized(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.Unauthorized] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusUnauthorized).SendString(err.Error())
}

func (e ErrorHandler) ErrAccess(ctx *fiber.Ctx, err error) error {
	if e.DevMode {
		log.Printf("[err.Access] error: %v\n", err)
	}

	return ctx.Status(fiber.StatusForbidden).SendString(err.Error())
}
