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
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!isAuthenticated || !username || !groupName) return

    ws.current = getSocket()

    ws.current.onopen = () => {
      console.log("✅ Connected to Group Chat WS")

      // 👉 JOIN event
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
        console.log("MESSAGE -> ", msg)
        if (!msg.event) return

        switch (msg.event) {
          case `Recieve-Message-${groupName}`:
            saveMessage(msg)
            break
          case `Group-Join-${groupName}`:
            showInfo(msg.data)
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
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        height: "100vh",
        backgroundColor: "#ece5dd", // WhatsApp background
      }}
    >
      {/* Header */}
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
  
      {/* Messages */}
      <div
        style={{
          flex: 1,
          overflowY: "auto",
          padding: "1rem",
          backgroundImage:
            "url('https://i.ibb.co/4pDNDk1/whatsapp-bg.png')", // WhatsApp-like wallpaper
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
                marginBottom: "0.5rem",
              }}
            >
              <div
                style={{
                  backgroundColor: isMe ? "#dcf8c6" : "white", // green for me, white for others
                  padding: "0.5rem 1rem",
                  borderRadius: "10px",
                  maxWidth: "60%",
                  boxShadow: "0 1px 1px rgba(0,0,0,0.1)",
                  position: "relative",
                }}
              >
                <span style={{ fontWeight: "bold", fontSize: "0.85rem" }}>
                  {msg.sender}
                </span>
                <div>{msg.message}</div>
                <span
                  style={{
                    fontSize: "0.7rem",
                    color: "gray",
                    position: "absolute",
                    bottom: "2px",
                    right: "6px",
                  }}
                >
                  {msg.time}
                </span>
              </div>
            </div>
          )
        })}
      </div>
  
      {/* Input */}
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
          }}
        >
          ➤
        </button>
      </div>
    </div>
  )
  

  // return (
  //   <div style={{ padding: "2rem" }}>
  //     <h2>💬 Group Chat: {groupName} Username: {username}</h2>
  //     <div
  //       style={{
  //         border: "1px solid gray",
  //         height: "300px",
  //         overflowY: "scroll",
  //         padding: "1rem",
  //         marginBottom: "1rem",
  //       }}
  //     >
  //       {messages.map((msg, idx) => (
  //         <p key={idx}>
  //           <strong>{msg.sender}:</strong> {msg.message}{" "}
  //           <span style={{ fontSize: "0.75rem", color: "gray" }}>
  //             {msg.time}
  //           </span>
  //         </p>
  //       ))}
  //     </div>
  //     <input
  //       type="text"
  //       placeholder="Type a message"
  //       value={newMessage}
  //       onChange={(e) => setNewMessage(e.target.value)}
  //     />
  //     <button onClick={sendMessage} style={{ marginLeft: "1rem" }}>
  //       Send
  //     </button>
  //   </div>
  // )
}

export default GroupChat
