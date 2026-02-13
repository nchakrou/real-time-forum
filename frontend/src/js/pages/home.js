import { headerButtons, CategoriesListener } from "../core/Listeners/Listeners.js";
import { fetchPosts } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { ProfileDropdown } from "../components/ProfileDropdown.js";
import { router } from "../core/Router.js";

const homePage = `
${Header}
<div class = "app-home">
<div class = "categories">
  <h2>Categories</h2>
  <ul id = "categories" class="list-categories">
  <li>FPS</li>
  <li>Battle Royale</li>
  <li>MOBA</li>
  <li>Esports</li>
  <li>RPG</li>
  <li>Strategy</li>
  <li>Simulation</li>
  </ul>
</div>
<div class = "posts">
  <h2>Posts</h2>
  <div id = "posts-container">
  </div>
</div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>
`

export function home() {
  document.body.innerHTML = homePage
  ProfileDropdown()
  OnlineUsers()
  headerButtons()
  CategoriesListener()
  fetchPosts("/api/posts")
}
export function OnlineUsers(users) {
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

    const initial = user.charAt(0).toUpperCase()

    userItem.innerHTML = `
      <div class="user-item-avatar">${initial}</div>
      <span class="user-item-name">${user}</span>
      <div class="user-status-dropdown"></div>
    `
    listUsers.appendChild(userItem)
  })

  usersContainer.appendChild(listUsers)
}




