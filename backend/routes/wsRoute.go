package routes

import (
	wsHandler "github.com/RugeFX/ruge-chat-app/handlers/ws"
	"github.com/gin-gonic/gin"
)

func RegisterWSRoute(r *gin.Engine) {
	ws := r.Group("/ws")

	ws.GET("/:id", gin.WrapF(wsHandler.HandleWs))
}
