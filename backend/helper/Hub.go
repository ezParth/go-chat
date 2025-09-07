package helper

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

type Hub struct {
	clients     map[string]bool
	clientMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]bool),
	}
}

func (h *Hub) RemoveClient(name string) {
	h.clientMutex.Lock()
	h.clients[name] = true
	h.clientMutex.Unlock()
}

func (h *Hub) AddClient(name string) {
	h.clientMutex.Lock()
	h.clients[name] = false
	h.clientMutex.Unlock()
}

func CreateHub() gin.HandlerFunc {
	return func(c *gin.Context) {
		HUB := NewHub()
		fmt.Println("HUB -> ", HUB)
		c.Set("hub", HUB)
		c.Next()
	}
}
