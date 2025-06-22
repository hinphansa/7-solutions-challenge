package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hinphansa/7-solutions-challenge/config"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, userHandler *UserHandler, authHandler *AuthHandler) {
	authMiddleware := AuthMiddleware(cfg.JWT.Secret)

	// Setup routes
	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			//// users endpoints
			{
				// Public users endpoints
				users := v1.Group("/users")
				users.Post("/", userHandler.Register)
				users.Get("/", userHandler.ListUsers)

				// Protected users endpoints
				authUsers := users.Group("/").Use(authMiddleware)
				authUsers.Get("/:id", userHandler.GetUser)
				authUsers.Put("/:id", userHandler.UpdateUser)
				authUsers.Delete("/:id", userHandler.DeleteUser)
			}

			//// auth endpoints
			auth := v1.Group("/auth")
			{
				auth.Post("/login", authHandler.Login)
			}
		}
	}
}
