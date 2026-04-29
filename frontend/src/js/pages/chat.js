import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
import { ws } from "../core/WebSocket/initWs.js";
import { router } from "../core/Router.js";
import { Popup } from "../components/Popup.js";
import { chatStates } from "../core/chatStates.js";
import { throttle } from "../utils/throttle.js";
import { updateUserList, formatTime } from "../utils/chatUtils.js";
import {
  setCurrentOpenChat,
  clearCurrentOpenChat,
  restoreUnreadDots,
} from "../core/WebSocket/shownotification.js";

const chatPage = `
${Header}
<main class="app-chat">
    <div class="chat-sidebar" id="chat-sidebar">
        <div class="sidebar-header">
            <h3>Conversations</h3>
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
                </div>
            </div>
            <div id="mobile-back-btn" class="mobile-back-btn">
                <button class="chat-action-btn" title="Information">
                    <img src="/frontend/src/assets/plus.svg" style="transform: rotate(45deg);">
                </button>
            </div>
        </div>

        <div class="chat-viewport" id="chat-viewport">
            <div class="empty-chat-state">
                <div class="empty-chat-icon">
                    <img src="/frontend/src/assets/chat.svg" alt="Chat Icon">
                </div>
                <h3>Your Conversations</h3>
                <p>Interact with other users in real-time. Select a contact from the sidebar to start chatting.</p>
            </div>
        </div>
         <div class="typing-indicator" id="typing-indicator" style="display: none;">
            <span class="typing-text"><span class="typing-name"></span> is typing</span>
            <div class="typing-dots">
                <span class="dot"></span>
                <span class="dot"></span>
                <span class="dot"></span>
            </div>
        </div>

        <div class="chat-footer">
            <div class="chat-input-row">

                <div class="input-container">
                    <textarea id="chat-message-input" placeholder="Write something..." rows="1" maxlength="1000"></textarea>
                </div>
                <button id="chat-send-btn" class="chat-send-btn" title="Send message">
                    <img src="/frontend/src/assets/send.svg" alt="Send">
                </button>
            </div>
        </div>
    </div>
</main>
`;

export function chat() {
  document.body.innerHTML = chatPage;
  pagesInit("/chat");

  clearCurrentOpenChat();

  setTimeout(() => restoreUnreadDots(), 100);

  if (window.location.search) {
    const urlParams = new URLSearchParams(window.location.search);
    const username = urlParams.get("username");
    if (username) {
      handleActiveChat(username);
      sentBtn();
    }
  }
}

function handleActiveChat(tagername) {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    Popup.show("Connection lost. Please refresh the page.");
    return;
  }

  setCurrentOpenChat(tagername);

  ws.send(
    JSON.stringify({
      type: "getChat",
      target: tagername,
      lastID: 0,
    }),
  );

  const chatMain = document.getElementById("chat-main");

  if (
    !getComputedStyle(chatMain).display ||
    getComputedStyle(chatMain).display === "none"
  ) {
    chatMain.style.display = "flex";
    const chatSidebar = document.getElementById("chat-sidebar");
    chatSidebar.style.display = "none";
    const mobileBackBtn = document.getElementById("mobile-back-btn");
    mobileBackBtn.style.display = "flex";
    mobileBackBtn.addEventListener("click", () => {
      clearCurrentOpenChat();
      router("/chat");
    });
  }

  const chatAvatar = document.getElementById("active-chat-avatar");
  const chatName = document.getElementById("active-chat-name");
  chatAvatar.textContent = tagername.charAt(0).toUpperCase();
  chatName.textContent = tagername;
  const emptyChatState = document.getElementsByClassName("empty-chat-state")[0];
  emptyChatState.style.display = "none";
  document.getElementById("chat-viewport").addEventListener(
    "scroll",
    throttle(() => {
      if (document.getElementById("chat-viewport").scrollTop <= 50) {
        if (chatStates.isEnd) {
          return;
        }
        ws.send(
          JSON.stringify({
            type: "getChat",
            target: tagername,
            lastID: chatStates.lastID,
          }),
        );
      }
    }, 200),
  );
}
function sentBtn() {
  const messageInput = document.getElementById("chat-message-input");
  const sendBtn = document.getElementById("chat-send-btn");

  if (!messageInput || !sendBtn) return;

  // Prevent duplicate listeners if sentBtn() is called more than once
  if (messageInput.dataset.typingBound === "1") return;
  messageInput.dataset.typingBound = "1";

  let typingTimeout = null;
  let isTyping = false;
  let lastTypingTarget = "";

  const getTarget = () => {
    const el = document.getElementById("active-chat-name");
    return el ? el.textContent.trim() : "";
  };

  const sendTyping = (target) => {
    if (!target || !ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(JSON.stringify({ type: "typing", target }));
  };

  const sendStopTyping = (target) => {
    if (!target || !ws || ws.readyState !== WebSocket.OPEN) return;
    ws.send(JSON.stringify({ type: "stop_typing", target }));
  };

  const stopTypingNow = () => {
    if (typingTimeout) {
      clearTimeout(typingTimeout);
      typingTimeout = null;
    }
    if (isTyping && lastTypingTarget) {
      sendStopTyping(lastTypingTarget);
    }
    isTyping = false;
    lastTypingTarget = "";
  };

  messageInput.addEventListener("input", () => {
    const target = getTarget();
    if (!target) return;

    if (!isTyping) {
      isTyping = true;
      lastTypingTarget = target;
      sendTyping(target);
    } else if (lastTypingTarget && lastTypingTarget !== target) {
      // switched chat while typing
      sendStopTyping(lastTypingTarget);
      lastTypingTarget = target;
      sendTyping(target);
    }

    if (typingTimeout) clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
      if (isTyping && lastTypingTarget) {
        sendStopTyping(lastTypingTarget);
      }
      isTyping = false;
      lastTypingTarget = "";
      typingTimeout = null;
    }, 2000);
  });

  messageInput.addEventListener("blur", stopTypingNow);

  window.addEventListener("beforeunload", stopTypingNow);
  document.addEventListener("visibilitychange", () => {
    if (document.hidden) stopTypingNow();
  });

  sendBtn.addEventListener("click", () => {
    const message = messageInput.value;
    if (!message.trim()) return;

    stopTypingNow();

    const target = getTarget();
    ws.send(
      JSON.stringify({
        type: "message",
        target,
        message,
      }),
    );

    messageInput.value = "";

    const chatViewport = document.getElementById("chat-viewport");
    const messageDiv = document.createElement("div");
    messageDiv.classList.add("Mymessage");

    const sender = document.createElement("h4");
    sender.textContent = "Me";

    const p = document.createElement("p");
    p.textContent = message;

    const timeSpan = document.createElement("span");
    timeSpan.classList.add("message-time");
    timeSpan.textContent = formatTime();

    messageDiv.appendChild(sender);
    messageDiv.appendChild(p);
    messageDiv.appendChild(timeSpan);

    chatViewport.appendChild(messageDiv);
    chatViewport.scrollTop = chatViewport.scrollHeight;

    const user = window.location.search;
    const urlParams = new URLSearchParams(user);
    const username = urlParams.get("username");
    updateUserList(username);
  });
}