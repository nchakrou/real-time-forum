import { router } from "../Router.js";

const toastQueue = [];
let isToastVisible = false;

export function showNotification(data, showToast = true) {
    storeNotification(data, showToast);
    if (showToast) enqueueToast(data);
}

function enqueueToast(data) {
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

    toast.addEventListener("click", (e) => {
        if (!e.target.classList.contains("toast-close")) {
            toast.remove();
            router(`/chat?username=${data.from}`);
            callback();
        }
    });

    toast.querySelector(".toast-close").addEventListener("click", (e) => {
        e.stopPropagation();
        toast.classList.add("toast-exit");
        setTimeout(() => { toast.remove(); callback(); }, 300);
    });

    setTimeout(() => {
        if (toast.parentElement) {
            toast.classList.add("toast-exit");
            setTimeout(() => { toast.remove(); callback(); }, 300);
        }
    }, 5000);
}

export function storeNotification(data, showToast = true) {
    const bar = document.getElementById("notification-list");
    if (!bar) return;

    let notif = bar.querySelector(`[data-user="${data.from}"]`);

    if (!notif) {
        notif = document.createElement("div");
        notif.classList.add("stored-notification");
        notif.dataset.user = data.from;
        notif.dataset.count = "1";
        notif.dataset.lastMessage = data.message;

        notif.innerHTML = `
            <div class="stored-notif-avatar">${getInitial(data.from)}</div>
            <div class="stored-notif-content">
                <div class="stored-notif-header">
                    <span class="stored-notif-name">${escapeHTML(data.from)}</span>
                    <span class="stored-notif-badge">1</span>
                </div>
                <div class="stored-notif-preview">${truncateUTF8(data.message, 40)}</div>
            </div>
        `;

        notif.addEventListener("click", () => {
            updateBadge(-parseInt(notif.dataset.count));
            notif.remove();
            router(`/chat?username=${data.from}`);
        });

        bar.prepend(notif);
    } else {
        const count = parseInt(notif.dataset.count) + 1;
        notif.dataset.count = count.toString();
        notif.dataset.lastMessage = data.message;
        notif.querySelector(".stored-notif-badge").textContent = count;
        notif.querySelector(".stored-notif-preview").textContent = truncateUTF8(data.message, 40);

        bar.prepend(notif);
    }

    if (showToast) updateBadge(1); // فقط إذا كان إشعار جديد
}

function removeStoredNotification(username) {
    const bar = document.getElementById("notification-list");
    if (!bar) return;
    const notif = bar.querySelector(`[data-user="${username}"]`);
    if (notif) {
        updateBadge(-parseInt(notif.dataset.count));
        notif.remove();
    }
}

function updateBadge(change) {
    const badge = document.getElementById("notification-badge");
    if (!badge) return;
    const current = parseInt(badge.textContent || "0") + change;
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