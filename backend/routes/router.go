package routes

import (
	"github.com/gin-gonic/gin"
)

// Sets up all of the routes for the API
func SetupRoutes(app *gin.Engine) {
	api := app.Group("/api")
	// /ws WebSocket endpoint
	RegisterWSRoute(app)
	// rest REST api endpoint
	RegisterUserRoute(api)
}
