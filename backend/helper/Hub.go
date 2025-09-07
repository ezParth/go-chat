package helper

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients     map[string]bool
	conns       map[string]*websocket.Conn
	clientMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]bool),
		conns:   make(map[string]*websocket.Conn),
	}
}

func (h *Hub) AddClient(name string, conn *websocket.Conn) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	h.clients[name] = true
	h.conns[name] = conn
}

func (h *Hub) RemoveClient(name string) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	h.clients[name] = false
	delete(h.conns, name)
}

func CreateHub() gin.HandlerFunc {
	return func(c *gin.Context) {
		HUB := NewHub()
		fmt.Println("HUB created ->", HUB)
		c.Set("hub", HUB)
		c.Next()
	}
}

func (h *Hub) Broadcast(msg interface{}) {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	for name, conn := range h.conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error broadcasting to %s: %v\n", name, err)
		}
	}
}
