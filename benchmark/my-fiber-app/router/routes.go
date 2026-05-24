package router

import (
	"my-fiber-app/internal/user"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	user.RegisterRoutes(api)
}
