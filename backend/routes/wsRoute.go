package routes

import (
	"log"
	"net/http"

	wsHandler "github.com/RugeFX/ruge-chat-app/handlers/ws"
	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func RegisterWSRoute(r *gin.Engine) {
	ws := r.Group("/ws")

	manager := wsHandler.NewManager()
	go manager.Listen()

	ws.GET("/:id", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		wsHandler.HandleWs(conn, manager)
	}))
}
