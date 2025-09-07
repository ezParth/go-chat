package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Event string `json:"event"`
	Room  string `json:"room,omitempty"`
	User  string `json:"user,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type RoomMember struct {
	Username string `json:"username"`
	Admin    bool   `json:"Admin"`
}

type Admin struct {
	Username string `json:"username"`
	RoomName string `json:"roomname"`
}

var Room = make(map[RoomMember]*websocket.Conn)

var User = make(map[string]*websocket.Conn)

func handleMessageEvent(conn *websocket.Conn, msg WSMessage, mt int) {
	fmt.Println("HANDLING MESSAGE EVENT")
	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Error in Type handling of data")
	}

	var Message WSMessage = WSMessage{
		Event: "Recieve-Message",
		Room:  "General",
		User:  "Alice",
		Data:  data,
	}

	if err := conn.WriteJSON(Message); err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}

}

func handleMessageEventForRoom(conn *websocket.Conn, msg WSMessage, mt int) {
	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Error in handling Type of data in handleMessageEventForRoom")
	}

	fmt.Println("data -> ", data)
}

func handleCreateRoom(conn *websocket.Conn, msg WSMessage, mt int) {

}

func handleJoinRoom(conn *websocket.Conn, msg WSMessage, mt int) {
	if msg.Room == "" {
		fmt.Println("Room cannot be empty bith")
	}
}

func handleJoin(conn *websocket.Conn, msg WSMessage, mt int) {
	User[msg.User] = conn
}

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error in websocket upgrader -> ", err)
	}

	fmt.Println("conn -> ", conn.LocalAddr())

	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error in Reading Messages -> ", err)
			break
		}

		fmt.Println("MESSAGE -> ", string(msg))

		var Message WSMessage
		if err = json.Unmarshal(msg, &Message); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}
		if Message.Event == "" {
			log.Println("Message.Event cannot be empty")
			continue
		}

		switch Message.Event {
		case "Message":
			handleMessageEvent(conn, Message, mt)

		case "Room-Message":
			handleMessageEventForRoom(conn, Message, mt)
		case "join":
			handleJoin(conn, Message, mt)
		}
	}
}
