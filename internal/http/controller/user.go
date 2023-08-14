package controller

import (
	"fmt"
	"strconv"
	"time"

	"github.com/automated-pen-testing/api/internal/http/request"
	"github.com/automated-pen-testing/api/internal/http/response"

	"github.com/gofiber/fiber/v2"
)

// UserRegister will create a new user into system
func (c Controller) UserRegister(ctx *fiber.Ctx) error {
	req := new(request.UserRegisterRequest)

	if err := ctx.BodyParser(req); err != nil {
		return c.ErrHandler.ErrBodyParser(ctx, fmt.Errorf("[controller.user.Register] failed to parse body error=%w", err))
	}

	if err := req.Validate(); err != nil {
		return c.ErrHandler.ErrValidation(ctx, fmt.Errorf("[controller.user.Register] failed to validate request error=%w", err))
	}

	if err := c.Models.Users.Create(req.ToModel()); err != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.user.Register] failed to create user error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// UserLogin logs in a user
func (c Controller) UserLogin(ctx *fiber.Ctx) error {
	req := new(request.UserRegisterRequest)

	if err := ctx.BodyParser(req); err != nil {
		return c.ErrHandler.ErrBodyParser(ctx, fmt.Errorf("[controller.user.Loing] failed to parse body error=%w", err))
	}

	if err := req.Validate(); err != nil {
		return c.ErrHandler.ErrValidation(ctx, fmt.Errorf("[controller.user.Login] failed to validate request error=%w", err))
	}

	userTmp, err := c.Models.Users.Validate(req.Name, req.Pass)
	if err != nil {
		return c.ErrHandler.ErrRecordNotFound(ctx, fmt.Errorf("[controller.user.Login] username and password don't match error=%w", err))
	}

	token, etime, err := c.JWTAuthenticator.GenerateToken(userTmp.Username, userTmp.Role)
	if err != nil {
		return c.ErrHandler.ErrLogical(ctx, fmt.Errorf("[controller.user.Loing] failed to create token error=%w", err))
	}

	if er := c.RedisConnector.Set(userTmp.Username, strconv.Itoa(int(userTmp.Role)), etime.Sub(time.Now())); er != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.user.Loing] failed to save token error=%w", er))
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Token{
		Token: token,
	})
}

// GetUser profile
func (c Controller) GetUser(ctx *fiber.Ctx) error {
	record, err := c.Models.Users.GetByName(ctx.Locals("name").(string))
	if err != nil {
		return c.ErrHandler.ErrRecordNotFound(ctx, fmt.Errorf("[controller.user.Get] username and password don't match error=%w", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response.UserResponse{}.DTO(record))
}

// GetUsersList returns the list of users
func (c Controller) GetUsersList(ctx *fiber.Ctx) error {
	list, err := c.Models.Users.Get()
	if err != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.user.List] failed to get users error=%w", err))
	}

	records := make([]*response.UserResponse, 0)

	for _, item := range list {
		records = append(records, response.UserResponse{}.DTO(item))
	}

	return ctx.Status(fiber.StatusOK).JSON(records)
}

// UpdateUser information
func (c Controller) UpdateUser(ctx *fiber.Ctx) error {
	req := new(request.UserRegisterRequest)

	if err := ctx.BodyParser(req); err != nil {
		return c.ErrHandler.ErrBodyParser(ctx, fmt.Errorf("[controller.user.Update] failed to parse body error=%w", err))
	}

	if er := c.Models.Users.UpdateInfo(ctx.Locals("name").(string), req.Name); er != nil {
		return c.ErrHandler.ErrRecordNotFound(ctx, fmt.Errorf("[controller.user.Update] failed to update user error=%w", er))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// UpdateUserRole changes the users role
func (c Controller) UpdateUserRole(ctx *fiber.Ctx) error {
	req := new(request.UserRoleUpdateRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return c.ErrHandler.ErrBodyParser(ctx, fmt.Errorf("[controller.user.Update] failed to parse body error=%w", err))
	}

	u, err := c.Models.Users.GetByID(req.UserID)
	if err != nil {
		return c.ErrHandler.ErrRecordNotFound(ctx, fmt.Errorf("[controller.user.Update] failed to find user error=%w", err))
	}

	u.Role = req.Role

	if er := c.Models.Users.Update(req.UserID, u); er != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.user.Update] failed to update user error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// DeleteUser removes user
func (c Controller) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := ctx.ParamsInt("user_id")

	if err := c.Models.Users.Delete(uint(id)); err != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.user.Delete] failed to delete user error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}
