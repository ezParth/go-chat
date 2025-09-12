package helper

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	conns       map[string]*websocket.Conn
	clientMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) AddClient(name string, conn *websocket.Conn) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	h.conns[name] = conn
}

func (h *Hub) RemoveClient(name string) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	if conn, ok := h.conns[name]; ok {
		conn.Close()
		delete(h.conns, name)
	}
}

func (h *Hub) Broadcast(msg interface{}) {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	for name, conn := range h.conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error broadcasting to %s: %v\n", name, err)
			// Better to remove bad connection
			go h.RemoveClient(name)
		}
	}
}

func (h *Hub) PrintHub() {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	for name := range h.conns {
		log.Println("Connected user ->", name)
	}
}
