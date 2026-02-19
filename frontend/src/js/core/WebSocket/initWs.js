import { OnlineUsers } from "../../components/pagesInit.js"
import { message } from "./messages.js"
import { chatHistory } from "./messages.js"
export let ws
export function OpenWS() {
    return new Promise((resolve, reject) => {
        ws = new WebSocket("ws://localhost:8081/ws")
        ws.onopen = () => {
            console.log("Connected to WebSocket server")
            resolve()
        }
        ws.onmessage = (event) => {
            console.log("Message from server:", event.data)
            const data = JSON.parse(event.data)
            if (data.type === "online_users") {
                OnlineUsers(data.users)
                console.log(data.users)
            } else if (data.type === "chat_history") {
                chatHistory(data)
            } else if (data.type === "message") {
                message(data)
            }
        }
        ws.onclose = () => {
            console.log("Disconnected from WebSocket server")
        }
        ws.onerror = (error) => {
            console.error("WebSocket error:", error)
            reject(error)
        }
    })
}