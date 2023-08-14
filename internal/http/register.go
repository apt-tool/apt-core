package http

import (
	"github.com/automated-pen-testing/api/internal/config"
	"github.com/automated-pen-testing/api/internal/http/controller"
	"github.com/automated-pen-testing/api/internal/http/controller/handler"
	"github.com/automated-pen-testing/api/internal/http/middleware"
	"github.com/automated-pen-testing/api/internal/storage/redis"
	"github.com/automated-pen-testing/api/internal/utils/jwt"
	"github.com/automated-pen-testing/api/pkg/client"
	"github.com/automated-pen-testing/api/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type Register struct {
	Config          config.Config
	RedisConnector  redis.Connector
	ModelsInterface *models.Interface
}

func (r Register) Create(app *fiber.App) {
	// create new jwt authenticator
	authenticator := jwt.New(r.Config.JWT)

	errHandler := handler.ErrorHandler{DevMode: r.Config.HTTP.DevMode}

	// create middleware and controller
	mid := middleware.Middleware{
		JWTAuthenticator: authenticator,
		Models:           r.ModelsInterface,
		RedisConnector:   r.RedisConnector,
		ErrHandler:       errHandler,
	}
	ctl := controller.Controller{
		Config:           r.Config,
		JWTAuthenticator: authenticator,
		Models:           r.ModelsInterface,
		RedisConnector:   r.RedisConnector,
		ErrHandler:       errHandler,
		Client:           client.NewClient(),
	}

	// register endpoints
	app.Post("/login", ctl.UserLogin)

	auth := app.Use(mid.Auth)

	// viewer routes
	viewerRoutes := auth.Group("/")

	viewerRoutes.Get("/profile", ctl.GetUser)
	viewerRoutes.Post("/profile", ctl.UpdateUser)
	viewerRoutes.Get("/namespaces", ctl.GetUserNamespaces)
	viewerRoutes.Get("/namespaces/:namespace_id", mid.UserNamespace, ctl.GetNamespace)
	viewerRoutes.Get("/namespaces/:namespace_id/projects/:project_id", mid.UserNamespace, ctl.GetProject)
	viewerRoutes.Get("/namespaces/:namespace_id/projects/:project_id/:document_id", mid.UserNamespace, ctl.DownloadProjectDocument)

	// user routesxx
	userRoutes := auth.Group("/user")

	userRoutes.Post("/namespaces/:namespace_id/projects", mid.UserNamespace, ctl.CreateProject)
	userRoutes.Post("/namespaces/:namespace_id/projects/:project_id", mid.UserNamespace, ctl.ExecuteProject)
	userRoutes.Delete("/namespaces/:namespace_id/projects/:project_id", mid.UserNamespace, ctl.DeleteProject)

	// admin routes
	adminRoutes := auth.Use(mid.Admin).Group("/admin")

	users := adminRoutes.Group("/users")

	users.Get("/", ctl.GetUsersList)
	users.Post("/", ctl.UserRegister)
	users.Put("/", ctl.UpdateUserRole)
	users.Delete("/:user_id", ctl.DeleteUser)

	namespaces := adminRoutes.Group("/namespaces")

	namespaces.Get("/", ctl.GetNamespaces)
	namespaces.Post("/", ctl.CreateNamespace)
	namespaces.Put("/", ctl.UpdateNamespace)
	namespaces.Get("/:namespace_id", ctl.GetNamespaceUsers)
	namespaces.Delete("/:namespace_id", ctl.DeleteNamespace)
}
