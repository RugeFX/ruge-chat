package ws

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
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
	conns      = make(Connections)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
	validate   = validator.New()
)

func HandleWs(c *websocket.Conn) {
	manager := NewManager()

	client := NewClient(c, manager)
	manager.register <- client

	go client.readMessages()
	go client.writeMessages()
}

func _HandleWs(c *websocket.Conn) {
	log.Println(c.Locals("allowed"))
	go _SocketListener()
	defer func() {
		_broadcastUserLeftMessage(c)
		unregister <- c
	}()

	register <- c
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return
		}
		if messageType == websocket.TextMessage {
			log.Println("got textmessage:", string(message))

			_broadcastJsonToChannel(c, message)
		} else {
			log.Println("received message of type:", messageType)
		}
	}
}

func _SocketListener() {
	for {
		select {
		case c := <-register:
			log.Println("case c := <-register")

			conns[c] = true
			remote := c.RemoteAddr().String()

			for ws, v := range conns {
				if !v || ws.Params("id") != c.Params("id") {
					continue
				}

				go func(ws *websocket.Conn) {
					if err := ws.WriteJSON(MessageResponse{Type: "message", From: "Server", Body: fmt.Sprintf("%v has joined the channel", remote)}); err != nil {
						log.Println("write error:", err)
					}
				}(ws)
			}

			log.Println("Appended connection")
			log.Printf("Connection Params: %s\n", c.Params("id"))
		case c := <-unregister:
			delete(conns, c)
			c.Close()
			log.Println("Closed connection")
		}
	}
}

func _broadcastJsonToChannel(c *websocket.Conn, m []byte) {
	var errors []*ErrorBody
	var jsonMsg MessageResponse
	if err := json.Unmarshal(m, &jsonMsg); err != nil {
		log.Println("unmarshal error:", err, ", message: ", m)
		return
	}

	jsonMsg.Type = "message"
	jsonMsg.From = c.RemoteAddr().String()

	err := validate.Struct(jsonMsg)

	if err != nil {
		log.Println("struct validation error")
		for _, err := range err.(validator.ValidationErrors) {
			var er ErrorBody
			er.Field = err.Field()
			er.Tag = err.Tag()
			er.Value = err.Param()
			errors = append(errors, &er)
		}
		c.WriteJSON(ErrorResponse{
			Type:   "error",
			Errors: errors,
		})
		return
	}

	for ws, v := range conns {
		if !v || ws.Params("id") != c.Params("id") {
			continue
		}

		go func(ws *websocket.Conn) {
			if err := ws.WriteJSON(jsonMsg); err != nil {
				log.Println("bj write error:", err)
			}
		}(ws)
	}
}

func _broadcastUserLeftMessage(c *websocket.Conn) {
	remote := c.RemoteAddr().String()
	for ws, v := range conns {
		if !v || ws.Params("id") != c.Params("id") {
			continue
		}
		// TODO : fix read and write error
		go func(ws *websocket.Conn) {
			if err := ws.WriteJSON(MessageResponse{Type: "message", From: "Server", Body: fmt.Sprintf("%v has left the channel", remote)}); err != nil {
				log.Println("bul write error:", err)
			}
		}(ws)
	}
}
