import { router } from "../core/Router.js"
import { headerButtons, CategoriesListener,PostButtonsListener } from "../core/Listeners.js"
const homePage = `
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
  headerButtons()
  CategoriesListener()
  getPosts()
}
async function getPosts() {
 
    const response = await fetch("/api/posts", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })

    if (response.ok) {
      const posts = await response.json()
      console.log(posts)
      const postsContainer = document.getElementById("posts-container")
    if (!posts ) {
     
      const noPosts = document.createElement("p")
      noPosts.textContent = "No posts yet"
      postsContainer.appendChild(noPosts)
    } else {
      
      posts.forEach((post) => {
        const postElement = document.createElement("div")
        postElement.classList.add("post")
        postElement.innerHTML = `
          <h3>${post.title}</h3>
          <p>${post.content}</p>
          <p><strong>Category:</strong> ${post.category}</p>
          <div class= "post-buttons">
          <button type="submit" class = "like_button">👍 <span>${post.likes}</span> Like</button>
          <button type="submit" class = "dislike_button">👎 <span>${post.dislikes}</span> Dislike</button>
          <button type="submit" class = "comment_button">💬<span> ${post.comments}</span> comment</button>
          </div>
          <div class = "comments-section hidden">
          <input type="text" placeholder="Write a comment...">
          <button type="submit">Submit</button>
          </div>
        `
        postsContainer.appendChild(postElement)
        
        
      })
PostButtonsListener()
    }
    } else {
      alert("Failed to fetch posts")
      return
    }
    
  }



