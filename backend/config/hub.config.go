package helper

// we have to implement channels here later on and also message queues too
// right now this hub is not safe and closes randomly so, we have to make the use of RabbitMQ/kafka/redis here in future

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	conns       map[string]*websocket.Conn
	names       map[*websocket.Conn]string
	groups      map[string][]string
	clientMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		conns:  make(map[string]*websocket.Conn),
		names:  make(map[*websocket.Conn]string),
		groups: make(map[string][]string),
	}
}

func CreateHub() *Hub {
	HUB := NewHub()
	log.Println("HUB created Successfully")
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
	name := h.names[conn]
	delete(h.conns, name)
	delete(h.names, conn)
	h.clientMutex.Unlock()

	conn.Close()
	log.Println(":: Client Disconnected ", name)
}

// func (h *Hub) RemoveClient(conn *websocket.Conn) {
// 	h.clientMutex.Lock()
// 	defer h.clientMutex.Unlock()

// 	name := h.names[conn]
// 	if conn, ok := h.conns[name]; ok {
// 		conn.Close()
// 		delete(h.conns, name)
// 		delete(h.names, conn)
// 		log.Println(":: Client Disconneted ", name)
// 	}
// }

func (h *Hub) Broadcast(msg interface{}) {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	for name, conn := range h.conns {
		go func(name string, conn *websocket.Conn) {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Error broadcasting to %s: %v\n", name, err)
				go h.RemoveClient(conn)
			}
		}(name, conn)
	}
}

func (h *Hub) BroadcastToGroup(groupName string, msg interface{}) {
	h.clientMutex.RLock()
	defer h.clientMutex.RUnlock()

	members, exists := h.groups[groupName]
	if !exists {
		fmt.Printf("Group '%s' does not exist", groupName)
		return
	}

	for _, username := range members {
		if conn, ok := h.conns[username]; ok {
			go func(conn *websocket.Conn) {
				if err := conn.WriteJSON(msg); err != nil {
					log.Printf("Error sending to %s in group %s: %v", username, groupName, err)
					go h.RemoveClient(conn)
				}
			}(conn)
		}
	}
}

func (h *Hub) AddToGroup(groupName, username string) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	h.groups[groupName] = append(h.groups[groupName], username)
}

func (h *Hub) RemoveFromGroup(groupName, username string) {
	h.clientMutex.Lock()
	defer h.clientMutex.Unlock()

	members := h.groups[groupName]
	for i, member := range members {
		if member == username {
			h.groups[groupName] = append(members[:i], members[i+1:]...)
			break
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
