import { OnlineUsers } from "../../pages/home.js"
export function OpenWS() {
    const ws = new WebSocket("ws://localhost:8081/ws")
    ws.onopen = () => {
        console.log("Connected to WebSocket server")
        ws.send(JSON.stringify({ type: "online_users" }))

    }
    ws.onmessage = (event) => {
        const res = JSON.parse(event.data)
        if (res.type === "online_users") {
            console.log("Online users:", res)
            OnlineUsers(res.users)
        }
    }
    ws.onclose = () => {
        console.log("Disconnected from WebSocket server")
    }
    ws.onerror = (error) => {
        console.error("WebSocket error:", error)
    }
}
