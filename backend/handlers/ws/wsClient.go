package ws

import (
	"encoding/json"
	"log"

	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
)

type Client struct {
	conn *websocket.Conn
	man  *Manager
	egr  chan []byte
}

func NewClient(c *websocket.Conn, m *Manager) *Client {
	return &Client{
		conn: c,
		man:  m,
		egr:  make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.man.unregister <- c
		c.conn.Close()
	}()

	for {
		mt, p, err := c.conn.ReadMessage()

		if err != nil {
			c.man.unregister <- c
			c.conn.Close()
			break
		}

		jsonMessage, _ := json.Marshal(&gin.H{"Content": string(p), "SenderIP": c.conn.RemoteAddr().String()})
		c.man.broadcast <- jsonMessage

		log.Println("mt :", mt, "p :", string(p))
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case m, y := <-c.egr:
			if !y {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed ", err)
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
				log.Println("Error sending msg: ", err)
			}
			log.Println("Message sent")
		}
	}
}
