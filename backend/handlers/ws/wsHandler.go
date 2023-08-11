package wsHandler

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Connections map[*websocket.Conn]bool

type Message struct {
	Body string `json:"body"`
}

var (
	conns      = make(map[*websocket.Conn]bool)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
)

func HandleWsMonitor(c *websocket.Conn) {
	go SocketListener()
	defer func() {
		unregister <- c
		c.Close()
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

			for ws, v := range conns {
				if !v || ws.Params("id") != c.Params("id") {
					continue
				}

				var jsonMsg Message
				if err := json.Unmarshal(message, &jsonMsg); err != nil {
					log.Println("unmarshal error:", err)
				}

				go func(ws *websocket.Conn) {
					if err := ws.WriteJSON(jsonMsg); err != nil {
						log.Println("write error:", err)
					}
				}(ws)
			}
		} else {
			log.Println("received message of type:", messageType)
		}
	}
}

func SocketListener() {
	for {
		select {
		case c := <-register:
			log.Println("case c := <-register")

			conns[c] = true

			log.Println("Appended connection")
			log.Printf("Connection Params: %s\n", c.Params("id"))
		case c := <-unregister:
			delete(conns, c)
			c.Close()
			log.Println("Closed connection")
		}
	}
}
