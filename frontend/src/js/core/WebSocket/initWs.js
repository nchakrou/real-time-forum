import { OnlineUsers } from "../../components/pagesInit.js"
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