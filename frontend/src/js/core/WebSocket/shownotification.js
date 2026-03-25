import { router } from "../Router.js";

const toastQueue = [];
let isToastVisible = false;
let currentOpenChat = null;
const unreadChats = new Set();

export function setCurrentOpenChat(username) {
    currentOpenChat = username;
    if (username) {
        markChatAsRead(username);
    }
}

export function clearCurrentOpenChat() {
    currentOpenChat = null;
}

function markChatAsRead(username) {
    unreadChats.delete(username);
    removeDotFromSidebar(username);
    removeStoredNotification(username);
}

function markChatAsUnread(username) {
    unreadChats.add(username);
    addDotToSidebar(username);
}

export function showNotification(data, showToast = true) {
    if (currentOpenChat === data.from) {
        return;
    }

    markChatAsUnread(data.from);
    storeNotification(data, showToast);

    if (isOnChatPage()) {
        return;
    }

    if (showToast) enqueueToast(data);
}

export function removeStoredNotification(username) {
    const bar = document.getElementById("notification-list");
    if (!bar) return;

    const notif = findNotifByUser(bar, username);
    if (notif) {
        const badgedCount = parseInt(notif.dataset.badgedCount || "0");
        updateBadge(-badgedCount);
        notif.remove();
    }
}

function addDotToSidebar(username) {
    const userElement = document.querySelector(
        `[data-username="${username}"]`
    );
    if (!userElement) return;
    if (userElement.querySelector(".unread-dot")) return;
    userElement.style.position = "relative";
    const dot = document.createElement("span");
    dot.classList.add("unread-dot");
    userElement.appendChild(dot);
}

function removeDotFromSidebar(username) {
    const userElement = document.querySelector(
        `[data-username="${username}"]`
    );
    if (!userElement) return;
    const dot = userElement.querySelector(".unread-dot");
    if (dot) dot.remove();
}

export function restoreUnreadDots() {
    unreadChats.forEach((username) => {
        addDotToSidebar(username);
    });
}

function isOnChatPage() {
    return window.location.pathname.startsWith("/chat");
}

function enqueueToast(data) {
    const existing = toastQueue.find((item) => item.from === data.from);
    if (existing) {
        existing.message = data.message;
        return;
    }
    toastQueue.push(data);
    if (!isToastVisible) showNextToast();
}

function showNextToast() {
    if (toastQueue.length === 0) {
        isToastVisible = false;
        return;
    }
    isToastVisible = true;
    const data = toastQueue.shift();
    displayToast(data, () => showNextToast());
}

function displayToast(data, callback) {
    let container = document.getElementById("toast-container");
    if (!container) {
        container = document.createElement("div");
        container.id = "toast-container";
        document.body.appendChild(container);
    }
    const toast = document.createElement("div");
    toast.classList.add("toast-notification");
    toast.innerHTML = `
        <div class="toast-avatar">${getInitial(data.from)}</div>
        <div class="toast-content">
            <div class="toast-sender">${escapeHTML(data.from)}</div>
            <div class="toast-message">${truncateUTF8(data.message, 50)}</div>
        </div>
        <div class="toast-close">✕</div>
    `;
    container.appendChild(toast);
    requestAnimationFrame(() => toast.classList.add("toast-enter"));
    let dismissed = false;
    let autoTimer = null;
    const dismiss = (navigate = false) => {
        if (dismissed) return;
        dismissed = true;
        if (autoTimer) clearTimeout(autoTimer);
        if (navigate) {
            toast.remove();
            router(`/chat?username=${data.from}`);
            callback();
            return;
        }
        toast.classList.add("toast-exit");
        setTimeout(() => {
            toast.remove();
            callback();
        }, 300);
    };
    toast.addEventListener("click", (e) => {
        if (!e.target.classList.contains("toast-close")) {
            dismiss(true);
        }
    });
    toast.querySelector(".toast-close").addEventListener("click", (e) => {
        e.stopPropagation();
        dismiss(false);
    });
    autoTimer = setTimeout(() => dismiss(false), 5000);
}

export function storeNotification(data, showToast = true) {
    const bar = document.getElementById("notification-list");
    if (!bar) return;

    let notif = findNotifByUser(bar, data.from);

    if (!notif) {
        notif = document.createElement("div");
        notif.classList.add("stored-notification");
        notif.dataset.user = data.from;
        notif.dataset.count = "0";
        notif.dataset.badgedCount = "0";
        notif.dataset.lastMessage = data.message;

        notif.innerHTML = `
            <div class="stored-notif-avatar">${getInitial(data.from)}</div>
            <div class="stored-notif-content">
                <div class="stored-notif-header">
                    <span class="stored-notif-name">${escapeHTML(data.from)}</span>
                    <span class="stored-notif-badge">0</span>
                </div>
                <div class="stored-notif-preview">${truncateUTF8(data.message, 40)}</div>
            </div>
        `;

        notif.addEventListener("click", () => {
            const badgedCount = parseInt(notif.dataset.badgedCount || "0");
            updateBadge(-badgedCount);
            notif.remove();
            router(`/chat?username=${data.from}`);
        });

        bar.prepend(notif);
    }

    const count = parseInt(notif.dataset.count) + 1;
    notif.dataset.count = count.toString();
    notif.dataset.lastMessage = data.message;
    notif.querySelector(".stored-notif-badge").textContent = count;
    notif.querySelector(".stored-notif-preview").textContent = truncateUTF8(
        data.message,
        40
    );
    bar.prepend(notif);

    if (showToast) {
        const badgedCount = parseInt(notif.dataset.badgedCount || "0") + 1;
        notif.dataset.badgedCount = badgedCount.toString();
        updateBadge(1);
    }
}

function findNotifByUser(container, username) {
    return (
        [...container.children].find((el) => el.dataset.user === username) ||
        null
    );
}

function updateBadge(change) {
    const badge = document.getElementById("notification-badge");
    if (!badge) return;
    const current = Math.max(0, parseInt(badge.textContent || "0") + change);
    if (current <= 0) {
        badge.style.display = "none";
        badge.textContent = "0";
    } else {
        badge.style.display = "flex";
        badge.textContent = current.toString();
    }
}

function truncateUTF8(str, max) {
    if (!str) return "";
    const chars = Array.from(str);
    return chars.length > max ? chars.slice(0, max).join("") + "..." : str;
}

function getInitial(name) {
    return Array.from(name)[0]?.toUpperCase() || "?";
}

function escapeHTML(str) {
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
}