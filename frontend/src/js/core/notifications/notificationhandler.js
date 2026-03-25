import { setCurrentOpenChat, clearCurrentOpenChat } from "../WebSocket/shownotification";

function openConversation(username) {
    setCurrentOpenChat(username);
}

function backToConversationList() {
    clearCurrentOpenChat();
}

function leaveChatPage() {
    clearCurrentOpenChat();
}