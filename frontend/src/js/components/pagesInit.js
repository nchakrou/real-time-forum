import { ws } from "../core/WebSocket/initWs.js"
import { router } from "../core/Router.js"
import { headerButtons } from "../core/Listeners/Listeners.js"
import { CategoriesListener } from "../core/Listeners/Listeners.js"
import { ProfileDropdown } from "../components/ProfileDropdown.js"
import { isLogged } from "../main.js"

export function pagesInit(path = "/") {
    ProfileDropdown()
    populateProfile()
    ws.send(JSON.stringify({ type: "online_users" }))
    headerButtons()
    if (path !== "/chat") {
        CategoriesListener(path)
    }
}
async function populateProfile() {
    try {
        const [user,log] = await isLogged()
        if (!log || !user) {
            throw new Error("User not logged in")
        }
        const logo = document.getElementById('user-initials')
        const Dropdownlogo = document.getElementById('user-initials-dropdown')
        const username = document.getElementById('user-name-dropdown')
        if (username) username.textContent = user 
        if (logo) logo.textContent = user.charAt(0).toUpperCase()
        if (Dropdownlogo) Dropdownlogo.textContent = user.charAt(0).toUpperCase()
    } catch (e) {
        console.log(e)
    }
}
export function OnlineUsers(users) {

    if (!ws) return
    const usersContainer = document.querySelector(".users")
    usersContainer.innerHTML = "<h2>Users</h2>"

    if (!users || users.length === 0) {
        const noUsers = document.createElement("p")
        noUsers.style.marginTop = "20px"
        noUsers.textContent = "No online users"
        usersContainer.appendChild(noUsers)
        return
    }

    const listUsers = document.createElement("div")
    listUsers.className = "list-users"

    users.forEach(user => {
        const userItem = document.createElement("div")
        userItem.className = "user-item"
        userItem.dataset.username = user
        userItem.addEventListener("click", () => {
            router(`/chat?username=${user}`)
        })

        const avatar = user.charAt(0).toUpperCase()

        userItem.innerHTML = `
      <div class="user-item-avatar">${avatar}</div>
      <span class="user-item-name">${user}</span>
      <div class="user-status-dropdown"></div>
    `
        listUsers.appendChild(userItem)
    })

    usersContainer.appendChild(listUsers)
}
