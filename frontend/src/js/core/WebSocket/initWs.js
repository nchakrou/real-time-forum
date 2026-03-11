import { OnlineUsers } from "../../components/pagesInit.js";
import { message } from "./messages.js";
import { chatHistory } from "./messages.js";
import { showNotification, storeNotification } from "../WebSocket/shownotification.js";

export let ws;
export let currentChatUser = null;

export function setCurrentChatUser(username) {
    currentChatUser = username;
}

export function OpenWS() {
    return new Promise((resolve, reject) => {
        ws = new WebSocket("ws://localhost:8081/ws");

        ws.onopen = () => {
            console.log("Connected to WebSocket server");
            ws.send(JSON.stringify({ type: "get_notifications" }));
            resolve();
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log("Message from server:", data);

            switch (data.type) {
                case "online_users":
                    OnlineUsers(data.users);
                    break;
                case "chat_history":
                    chatHistory(data);
                    break;
                case "private_message":
                    handlePrivateMessage(data);
                    break;
                case "notifications_history":
                    if (data.data && Array.isArray(data.data)) {
                        data.data.forEach(n => storeNotification(n, false));
                    }
                    break;
            }
        };

        ws.onclose = () => console.log("Disconnected from WebSocket server");
        ws.onerror = (error) => {
            console.error("WebSocket error:", error);
            reject(error);
        };
    });
}

function handlePrivateMessage(data) {
    const currentChat = new URLSearchParams(window.location.search).get("username");

    if (currentChat === data.from) {
        message(data);
    } else {
        showNotification(data, true);
    }
}