package routes

import (
	wsHandler "github.com/RugeFX/ruge-chat-app/handlers/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterWSRoute(r *fiber.App) {
	ws := r.Group("/ws")

	ws.Use(func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.Get("/:id", websocket.New(wsHandler.HandleWs))

	// ws.Get("/:id", websocket.New(wsHandler.HandleWs, websocket.Config{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }))
}
