import { ws } from "../core/WebSocket/initWs.js";
import { router } from "../core/Router.js";
import { headerButtons } from "../core/Listeners/Listeners.js";
import { CategoriesListener } from "../core/Listeners/Listeners.js";
import { ProfileDropdown } from "../components/ProfileDropdown.js";
import { isLogged } from "../main.js";
import { Popup } from "../components/Popup.js";
import {
  renderStoredNotifications,
  restoreUnreadDots,
} from "../core/WebSocket/shownotification.js";

export function pagesInit(path = "/") {
  ProfileDropdown();
  populateProfile();
  headerButtons();
  renderStoredNotifications();
  restoreUnreadDots();

  if (path === "/chat") {
    ws.send(JSON.stringify({ type: "get_chat_users" }));
  } else {
    CategoriesListener("/");
  }
  ws.send(JSON.stringify({ type: "online_users" }));

  initMobileToggles();
}

function initMobileToggles() {
  const toggleCategoriesBtn = document.getElementById("toggle-categories");
  const toggleUsersBtn = document.getElementById("toggle-users");
  const mobileCategories = document.querySelector(".categories");
  const mobileUsers = document.querySelector(".online-users");

  if (toggleCategoriesBtn && mobileCategories) {
    // Remove old listeners to avoid duplicates if re-rendering
    const newCatBtn = toggleCategoriesBtn.cloneNode(true);
    toggleCategoriesBtn.parentNode.replaceChild(newCatBtn, toggleCategoriesBtn);

    newCatBtn.addEventListener("click", () => {
      mobileCategories.classList.toggle("show-mobile");
      if (mobileUsers) mobileUsers.classList.remove("show-mobile");
    });
  }

  if (toggleUsersBtn && mobileUsers) {
    const newUserBtn = toggleUsersBtn.cloneNode(true);
    toggleUsersBtn.parentNode.replaceChild(newUserBtn, toggleUsersBtn);

    newUserBtn.addEventListener("click", () => {
      mobileUsers.classList.toggle("show-mobile");
      if (mobileCategories) mobileCategories.classList.remove("show-mobile");
    });
  }
}

async function populateProfile() {
  try {
    const [user, log] = await isLogged();

    if (!log || !user) {
      throw new Error("User not logged in");
    }

    const logo = document.getElementById("user-initials");
    const Dropdownlogo = document.getElementById("user-initials-dropdown");
    const username = document.getElementById("user-name-dropdown");

    if (username) username.textContent = user;
    if (logo) logo.textContent = user.charAt(0).toUpperCase();
    if (Dropdownlogo) Dropdownlogo.textContent = user.charAt(0).toUpperCase();
  } catch (e) {
    Popup.show("Something went wrong");
  }
}

export function OnlineUsers(users) {
  if (!ws) return;
  if (window.location.pathname.includes("/chat")) {
    if (!users) return;
    users.forEach((user) => {
      const div = document.querySelector(`.user-item[data-username="${user}"]`);
      if (div) {
        div.classList.add("online");
      }
    });
  }
  const usersContainer = document.querySelector(".online-users");
  if (!usersContainer) return;

  usersContainer.innerHTML = "<h2>Users</h2>";

  if (!users || users.length === 0) {
    const noUsers = document.createElement("p");
    noUsers.style.marginTop = "20px";
    noUsers.textContent = "No online users";
    usersContainer.appendChild(noUsers);
    return;
  }

  const listUsers = document.createElement("div");
  listUsers.className = "list-users";

  users.forEach((user) => {
    const userItem = document.createElement("div");
    userItem.className = "user-item";
    userItem.dataset.username = user;
    userItem.addEventListener("click", () => {
      router(`/chat?username=${user}`);
    });

    const avatar = user.charAt(0).toUpperCase();
    userItem.innerHTML = `
            <div class="user-item-avatar">${avatar}</div>
            <span class="user-item-name">${user}</span>
        `;
    listUsers.appendChild(userItem);
  });

  usersContainer.appendChild(listUsers);
}

export function ChatUsers(chats) {
  const usersContainer = document.querySelector(".users");
  if (!usersContainer) return;
  if (document.getElementById("no-chats")) {
    document.getElementById("no-chats").remove();
  }
  let users = false;
  if (document.querySelector(".list-users")) {
    users = true;
  }
  if (!chats || chats.length === 0) {
    const noChats = document.createElement("p");
    noChats.style.marginTop = "20px";
    noChats.textContent = "No conversations yet";
    noChats.id = "no-chats";
    usersContainer.appendChild(noChats);
    return;
  }

  if (!chats || chats.length === 0) {
    const noChats = document.createElement("p");
    noChats.style.marginTop = "20px";
    noChats.textContent = "No conversations yet";
    usersContainer.appendChild(noChats);
    return;
  }
  let listUsers;
  if (!users) {
    listUsers = document.createElement("div");
    listUsers.className = "list-users";
  } else {
    listUsers = document.querySelector(".list-users");
  }

  chats.forEach(({ target }) => {
    const userItem = document.createElement("div");
    userItem.className = "user-item";
    userItem.dataset.username = target;
    userItem.addEventListener("click", () => {
      router(`/chat?username=${target}`);
    });

    const avatar = target.charAt(0).toUpperCase();
    userItem.innerHTML = `
            <div class="user-item-avatar">${avatar}</div>
            <div class="user-item-info">
                <span class="user-item-name">${target}</span>
            </div>
        `;
    listUsers.appendChild(userItem);
  });

  usersContainer.appendChild(listUsers);
}
