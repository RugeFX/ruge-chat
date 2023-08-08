package routes

import (
	"github.com/gofiber/fiber/v2"
)

// Sets up all of the routes for the API
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	RegisterUserRoute(api)
}
