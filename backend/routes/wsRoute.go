package routes

import (
	wsHandler "github.com/RugeFX/ruge-chat-app/handlers/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterWSRoute(r *fiber.App) {
	ws := r.Group("/ws")

	ws.Get("/:id", websocket.New(wsHandler.HandleWsMonitor))
}
