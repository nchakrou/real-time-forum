import { OnlineUsers, ChatUsers } from "../../components/pagesInit.js";
import { message } from "./messages.js";
import { chatHistory } from "./messages.js";
import {
  showNotification,
  storeNotification,
} from "../WebSocket/shownotification.js";
import { Popup } from "../../components/Popup.js";
import { updateUserList } from "../../utils/chatUtils.js";

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
        case "chat_users":
          ChatUsers(data.chat);
          break;
        case "notifications_history":
          if (data.data && Array.isArray(data.data)) {
            data.data.forEach((n) => storeNotification(n, false));
          }
          break;
        case "join":
          if (window.location.pathname !== "/chat") {
            handleJoin(data);
          }
          break;
      }
    };

    ws.onclose = () => console.log("Disconnected from WebSocket server");
    ws.onerror = (error) => {
      Popup.show("Error connecting");
      reject(error);
    };
  });
}

function handlePrivateMessage(data) {
  const currentChat = new URLSearchParams(window.location.search).get(
    "username",
  );
  if (window.location.pathname === "/chat") {
    updateUserList(data.from);
  }
  if (currentChat === data.from) {
    message(data);
  } else {
    showNotification(data, true);
  }
}

function handleJoin(data) {
  const users = document.getElementsByClassName("online-users")[0];
  const msg = users.querySelector("p");
  let list;
  if (msg) {
    msg.remove();
    list = document.createElement("div");
    list.className = "list-users";
    users.appendChild(list);
  }
  const flag = document.querySelector(
    `.user-item[data-username="${data.from}"]`,
  );
  if (flag) return;

  const userItem = document.createElement("div");
  userItem.className = "user-item";
  userItem.dataset.username = data.from;
  userItem.addEventListener("click", () => {
    router(`/chat?username=${data.from}`);
  });

  const avatar = data.from.charAt(0).toUpperCase();
  userItem.innerHTML = `
    <div class="user-item-avatar">${avatar}</div>
    <div class="user-item-info">
      <span class="user-item-name">${data.from}</span>
    </div>`;
  if (!list) {
    list = document.getElementsByClassName("list-users")[0];
  }
  list.prepend(userItem);
}
