export function message(data) {
    if (window.location.search.includes(data.from)) {
        const chatViewport = document.getElementById("chat-viewport")
        const message = document.createElement("div")
        message.classList.add("message")
        message.textContent = data.message
        chatViewport.appendChild(message)
    } else {
        MessageNotification(data.from)
    }
}
function MessageNotification(username) {
    const notificationBar = document.getElementById("notification-list")
    let notification = notificationBar.querySelector(`[data-user="${username}"]`)
    const badge = document.getElementById("notification-badge")
    if (window.getComputedStyle(badge).display === "none") {
        badge.style.display = "flex"
    }
    badge.textContent = parseInt(badge.textContent) + 1
    if (notification) {
        const counter = ++notification.dataset.counter
        notification.textContent = `You have ${counter} message from ${username}`
    } else {
        notification = document.createElement("div")
        notification.classList.add("notification")
        notification.dataset.user = username
        notification.dataset.counter = 1
        notification.textContent = `You have 1 message from ${username}`
        notificationBar.appendChild(notification)
    }


}
export function chatHistory(data) {
    const chatViewport = document.getElementById("chat-viewport")
    if (data.Messages && data.Messages.length > 0) {
        
        data.Messages.forEach(message => {
            const messageDiv = document.createElement("div")
            const target = new URLSearchParams(window.location.search).get("username")
            console.log("hadi",message.from, target);
            
            if (message.from === target) {
                messageDiv.classList.add("message")
            } else {
                messageDiv.classList.add("Mymessage")
            }
            messageDiv.textContent = message.message
            chatViewport.appendChild(messageDiv)
        })
    }else{
        chatViewport.innerHTML = `
        <div class="empty-chat-state">
            <div class="empty-chat-icon">
                <img src="src/assets/chat.svg" alt="Chat Icon">
            </div>
            <h3>No messages yet</h3>
            <p>Start a conversation</p>
        </div>
        `
    }
}