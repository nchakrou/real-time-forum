import { ws } from "../core/WebSocket/initWs.js";
import { router } from "../core/Router.js";
import { headerButtons } from "../core/Listeners/Listeners.js";
import { CategoriesListener } from "../core/Listeners/Listeners.js";
import { ProfileDropdown } from "../components/ProfileDropdown.js";
import { isLogged } from "../main.js";
import { Popup } from "../components/Popup.js";

export function pagesInit(path = "/") {
    ProfileDropdown();
    populateProfile();
    headerButtons();

    if (path === "/chat") {
        ws.send(JSON.stringify({ type: "get_chat_users" }));
    } else {
        ws.send(JSON.stringify({ type: "online_users" }));
        CategoriesListener(path);  
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
        if (Dropdownlogo)
            Dropdownlogo.textContent = user.charAt(0).toUpperCase();
    } catch (e) {
        Popup.show("Something went wrong");
    }
}


export function OnlineUsers(users) {
    if (!ws) return;
    const usersContainer = document.querySelector(".users");
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

    usersContainer.innerHTML = "<h3>Messages</h3>";

    if (!chats || chats.length === 0) {
        const noChats = document.createElement("p");
        noChats.style.marginTop = "20px";
        noChats.textContent = "No conversations yet";
        usersContainer.appendChild(noChats);
        return;
    }

    const listUsers = document.createElement("div");
    listUsers.className = "list-users";

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