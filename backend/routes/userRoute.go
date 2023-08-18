package routes

import (
	userHandler "github.com/RugeFX/ruge-chat-app/handlers/user"
	"github.com/gin-gonic/gin"
)

// Registers the /users endpoint into the router
func RegisterUserRoute(r *gin.RouterGroup) {
	user := r.Group("/users")
	user.GET("/", userHandler.GetAllUsers)
	user.GET("/:username", userHandler.GetUserByUsername)
	user.POST("/", userHandler.CreateUser)
	user.DELETE("/:id", userHandler.DeleteUserByID)
	user.PUT("/:id", userHandler.UpdateUser)
}
