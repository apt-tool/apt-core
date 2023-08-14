package controller

import (
	"fmt"

	"github.com/automated-pen-testing/api/internal/http/request"
	"github.com/automated-pen-testing/api/internal/http/response"
	"github.com/automated-pen-testing/api/internal/utils/crypto"

	"github.com/gofiber/fiber/v2"
)

// CreateProject into system
func (c Controller) CreateProject(ctx *fiber.Ctx) error {
	req := new(request.ProjectRequest)
	namespaceID, _ := ctx.ParamsInt("namespace_id", 0)

	if err := ctx.BodyParser(&req); err != nil {
		return c.ErrHandler.ErrBodyParser(ctx, fmt.Errorf("[controller.project.Create] failed to parse body error=%w", err))
	}

	if err := c.Models.Projects.Create(req.ToModel(uint(namespaceID), ctx.Locals("name").(string))); err != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.project.Create] failed to create project error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// GetProject by its id
func (c Controller) GetProject(ctx *fiber.Ctx) error {
	projectID, _ := ctx.ParamsInt("project_id", 0)

	project, err := c.Models.Projects.GetByID(uint(projectID))
	if err != nil {
		return c.ErrHandler.ErrRecordNotFound(ctx, fmt.Errorf("[controller.project.Get] record not found error=%w", err))
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ProjectResponse{}.DTO(project))
}

// DeleteProject by its id
func (c Controller) DeleteProject(ctx *fiber.Ctx) error {
	projectID, _ := ctx.ParamsInt("project_id", 0)

	if err := c.Models.Projects.Delete(uint(projectID)); err != nil {
		return c.ErrHandler.ErrDatabase(ctx, fmt.Errorf("[controller.project.Create] failed to delete project error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// ExecuteProject will send http request to core
func (c Controller) ExecuteProject(ctx *fiber.Ctx) error {
	projectID, _ := ctx.ParamsInt("project_id", 0)
	url := fmt.Sprintf("%s/%d", c.Config.HTTP.Core, projectID)

	rsp, err := c.Client.Get(url, fmt.Sprintf("x-secure:%s", c.Config.Core.Secret))
	if err != nil {
		return c.ErrHandler.ErrLogical(ctx, fmt.Errorf("[controller.project.Execute] failed to execute project error=%w", err))
	}

	if rsp.StatusCode != 200 {
		return c.ErrHandler.ErrLogical(ctx, fmt.Errorf("[controller.project.Execute] failed to execute project error=%w", err))
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// DownloadProjectDocument will download the project document
func (c Controller) DownloadProjectDocument(ctx *fiber.Ctx) error {
	documentID, _ := ctx.ParamsInt("document_id", 0)

	cypher := crypto.GetMD5Hash(fmt.Sprintf("%s%d", c.Config.FTP.Access, documentID))

	url := fmt.Sprintf("%s/download?path=%d&token=%s", c.Config.FTP.Host, documentID, cypher)

	return ctx.Redirect(url, fiber.StatusPermanentRedirect)
}
