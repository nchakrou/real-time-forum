
import { router } from "../core/Router.js"
import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
const createPostPage = `${Header}
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
    <div data-id="1" >FPS</div>
    <div data-id="2">Battle Royale</div>
    <div data-id="3">MOBA</div>
    <div data-id="4">Esports</div>
    <div data-id="5">RPG</div>
    <div data-id="6">Strategy</div>
    <div data-id="7">Simulation</div>
  </div>
  </div>
  <button id = "submitpost">Create Post</button>
    
  </div>
<div class = "users">
  <h2>Users</h2>
</div>
</div>`
export function createPost() {
  document.body.innerHTML = createPostPage
    pagesInit()
  handleCreatePost()
  createPostCategoriesListener()
}
function createPostCategoriesListener() {
  document.getElementById("CreatePostCategories").addEventListener("click", (e) => {
    if (e.target.dataset.id) {
      e.target.classList.toggle("active")
    }

  })
}
function handleCreatePost() {
  document.getElementById("submitpost").addEventListener("click", async () => {
    
    const title = document.getElementById("post").value
    const content = document.getElementById("content").value
    const CategoriesElements = document.querySelectorAll(
      "#CreatePostCategories .active"
    )
    const selectedCategories = Array.from(CategoriesElements).map(
      (el) => Number(el.dataset.id)
    )

    if (!title || !content || selectedCategories.length === 0) {
      alert("Please fill in all fields and select at least one category.");
      return;
    }

    const postData = {
      title: title,
      content: content,
      categories: selectedCategories,
    }
    
      const response = await fetch("/api/createpost", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(postData),
      })
  console.log(response.statusText);

      if (response.ok) {
        
        router("/")
      } else {
        const errorData = await response.json()
        alert(`Error: ${errorData.message}`)
      }
    
  })
}