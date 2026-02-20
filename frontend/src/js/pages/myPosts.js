import { fetchPosts } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
const myPostsPage = `
${Header}
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
  <div id = "posts-container">
  </div>
</div>
<div class = "users">
  <h2>Users</h2>
</div>
</div>
`

export function myPosts() {
    document.body.innerHTML = myPostsPage
   
    pagesInit("/myPosts")
    fetchPosts("/api/myposts")
}
