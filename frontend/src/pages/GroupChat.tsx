/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState } from "react"
import { useParams } from "react-router-dom"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"

interface WSMessage {
  event: string
  room?: string
  user?: string
  data?: any
}

interface ChatMessage {
  sender: string
  message: string
  time?: string
}

const GroupChat = () => {
  const { groupName } = useParams<{ groupName: string }>()
  const { username, isAuthenticated } = useSelector((state: RootState) => state.auth)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [newMessage, setNewMessage] = useState("")
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!isAuthenticated || !username || !groupName) return

    // Connect WebSocket
    ws.current = new WebSocket("ws://localhost:8080/ws")

    ws.current.onopen = () => {
      console.log("Connected to Group Chat WS")

      // Join the group room
      const joinPayload: WSMessage = {
        event: "join",
        room: groupName,
        user: username,
      }
      ws.current?.send(JSON.stringify(joinPayload))
    }

    ws.current.onmessage = (event: MessageEvent) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        if (!msg.event) return

        switch (msg.event) {
          case "Recieve-Message":
            setMessages((prev) => [
              ...prev,
              { sender: msg.user || "Unknown", message: String(msg.data) },
            ])
            break

          default:
            console.log("Unhandled event:", msg.event)
        }
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      } catch (err) {
        console.error("Invalid WS message:", event.data)
      }
    }

    return () => {
      ws.current?.close()
    }
  }, [groupName, username, isAuthenticated])

  const sendMessage = () => {
    if (!newMessage.trim() || !groupName || !username) return

    const payload: WSMessage = {
      event: "Message",
      room: groupName,
      user: username,
      data: newMessage,
    }

    ws.current?.send(JSON.stringify(payload))
    setMessages([...messages, { sender: "You", message: newMessage }])
    setNewMessage("")
  }

  if (!isAuthenticated) return <p>Please log in to chat</p>

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Group Chat: {groupName}</h2>
      <div
        style={{
          border: "1px solid gray",
          height: "300px",
          overflowY: "scroll",
          padding: "1rem",
          marginBottom: "1rem",
        }}
      >
        {messages.map((msg, idx) => (
          <p key={idx}>
            <strong>{msg.sender}:</strong> {msg.message}
          </p>
        ))}
      </div>
      <input
        type="text"
        placeholder="Type a message"
        value={newMessage}
        onChange={(e) => setNewMessage(e.target.value)}
      />
      <button onClick={sendMessage} style={{ marginLeft: "1rem" }}>
        Send
      </button>
    </div>
  )
}

export default GroupChat
