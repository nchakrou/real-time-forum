import { OnlineUsers, ChatUsers } from "../../components/pagesInit.js";
import { message } from "./messages.js";
import { chatHistory } from "./messages.js";
import {
  showNotification,
  storeNotification,
  renderStoredNotifications,
} from "../WebSocket/shownotification.js";
import { Popup } from "../../components/Popup.js";
import { updateUserList, formatTime } from "../../utils/chatUtils.js";
import { router } from "../Router.js";
import { ErrorPage } from "../../pages/Error.js";

export let ws;
export let currentChatUser = null;

export function setCurrentChatUser(username) {
  currentChatUser = username;
}

export function OpenWS() {
  return new Promise((resolve, reject) => {
    if (
      ws &&
      (ws.readyState === WebSocket.OPEN ||
        ws.readyState === WebSocket.CONNECTING)
    ) {
      resolve();
      return;
    }

    ws = new WebSocket("ws://localhost:8081/ws");

    ws.onopen = () => {
      console.log("Connected to WebSocket server");
      ws.send(JSON.stringify({ type: "get_notifications" }));
      
      resolve();
    };

    ws.onmessage = (event) => {
      let data;

      try {
        data = JSON.parse(event.data);
      } catch (e) {
        console.error("Invalid WS message:", e);
        return;
      }

      if (data.type === "notification") {
        storeNotification(data); 
      }

      switch (data.type) {
        case "typing":
          handleTypingStatus(data);
          break;

        case "online_users":
          OnlineUsers(data.users);
          break;

        case "chat_history":
          chatHistory(data);
          break;

        case "private_message":
          handlePrivateMessage(data);
          break;

        case "chat_users":
          ChatUsers(data.chat);
          break;

        case "notifications_history":
          if (Array.isArray(data.data)) {
            data.data.forEach((n) => storeNotification(n, false));
          }
          renderStoredNotifications();
          break;

        case "join":
          ws.send(JSON.stringify({ type: "get_chat_users" }));
          ws.send(JSON.stringify({ type: "online_users" }));
          break;

        case "user_offline": {
          const userItem = document.querySelector(
            `.user-item[data-username="${data.from}"]`
          );

          if (!userItem) break;

          if (window.location.pathname !== "/chat") {
            userItem.remove();
          } else {
            userItem.classList.remove("online");
          }
          break;
        }

        default:
          console.warn("Unhandled WS type:", data.type);
      }
    };

    ws.onerror = (error) => {
      Popup.show("Connection issue...");
    };

    ws.onclose = (event) => {
      console.log("CLOSE EVENT:", event);
      if (!event.wasClean) {
        ErrorPage("Server disconnected");
      }
    };
  }); 
} 

function handlePrivateMessage(data) {
  const currentChat = new URLSearchParams(window.location.search).get(
    "username"
  );

  if (data.to) {
    if (window.location.pathname === "/chat" && currentChat === data.to) {
      const chatViewport = document.getElementById("chat-viewport");
      if (!chatViewport) return;

      const div = document.createElement("div");
      div.classList.add("Mymessage");

      const h4 = document.createElement("h4");
      h4.textContent = "Me";

      const p = document.createElement("p");
      p.textContent = data.message;

      const timeSpan = document.createElement("span");
      timeSpan.classList.add("message-time");
      timeSpan.textContent = formatTime(data.CreatedAt);

      div.appendChild(h4);
      div.appendChild(p);
      div.appendChild(timeSpan);
      chatViewport.appendChild(div);
      chatViewport.scrollTop = chatViewport.scrollHeight;
    }
    return;
  }

  if (window.location.pathname === "/chat") {
    updateUserList(data.from);
  }

  if (currentChat === data.from) {
    message(data);
  } else {
    showNotification(data, true);
  }
}

function handleTypingStatus(data) {
  const currentChat = new URLSearchParams(window.location.search).get("username");
  if (window.location.pathname !== "/chat" || currentChat !== data.from) return;

  const typingIndicator = document.getElementById("typing-indicator");
  if (!typingIndicator) return;

  const isTyping = Boolean(data.is_typing ?? data.isTyping);
  const nameEl = typingIndicator.querySelector(".typing-name");

  if (isTyping) {
    if (nameEl) {
      nameEl.textContent = data.from;
    } else {
      typingIndicator.textContent = `${data.from} is typing...`;
    }
    typingIndicator.style.display = "flex";
  } else {
    if (nameEl) {
      nameEl.textContent = "";
    } else {
      typingIndicator.textContent = "";
    }
    typingIndicator.style.display = "none";
  }
}