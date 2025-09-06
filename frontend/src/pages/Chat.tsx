/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useState } from "react";

interface WSMessage {
    Event: string
    Room?: string
    User?: string
    Data?: any
}

export default function Chat() {

    const ws = useRef<WebSocket | null>(null)
    const [messages, setMessages] = useState<Array<string>>([])
    const [message, setMessage] = useState<string>("")

    useEffect(() => {
        ws.current = new WebSocket("http://localhost:8080/ws")

        ws.current.onopen = () => {
            console.log("Client Connected Successfully")
            
            const data: WSMessage = {
                Event: "join",
                User: "alice",
            }

            ws.current?.send(JSON.stringify(data))
        }

        ws.current.onmessage = (event: any) => {
            const msg: WSMessage = JSON.parse(event)
            console.log(msg)

            if(!msg.Event) {
                console.log("Event Cannot be Empty")
                return
            }

            switch (msg.Event) {
                case "Message":
                    setMessages((prev) => [...prev, String(msg?.Data)])
                    break;
            
                default:
                    break;
            }
        }
    }, [])

    const handleMessage = () => {
        if(message == "") {
            console.log("Message Can't be Empty")
            return
        }

        const sendMessage: WSMessage = {
            Event: "Message",
            Room: "General",
            User: "Alice",
            Data: message
        }

        ws.current?.send(JSON.stringify(sendMessage))
    }

    return (
        <div>
            <div className="">
                <input type="text" placeholder="Enter a message..." value={message} onChange={(e) => setMessage(e.target.value)} />
                <button onChange={handleMessage}>Submit</button>
            </div>

            <div>
                {messages.map((msg, key) => (
                    <li key={key}>{msg}</li>
                ))}
            </div>
        </div>
    )
}
