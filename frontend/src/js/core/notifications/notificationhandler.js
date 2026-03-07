import { router } from "../Router.js";
import { ws } from "../WebSocket/initWs.js";

let notifications = {}; 

export function addNotification(data) {
  const username = data.from;
  const list = document.getElementById("notification-list");
  const badge = document.getElementById("notification-badge");

  if (notifications[username]) {
    notifications[username].counter++;
    notifications[username].element.querySelector(".notif-count").textContent =
      notifications[username].counter;

    notifications[username].element.querySelector(".notif-message").textContent =
      data.message;

    updateBadge();
    return;
  }

  const item = document.createElement("div");
  item.classList.add("notification-item");
  item.innerHTML = `
    <div class="notif-avatar">${username.charAt(0).toUpperCase()}</div>
    <div class="notif-content">
      <strong>${username}</strong>
      <p class="notif-message">${data.message}</p>
      <span class="notif-count">1</span> new
    </div>
    <div class="notif-actions">
      <button class="notif-reply">Reply</button>
      <button class="notif-close">✕</button>
    </div>
  `;

  item.addEventListener("click", (e) => {
    if (e.target.classList.contains("notif-reply")) return;
    if (e.target.classList.contains("notif-close")) return;

    removeNotification(username);
    router(`/chat?username=${username}`);
  });

  item.querySelector(".notif-reply").onclick = (e) => {
    e.stopPropagation();
    const reply = prompt(`Reply to ${username}`);
    if (!reply) return;
    ws.send(
      JSON.stringify({
        type: "message",
        target: username,
        message: reply,
      }),
    );

    removeNotification(username);
  };

  item.querySelector(".notif-close").onclick = (e) => {
    e.stopPropagation();
    removeNotification(username);
  };

  list.prepend(item);

  notifications[username] = {
    element: item,
    counter: 1,
  };

  updateBadge();
}

function removeNotification(username) {
  if (!notifications[username]) return;
  notifications[username].element.remove();
  delete notifications[username];
  updateBadge();
}

function updateBadge() {
  const badge = document.getElementById("notification-badge");
  const count = Object.keys(notifications).length;
  if (count === 0) {
    badge.style.display = "none";
    return;
  }
  badge.style.display = "flex";
  badge.textContent = count;
}