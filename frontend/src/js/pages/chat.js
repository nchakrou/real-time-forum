import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
import { ws } from "../core/WebSocket/initWs.js"
import { router } from "../core/Router.js"
const chatPage = `
${Header}
<main class="app-chat">
    <div class="chat-sidebar" id="chat-sidebar">
        <div class="sidebar-header">
            <h3>Messages</h3>
        </div>
        <div class="users chat-users-list">
        </div>
    </div>

    <div class="chat-main" id="chat-main">
        <div class="chat-header">
            <div class="current-contact-info">
                <div class="contact-avatar-lg" id="active-chat-avatar">?</div>
                <div class="contact-meta">
                    <h4 id="active-chat-name">Select a conversation</h4>
                    <span class="online-status" id="active-chat-status">Offline</span>
                </div>
            </div>
            <div id="mobile-back-btn" class="mobile-back-btn">
                <button class="chat-action-btn" title="Information">
                    <img src="src/assets/plus.svg" style="transform: rotate(45deg);">
                </button>
            </div>
        </div>

        <div class="chat-viewport" id="chat-viewport">
            <div class="empty-chat-state">
                <div class="empty-chat-icon">
                    <img src="src/assets/chat.svg" alt="Chat Icon">
                </div>
                <h3>Your Conversations</h3>
                <p>Interact with other users in real-time. Select a contact from the sidebar to start chatting.</p>
            </div>
        </div>

        <div class="chat-footer">
            <div class="chat-input-row">

                <div class="input-container">
                    <textarea id="chat-message-input" placeholder="Write something..." rows="1"></textarea>
                </div>
                <button id="chat-send-btn" class="chat-send-btn" title="Send message">
                    <img src="src/assets/send.svg" alt="Send">
                </button>
            </div>
        </div>
    </div>
</main>
`

export function chat() {
    document.body.innerHTML = chatPage
    pagesInit("/chat")
    console.log(window.location.pathname);
    if (window.location.search) {
        const urlParams = new URLSearchParams(window.location.search);
        const username = urlParams.get('username');
        if (username) {
            handleActiveChat(username)
        }

    }
}

function handleActiveChat(username) {
    if (!ws) {
        alert("WebSocket is not open")
        return
    }
    ws.send(JSON.stringify({
        type: "getChat",
        target: username
    }))
    //mobile handle
    const chatMain = document.getElementById("chat-main")

    console.log(getComputedStyle(chatMain).display);
    if (!getComputedStyle(chatMain).display || getComputedStyle(chatMain).display === "none") {
        chatMain.style.display = "flex"
        const chatSidebar = document.getElementById("chat-sidebar")
        chatSidebar.style.display = "none"
        const mobileBackBtn = document.getElementById("mobile-back-btn")
        mobileBackBtn.style.display = "flex"
        mobileBackBtn.addEventListener("click", () => {
            router("/chat")
        })
    }

}