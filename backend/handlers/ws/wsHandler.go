package ws

import (
	"log"
	"net/http"

	"github.com/fasthttp/websocket"
	"github.com/go-playground/validator/v10"
)

type Connections map[*websocket.Conn]bool

type ErrorBody struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type MessageResponse struct {
	Type string `json:"type" validate:"required"`
	From string `json:"from" validate:"required"`
	Body string `json:"body" validate:"required,min=1,max=255"`
}

type ErrorResponse struct {
	Type   string       `json:"type" validate:"required"`
	Errors []*ErrorBody `json:"errors" validate:"required"`
}

var (
	conns             = make(Connections)
	register          = make(chan *websocket.Conn)
	unregister        = make(chan *websocket.Conn)
	validate          = validator.New()
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func HandleWs(w http.ResponseWriter, r *http.Request) {
	manager := NewManager()
	go manager.listen()

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, manager)
	manager.register <- client

	go client.readMessages()
	go client.writeMessages()
}
