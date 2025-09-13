package ws

import (
	helper "backend/config"
	"backend/controllers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	User  string `json:"user"`
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

func handleMessageEvent(conn *websocket.Conn, msg WSMessage, mt int, hub *helper.Hub) {
	fmt.Println("HANDLING MESSAGE EVENT")
	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Error in Type handling of data")
	}

	var Message WSMessage = WSMessage{
		Event: "Recieve-Message",
		Room:  "General",
		User:  msg.User,
		Data:  data,
	}

	hub.PrintHub()

	hub.Broadcast(Message)
}

func handleMessageEventForRoom(_ *websocket.Conn, msg WSMessage, _ int) {
	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Error in handling Type of data in handleMessageEventForRoom")
	}

	fmt.Println("data -> ", data)
}

func handleGroupChat(msg WSMessage, _ int, hub *helper.Hub) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	fmt.Println("Message from group ", msg.Room, " -> ", msg)
	if msg.Room == "" {
		log.Println("Group name (room) cannot be empty")
		return
	}
	if msg.User == "" {
		log.Println("User cannot be empty in group chat")
		return
	}

	fmt.Println("PrintingHub")
	hub.PrintHub()
	fmt.Println("\nDone Printing Hub")
	// 1. Convert data to string (message text)
	text, ok := msg.Data.(string)
	if !ok {
		log.Println("Invalid message data type in handleGroupChat")
		return
	}

	err := controllers.SaveGroupChatLogic(ctx, msg.User, msg.Room, text)
	if err != nil {
		fmt.Println("Error saving the message", err)
	}

	if err == nil {
		fmt.Println("GroupChat saved successfully")
	}

	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Wrong data type of msg.Data")
	}

	EventName := "Recieve-Message" + "-" + msg.Room
	var Message WSMessage = WSMessage{
		Event: EventName,
		// Event: "Recieve-Message",
		// Room: "General",
		Room: msg.Room,
		User: msg.User,
		Data: data,
	}

	hub.PrintHub()

	hub.Broadcast(Message)
}

func handleCreateRoom(conn *websocket.Conn, msg WSMessage, mt int) {

}

func handleJoinRoom(conn *websocket.Conn, msg WSMessage, mt int) {
	if msg.Room == "" {
		fmt.Println("Room cannot be empty bith")
	}
}

func handleJoin(conn *websocket.Conn, msg WSMessage, mt int, hub *helper.Hub) {
	// helper.User[msg.User] = conn
	hub.AddClient(msg.User, conn)
	fmt.Println("Client Connected: ", msg.User)
}

func handleClose(conn *websocket.Conn, msg WSMessage, mt int, hub *helper.Hub) {
	hub.RemoveClient(conn)
	fmt.Println("Client Disconnected: ", msg.User)
}

func handleGroupJoin(conn *websocket.Conn, msg WSMessage, mt int, hub *helper.Hub) {
	if msg.Room == "" {
		fmt.Println("Room can't be empty dude")
		return
	}

	Message := &WSMessage{
		Event: "Group-Join-" + msg.Room,
		User:  msg.User,
		Data:  msg.User + " Joined the Group",
		Room:  msg.Room,
	}

	hub.Broadcast(Message)
}

func WsHandler(c *gin.Context) {
	value, exists := c.Get("hub")

	hub, ok := value.(*helper.Hub)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid hub type"})
		return
	}

	if !exists {
		c.JSON(500, gin.H{"error": "hub not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error in websocket upgrader -> ", err)
	}

	fmt.Println("conn -> ", conn.LocalAddr())

	// var initialUser = "unknown"
	// hub.AddClient(initialUser, conn)
	defer func() {
		hub.RemoveClient(conn)
		conn.Close()
	}()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Client disconnected:", err)
			} else {
				log.Println("Unexpected websocket error:", err)
			}
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

		// changed: update initialUser when first message provides username
		// if initialUser == "unknown" && Message.User != "" {
		// 	initialUser = Message.User
		// }

		switch Message.Event {
		case "Message":
			handleMessageEvent(conn, Message, mt, hub)
		case "Room-Message":
			handleMessageEventForRoom(conn, Message, mt)
		case "join":
			handleJoin(conn, Message, mt, hub)
		case "Send-Message":
			handleGroupChat(Message, mt, hub)
		case "leave":
			handleClose(conn, Message, mt, hub)
		case "group-join":
			handleGroupJoin(conn, Message, mt, hub)
		}
	}
}

// func WsHandler(c *gin.Context) {
// 	value, exists := c.Get("hub")

// 	hub, ok := value.(*helper.Hub)
// 	if !ok {
// 		c.JSON(500, gin.H{"error": "invalid hub type"})
// 		return
// 	}

// 	if !exists {
// 		c.JSON(500, gin.H{"error": "hub not found"})
// 		return
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		// log.Fatal("Error in websocket upgrader -> ", err)

// 	}

// 	fmt.Println("conn -> ", conn.LocalAddr())

// 	defer conn.Close()

// 	for {
// 		mt, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error in Reading Messages -> ", err)
// 			break
// 		}

// 		fmt.Println("MESSAGE -> ", string(msg))

// 		var Message WSMessage
// 		if err = json.Unmarshal(msg, &Message); err != nil {
// 			log.Println("Invalid message format:", err)
// 			continue
// 		}
// 		if Message.Event == "" {
// 			log.Println("Message.Event cannot be empty")
// 			continue
// 		}

// 		hub.AddClient(Message.User, conn)
// 		defer hub.RemoveClient(Message.User)

// 		switch Message.Event {
// 		case "Message":
// 			handleMessageEvent(conn, Message, mt, hub)
// 		case "Room-Message":
// 			handleMessageEventForRoom(conn, Message, mt)
// 		case "join":
// 			handleJoin(conn, Message, mt, hub)
// 		case "/:groupName/message":
// 			handleGroupChat(conn, Message, mt, hub)
// 		}
// 	}
// }
