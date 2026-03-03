export function showNotification(data) {
    const notificationBar = document.getElementById("notification-list");
    if (!notificationBar) return;

    const notif = document.createElement("div");
    notif.classList.add("notification");
    notif.textContent = `${data.from} ${data.content}`;

    if (data.post_id) {
        notif.addEventListener("click", () => {
            window.location.href = `/post?id=${data.post_id}`;
        });
    }
    notificationBar.prepend(notif);

    const badge = document.getElementById("notification-badge");
    badge.style.display = "flex";
    badge.textContent = parseInt(badge.textContent || "0") + 1;
}