package ws

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type ClientMap map[*Client]bool

type Manager struct {
	clients    ClientMap
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(ClientMap),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) listen() {
	for {
		select {
		case c := <-m.register:
			m.addClient(c)
			jsonMessage, _ := json.Marshal(&gin.H{"Content": "/A new socket has connected. ", "SenderIP": c.conn.RemoteAddr().String()})
			m.sendMessage(jsonMessage, c)
		case c := <-m.unregister:
			m.removeClient(c)
			jsonMessage, _ := json.Marshal(&gin.H{"Content": "/A socket has disconnected. ", "SenderIP": c.conn.RemoteAddr().String()})
			m.sendMessage(jsonMessage, c)
		case message := <-m.broadcast:
			for conn := range m.clients {
				select {
				case conn.egr <- message:
				default:
					conn.man.removeClient(conn)
				}
			}
		}
	}
}

func (m *Manager) sendMessage(message []byte, ignore *Client) {
	for c := range m.clients {
		if c != ignore {
			c.egr <- message
		}
	}
}

func (m *Manager) addClient(c *Client) {
	log.Println("New client connected")

	m.clients[c] = true
}

func (m *Manager) removeClient(c *Client) {
	if _, y := m.clients[c]; y {
		close(c.egr)
		c.conn.Close()
		delete(m.clients, c)
		log.Println("Client disconnected")
	}
}
