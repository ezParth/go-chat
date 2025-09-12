/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState } from "react";
import { useSelector } from "react-redux";
import type { RootState } from "../store/store";
import { useNavigate } from "react-router-dom";
import { getSocket } from "../servies/getSocket";

interface WSMessage {
  event: string;
  room?: string;
  user?: string;
  data?: any;
}

export default function Chat() {
  const ws = useRef<WebSocket | null>(null);
  const [messages, setMessages] = useState<string[]>([]);
  const [message, setMessage] = useState<string>("");
  const { username, isAuthenticated } = useSelector((state: RootState) => state.auth)
  const [peopleOnline, setPeopleOnline] = useState<string[]>([])
  const nav = useNavigate()

  useEffect(() => {
    if (!isAuthenticated || !username) {
      alert("Please login");
      nav("/login");
    }
  }, [isAuthenticated, nav]);



  useEffect(() => {
    // ws.current = new WebSocket("ws://localhost:8080/ws");
    ws.current = getSocket()

    ws.current.onopen = () => {
      console.log("Client Connected Successfully");

      const data: WSMessage = {
        event: "join",
        user: username ?? "",
      };

      ws.current?.send(JSON.stringify(data));
    };

    ws.current.onmessage = (event: MessageEvent) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);

        if (!msg.event) {
          console.warn("Event cannot be empty");
          return;
        }

        switch (msg.event) {
          case "Recieve-Message":
            setMessages((prev) => [...prev, String(msg.data)]);
            break;

          default:
            console.log("Unhandled event:", msg.event);
            break;
        }
      } catch (err) {
        console.error("Invalid message received:", event.data);
      }
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  const handleMessage = () => {
    if (message.trim() === "") {
      console.warn("Message can't be empty");
      return;
    }

    const sendMessage: WSMessage = {
      event: "Message",
      room: "General",
      user: username ?? "",
      data: message,
    };

    ws.current?.send(JSON.stringify(sendMessage));
    setMessage("");
  };

  return (
    <div>
      <p>welcome - {username}</p>
      <p>
        {peopleOnline.map((people, index) => {
          return <li key={index}>{people}</li>
        })}
      </p>
      <div>
        <input
          type="text"
          placeholder="Enter a message..."
          value={message}
          onChange={(e) => setMessage(e.target.value)}
        />
        <button onClick={handleMessage}>Submit</button>
      </div>

      <div>
        <ul>
          {messages.map((msg, key) => (
            <li key={key}>{msg}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}
