import { states } from "../Listeners/postListners.js";
import { showNotification } from "../WebSocket/shownotification.js";

export function message(data) {
  const currentChat = new URLSearchParams(window.location.search).get(
    "username",
  );

  if (window.location.pathname === "/chat" && currentChat === data.from) {
    const chatViewport = document.getElementById("chat-viewport");
    if (!chatViewport) return;

    const div = document.createElement("div");
    div.classList.add("message");

    const h4 = document.createElement("h4");
    h4.textContent = data.from;

    const p = document.createElement("p");
    p.textContent = data.message;

    div.appendChild(h4);
    div.appendChild(p);
    chatViewport.appendChild(div);
    chatViewport.scrollTop = chatViewport.scrollHeight;
  } else {
    showNotification(data);
  }
}

export function chatHistory(data) {
  console.log("end", states.isEnd);

  const chatViewport = document.getElementById("chat-viewport");
  if (states.offset === 0) {
    chatViewport.innerHTML = "";
  }
  states.offset += 10;
  if (data.Messages) {
    if (data.Messages.length < 10) {
      states.isEnd = true;
    }
    data.Messages.forEach((message) => {
      const messageDiv = document.createElement("div");
      const target = new URLSearchParams(window.location.search).get(
        "username",
      );
      console.log("hadi", message.from, target);
      const p = document.createElement("p");
      const sender = document.createElement("h4");

      sender.textContent = message.from;
      p.textContent = message.message;
      messageDiv.appendChild(sender);
      messageDiv.appendChild(p);
      if (message.from === target) {
        messageDiv.classList.add("message");
      } else {
        messageDiv.classList.add("Mymessage");
      }
      chatViewport.prepend(messageDiv);
    });

    if (states.offset === 10) {
      chatViewport.scrollTop = chatViewport.scrollHeight;
    }
  } else if (
    data.Messages &&
    data.Messages.length === 0 &&
    states.offset === 0
  ) {
    chatViewport.innerHTML = `
        <div class="empty-chat-state">
            <div class="empty-chat-icon">
                <img src="src/assets/chat.svg" alt="Chat Icon">
            </div>
            <h3>No messages yet</h3>
            <p>Start a conversation</p>
        </div>
        `;
  } else {
    states.isEnd = true;
  }
}
