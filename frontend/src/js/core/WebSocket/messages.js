import { chatStates } from "../chatStates.js";
import { updateUserList, formatTime } from "../../utils/chatUtils.js";

export function message(data) {
  const currentChat = new URLSearchParams(window.location.search).get(
    "username",
  );

  if (window.location.pathname === "/chat") {
    updateUserList(data.from);
  }
  if (window.location.pathname === "/chat" && currentChat === data.from) {
    const chatViewport = document.getElementById("chat-viewport");
    if (!chatViewport) return;
    const div = document.createElement("div");
    div.classList.add("message");
    const h4 = document.createElement("h4");
    h4.textContent = data.from;
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
}


export function chatHistory(data) {
  const chatViewport = document.getElementById("chat-viewport");
  const isFirstLoad = chatStates.lastID === 0;

  if (isFirstLoad) {
    chatViewport.innerHTML = "";
  }

  if (data.Messages && data.Messages.length > 0) {

    chatStates.lastID = data.Messages[data.Messages.length - 1].id;

    if (data.Messages.length < 10) {
      chatStates.isEnd = true;
    }
    data.Messages.forEach((message) => {
      const messageDiv = document.createElement("div");
      const target = new URLSearchParams(window.location.search).get(
        "username",
      );

      const sender = document.createElement("h4");

      sender.textContent = message.from;

      const p = document.createElement("p");
      p.textContent = message.message;

      const timeSpan = document.createElement("span");
      timeSpan.classList.add("message-time");
      timeSpan.textContent = formatTime(message.CreatedAt);

      messageDiv.appendChild(sender);
      messageDiv.appendChild(p);
      messageDiv.appendChild(timeSpan);

      if (message.from === target) {
        messageDiv.classList.add("message");
      } else {
        messageDiv.classList.add("Mymessage");
        sender.textContent = "Me";
      }
      chatViewport.prepend(messageDiv);
    });

    if (isFirstLoad) {
      chatViewport.scrollTop = chatViewport.scrollHeight;
    }
  } else if (
    data.Messages &&
    data.Messages.length === 0 &&
    chatStates.lastID === 0
  ) {
    chatViewport.innerHTML = `
        <div class="empty-chat-state">
            <div class="empty-chat-icon">
                <img src="/frontend/src/assets/chat.svg" alt="Chat Icon">
            </div>
            <h3>No messages yet</h3>
            <p>Start a conversation</p>
        </div>
        `;
  } else {
    chatStates.isEnd = true;
  }
}
