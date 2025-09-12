/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState } from "react"
import { useParams } from "react-router-dom"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"
import { getSocket } from "../servies/getSocket"
import { showError } from "../servies/toast"
import { groupApi } from "../api/group"

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
    // ws.current = new WebSocket("ws://localhost:8080/ws")
    ws.current = getSocket()

    ws.current.onopen = () => {
      console.log("✅ Connected to Group Chat WS")

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
        console.log("MESSAGE -> ", msg)
        if (!msg.event) return

        switch (msg.event) {
          case `Recieve-Message-${groupName}`:
          // case `Recieve-Message`:
            saveMessage(msg)
            break

          default:
            console.log("⚠️ Unhandled event:", msg.event)
        }
      } catch (err) {
        console.error("❌ Invalid WS message:", event.data)
      }
    }

    ws.current.onclose = () => {
      console.log("🔌 Disconnected from WS")
    }

    return () => {
      ws.current?.close()
    }
  }, [groupName, username, isAuthenticated])

  useEffect(() => {
    const fetchGroupChat = async () => {
      try {
        if(!groupName) {
          showError("GroupName not found")
          return
        }

        const res = await groupApi.getGroupChat(groupName)
        console.log("response of getGroupChat -> ", res)
  
        if (!res?.data.success) {
          showError("Cannot get group Chat")
        } else {
          console.log("chats -> ", res?.data.chats)
          setMessages(res?.data.chats)
        }
      } catch (err) {
        showError("Error fetching group chat")
        console.error(err)
      }
    }
  
    fetchGroupChat()
  }, [groupName])
  
  const saveMessage = (msg: WSMessage) => {
    setMessages((prev) => [
      ...prev,
      {
        sender: msg.user || "Unknown",
        message: String(msg.data),
        time: new Date().toLocaleTimeString(),
      },
    ])
  }

  const sendMessage = () => {
    if (!newMessage.trim() || !groupName || !username) return

    const payload: WSMessage = {
      event: `Send-Message`,
      room: groupName,
      user: username,
      data: newMessage,
    }

    ws.current?.send(JSON.stringify(payload))
    // setMessages((prev) => [
    //   ...prev,
    //   { sender: "You", message: newMessage, time: new Date().toLocaleTimeString() },
    // ])
    setNewMessage("")
  }

  if (!isAuthenticated) return <p>Please log in to chat</p>

  return (
    <div style={{ padding: "2rem" }}>
      <h2>💬 Group Chat: {groupName}</h2>
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
            <strong>{msg.sender}:</strong> {msg.message}{" "}
            <span style={{ fontSize: "0.75rem", color: "gray" }}>
              {msg.time}
            </span>
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
