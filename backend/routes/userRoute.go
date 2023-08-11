package routes

import (
	userHandler "github.com/RugeFX/ruge-chat-app/handlers/user"
	"github.com/gofiber/fiber/v2"
)

// Registers the /users endpoint into the router
func RegisterUserRoute(r fiber.Router) {
	user := r.Group("/users")
	user.Get("/", userHandler.GetAllUsers)
	user.Get("/:username", userHandler.GetUserByUsername)
	user.Post("/", userHandler.CreateUser)
	user.Delete("/:id", userHandler.DeleteUserByID)
	user.Put("/:id", userHandler.UpdateUser)
}
