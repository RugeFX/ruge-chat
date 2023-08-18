package ws

import (
	"encoding/json"
	"fmt"
	"log"
)

type ClientMap map[*Client]bool

type Manager struct {
	clients    ClientMap
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewManager() *Manager {
	log.Println("New manager")
	return &Manager{
		clients:    make(ClientMap),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) Listen() {
	for {
		select {
		case c := <-m.register:
			m.addClient(c)
			jsonMessage, _ := json.Marshal(&MessageResponse{Type: "message", From: "Server", Body: fmt.Sprintf("%v has connected to the room", c.conn.RemoteAddr().String())})
			m.sendMessage(jsonMessage, c)
		case c := <-m.unregister:
			m.removeClient(c)
			jsonMessage, _ := json.Marshal(&MessageResponse{Type: "message", From: "Server", Body: fmt.Sprintf("%v has disconnected from the room", c.conn.RemoteAddr().String())})
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
