import { states } from "../Listeners/postListners.js";
export function message(data) {
  if (window.location.search.includes(data.from)) {
    const chatViewport = document.getElementById("chat-viewport");
    const message = document.createElement("div");
    message.classList.add("message");
    const p = document.createElement("p");
    p.textContent = data.message;
    const sender = document.createElement("h4");
    sender.textContent = data.from;
    message.appendChild(sender);
    message.appendChild(p);
    chatViewport.appendChild(message);
  } else {
    MessageNotification(data.from);
  }
}
function MessageNotification(username) {
  const notificationBar = document.getElementById("notification-list");
  let notification = notificationBar.querySelector(`[data-user="${username}"]`);
  const badge = document.getElementById("notification-badge");
  if (window.getComputedStyle(badge).display === "none") {
    badge.style.display = "flex";
  }
  badge.textContent = parseInt(badge.textContent) + 1;
  if (notification) {
    const counter = ++notification.dataset.counter;
    notification.textContent = `You have ${counter} message from ${username}`;
  } else {
    notification = document.createElement("div");
    notification.classList.add("notification");
    notification.dataset.user = username;
    notification.dataset.counter = 1;
    notification.textContent = `You have 1 message from ${username}`;
    notificationBar.appendChild(notification);
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
  } else if (states.offset === 0) {
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
