import { router } from "../core/Router.js";

export function updateUserList(username) {
  const usersDiv = document.getElementsByClassName("users")[0];
  if (!usersDiv) {
    return;
  }

  const target = document.querySelector(`[data-username="${username}"]`);

  if (target) {
    usersDiv.prepend(target);
  } else {
    const noChats = document.getElementById("no-chats");
    if (noChats) {
      noChats.remove();
    }
    const userItem = document.createElement("div");
    userItem.className = "user-item";
    userItem.dataset.username = username;
    userItem.addEventListener("click", () => {
      router(`/chat?username=${username}`);
    });

    const avatar = username.charAt(0).toUpperCase();
    userItem.innerHTML = `
      <div class="user-item-avatar">${avatar}</div>
      <span class="user-item-name">${username}</span>
    `;
    usersDiv.prepend(userItem);
  }
}

export function formatTime(timestamp) {
  const date = timestamp ? new Date(timestamp) : new Date();

  return date.toLocaleString([], {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}
