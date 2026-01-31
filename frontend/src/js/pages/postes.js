
import { headerButtons } from "../core/Listeners.js"
const createPostPage = `<header>
<h2>Forum</h2>
<div id = "header-buttons">
<button id = "home">Home</button>
<button id ="createpost">Create Post</button>
<button id ="myPosts">My Posts</button>
<button id ="likedPosts">Liked Posts</button>
<button id = "logout">logout</button>
</div>
</header>
<div class = "app-home">
<div id = "categories" class = "categories">
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
  <div class = "postTitle">
  <p>Post Title</p>
  <input type = "text" id = "post" >
  </div>
  <div class = "postContent">
  <p>Post Content</p>
  <textarea id = "content"></textarea>
  </div>
  <div class = "postCategory">
  <p>Post Category</p>
  <div class = "CreatePostCategories" id = "CreatePostCategories">
    <div>Cybersecurity </div>
    <div>Esports</div>
    <div>MOBA</div>
    <div>RPG</div>
    <div>Strategy</div>
    <div>Simulation</div>
    <div>FPS</div>
    <div>Battle Royale</div>
  </div>
  </div>
  <button id = "createpost">Create Post</button>
    
  </div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>`
export function createPost() {
  document.body.innerHTML = createPostPage
  headerButtons()
  createPostCategoriesListener()
}
function createPostCategoriesListener() {
  document.getElementById("CreatePostCategories").addEventListener("click", (e) => {
    if (e.target.tagName === "DIV") {
      e.target.classList.toggle("active")
    }

  })
}