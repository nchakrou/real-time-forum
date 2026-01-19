import { router } from "../Router.js"
const homePage = `
<header>
<h2>Forum</h2>
<div>
<button id ="creatpost">Create Post</button>
  <button>created posts</button>
  <button>liked posts</button>
  <button id = "logout">logout</button>
</div>
</header>
<div class = "app-home">
<div class = "categories">
  <h2>Categories</h2>
  <ul class="list-categories">
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
  <p>No posts yet</p>
</div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>
`
const createPostPage = `<header>
<h2>Forum</h2>
<div>
<button id ="createpost">Create Post</button>
  <button>created posts</button>
  <button>liked posts</button>
  <button id = "logout">logout</button>
</div>
</header>
<div class = "app-home">
<div class = "categories">
  <h2>Categories</h2>
  <ul class="list-categories">
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
  <h2>Create Post</h2>
  <input type = "text" id = "post" >
  <input type="checkbox" name="categories" value="Cybersecurity">Cybersecurity</label>
</div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>`
export async function home() {
    document.body.innerHTML = homePage
    logoutListener()
    creatPostListener()
}
async function logoutListener() {
    document.getElementById("logout").addEventListener("click", async () => {
        const req = await fetch("/api/logout", {
            method: "GET"
        })
        router("/")
    })
}
function creatPostListener() {
    document.getElementById("creatpost").addEventListener("click", () => {
        router("/createpost")
    })

}
export function createPost() {
    document.body.innerHTML = createPostPage
}