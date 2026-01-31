import { headerButtons } from "../core/Listeners.js";

const myPostsPage = `
<header>
<h2>Forum</h2>
<div id  = "header-buttons">
<button id = "home">Home</button>
<button id ="createpost">Create Post</button>
<button id ="myPosts">My Posts</button>
<button id ="likedPosts">Liked Posts</button>
<button id = "logout">logout</button>
</div>
</header>
<div class = "app-home">
<div id = "categories"class = "categories">
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
  <h2>My Posts</h2>
  <p>No posts yet</p>
</div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>
`

export function myPosts() {
    document.body.innerHTML = myPostsPage
    headerButtons()
}