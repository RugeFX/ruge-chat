package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
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
			m.clients[c] = true
			jsonMessage, _ := json.Marshal(&fiber.Map{"Content": "/A new socket has connected. ", "SenderIP": c.conn.RemoteAddr().String()})
			m.sendMessage(jsonMessage, c)
		case c := <-m.unregister:
			if _, ok := m.clients[c]; ok {
				close(c.egr)
				delete(m.clients, c)
				jsonMessage, _ := json.Marshal(&fiber.Map{"Content": "/A socket has disconnected. ", "SenderIP": c.conn.RemoteAddr().String()})
				m.sendMessage(jsonMessage, c)
			}
		case message := <-m.broadcast:
			for conn := range m.clients {
				select {
				case conn.egr <- message:
				default:
					close(conn.egr)
					delete(m.clients, conn)
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
		c.conn.Close()
		delete(m.clients, c)

		log.Println("Client disconnected")
	}
}
