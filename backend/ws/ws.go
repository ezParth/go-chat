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

func handleMessageEvent(conn *websocket.Conn, msg WSMessage, mt int) {
	data, ok := msg.Data.(string)
	if !ok {
		fmt.Println("Error in Type handling of data")
	}

	conn.WriteMessage(mt, []byte(data))
}

func handleMessageEventForRoom(conn *websocket.Conn, msg WSMessage, mt int) {

}

func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Error in websocket upgrader -> ", err)
	}

	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Error in Reading Messages -> ", err)
		}

		var Message WSMessage
		if err = json.Unmarshal(msg, &Message); err != nil {
			log.Fatal("Error in Unmarshelling Message -> ", err)
		}

		if Message.Event == "" {
			log.Fatal("Message.Event Cannot be nil or Empty -> ", Message.Event)
		}

		switch Message.Event {
		case "Message":
			handleMessageEvent(conn, Message, mt)

		case "Room-Message":
			handleMessageEventForRoom(conn, Message, mt)
		}
	}
}
