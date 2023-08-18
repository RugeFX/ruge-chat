package ws

import (
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
	validate = validator.New()
)

func HandleWs(c *websocket.Conn, m *Manager) {
	client := NewClient(c, m)
	m.register <- client

	go client.readMessages()
	go client.writeMessages()
}
