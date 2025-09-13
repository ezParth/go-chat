/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState } from "react"
import { useParams } from "react-router-dom"
import { useSelector } from "react-redux"
import type { RootState } from "../store/store"
import { getSocket } from "../servies/getSocket"
import { showError, showInfo } from "../servies/toast"
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
  const [onlineUsers, setOnlineUsers] = useState<string[]>([])
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!isAuthenticated || !username || !groupName) return

    ws.current = getSocket()

    ws.current.onopen = () => {
      console.log("✅ Connected to Group Chat WS")

      const joinPayload: WSMessage = {
        event: "join",
        room: groupName,
        user: username,
      }
      ws.current?.send(JSON.stringify(joinPayload))

      const groupPayload: WSMessage = {
        event: "group-join",
        room: groupName,
        user: username
      }
      ws.current?.send(JSON.stringify(groupPayload))
    }

    ws.current.onmessage = (event: MessageEvent) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        if (!msg.event) return

        switch (msg.event) {
          case `Recieve-Message-${groupName}`:
            saveMessage(msg)
            break
          case `Group-Join-${groupName}`:
            showInfo(msg.data)
            setOnlineUsers((prev) => msg.user && !prev.includes(msg.user) ? [...prev, msg.user] : prev)
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
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        const leavePayload: WSMessage = {
          event: "leave",
          room: groupName,
          user: username,
        }
        ws.current.send(JSON.stringify(leavePayload))
      }
      ws.current?.close()
    }
  }, [groupName, username, isAuthenticated])

  useEffect(() => {
    const fetchGroupChat = async () => {
      try {
        if (!groupName) {
          showError("GroupName not found")
          return
        }
        const res = await groupApi.getGroupChat(groupName)
        if (res?.data.success) {
          setMessages(res?.data.chats)
        } else {
          showError("Cannot get group Chat")
        }
      } catch (err) {
        showError("Error fetching group chat")
      }
    }

    const fetchOnlineUsers = async () => {
      try {
        if (!groupName) {
          showError("GroupName Not Found!")
          return
        }
        const res = await groupApi.getUsersByGroupName(groupName)
        if (res?.data.success) {
          setOnlineUsers(res?.data.members)
        } else {
          showError("Cannot get online users")
        }
      } catch (error) {
        showError("Error fetching online users")
      }
    }

    fetchGroupChat()
    fetchOnlineUsers()
  }, [groupName])
  
  const saveMessage = (msg: WSMessage) => {
    setMessages((prev) => [
      ...prev,
      {
        sender: msg.user || "Unknown",
        message: String(msg.data),
        time: new Date().toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"}),
      },
    ])
  }

  const sendMessage = () => {
    if (!newMessage.trim() || !groupName || !username) return

    const payload: WSMessage = {
      event: "Send-Message",
      room: groupName,
      user: username,
      data: newMessage,
    }

    ws.current?.send(JSON.stringify(payload))
    setNewMessage("")
  }

  if (!isAuthenticated) return <p>Please log in to chat</p>

  return (
    <div style={{ display: "flex", height: "100vh", backgroundColor: "#ece5dd" }}>
      {/* Sidebar for online users */}
      <div
        style={{
          width: "220px",
          backgroundColor: "#ffffff",
          borderRight: "1px solid #ddd",
          padding: "1rem",
          overflowY: "auto",
        }}
      >
        <h3 style={{ marginBottom: "1rem", color: "#075E54" }}>🟢 Online</h3>
        {onlineUsers?.length > 0 ? (
          onlineUsers.map((user, idx) => (
            <div
              key={idx}
              style={{
                padding: "0.5rem",
                borderBottom: "1px solid #f0f0f0",
                fontWeight: user === username ? "bold" : "normal",
                color: user === username ? "#075E54" : "#333",
              }}
            >
              {user}
            </div>
          ))
        ) : (
          <p style={{ color: "gray" }}>No users online</p>
        )}
      </div>

      {/* Chat area */}
      <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
        <div
          style={{
            backgroundColor: "#075E54",
            color: "white",
            padding: "1rem",
            fontWeight: "bold",
          }}
        >
          💬 {groupName} — {username}
        </div>
  
        <div
          style={{
            flex: 1,
            overflowY: "auto",
            padding: "1rem",
            backgroundImage: "url('https://i.ibb.co/4pDNDk1/whatsapp-bg.png')",
            backgroundSize: "cover",
          }}
        >
          {messages.map((msg, idx) => {
            const isMe = msg.sender === username
            return (
              <div
                key={idx}
                style={{
                  display: "flex",
                  justifyContent: isMe ? "flex-end" : "flex-start",
                  marginBottom: "1rem",
                }}
              >
                <div
                  style={{
                    backgroundColor: isMe ? "#dcf8c6" : "white",
                    padding: "0.75rem 1rem",
                    borderRadius: "12px",
                    maxWidth: "65%",
                    boxShadow: "0 2px 4px rgba(0,0,0,0.15)",
                  }}
                >
                  {/* Sender */}
                  <div
                    style={{
                      fontWeight: "bold",
                      fontSize: "0.85rem",
                      textAlign: "center",
                      marginBottom: "0.25rem",
                      color: isMe ? "#075E54" : "#333",
                    }}
                  >
                    {msg.sender}
                  </div>

                  {/* Message */}
                  <div style={{ fontSize: "1rem", marginBottom: "0.4rem" }}>
                    {msg.message}
                  </div>

                  {/* Time */}
                  <div
                    style={{
                      fontSize: "0.8rem",
                      fontWeight: "bold",
                      color: "gray",
                      textAlign: "right",
                    }}
                  >
                    {msg.time}
                  </div>
                </div>
              </div>
            )
          })}
        </div>
  
        <div
          style={{
            display: "flex",
            padding: "0.5rem",
            backgroundColor: "#f0f0f0",
            borderTop: "1px solid #ddd",
          }}
        >
          <input
            type="text"
            placeholder="Type a message"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            style={{
              flex: 1,
              padding: "0.75rem",
              borderRadius: "20px",
              border: "1px solid #ccc",
              outline: "none",
              fontSize: "1rem",
            }}
          />
          <button
            onClick={sendMessage}
            style={{
              marginLeft: "0.5rem",
              padding: "0.75rem 1.25rem",
              borderRadius: "50%",
              border: "none",
              backgroundColor: "#075E54",
              color: "white",
              fontWeight: "bold",
              cursor: "pointer",
              fontSize: "1.2rem",
            }}
          >
            ➤
          </button>
        </div>
      </div>
    </div>
  )
}

export default GroupChat
