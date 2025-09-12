package helper

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	conns       map[string]*websocket.Conn
	names       map[*websocket.Conn]string
	clientMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		conns: make(map[string]*websocket.Conn),
	}
}

func CreateHub() *Hub {
	HUB := NewHub()
	log.Println("HUB created ->", HUB)
	return HUB
}

func (h *Hub) AddClient(name string, conn *websocket.Conn) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	h.conns[name] = conn
	h.names[conn] = name
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	name := h.names[conn]
	if conn, ok := h.conns[name]; ok {
		conn.Close()
		delete(h.conns, name)
		delete(h.names, conn)
		log.Println(":: Client Disconneted ", name)
	}
}

func (h *Hub) Broadcast(msg interface{}) {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	for name, conn := range h.conns {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error broadcasting to %s: %v\n", name, err)
			// Better to remove bad connection
			go h.RemoveClient(conn)
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
