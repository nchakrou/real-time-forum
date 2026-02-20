import { fetchPosts } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
import { Popup } from "../core/Popup.js";
const homePage = `
${Header}
<div class = "app-home">
<div class = "categories">
  <h2>Categories</h2>
  <ul id = "categories" class="list-categories">
  <li data-category="FPS">FPS</li>
  <li data-category="Battle Royale">Battle Royale</li>
  <li data-category="MOBA">MOBA</li>
  <li data-category="Esports">Esports</li>
  <li data-category="RPG">RPG</li>
  <li data-category="Strategy">Strategy</li>
  <li data-category="Simulation">Simulation</li>
  </ul>
</div>
<div class = "posts">
  <h2>Posts</h2>
  <div id = "posts-container">
  </div>
</div>
<div class = "users">
  <h2>Users</h2>
</div>
</div>
<div id ="popup" class="popup">
<div id ="popup-content">
<p>test</p>
</div>
</div>
`

export function home() {
  document.body.innerHTML = homePage
  Popup.show("test")
  pagesInit()
  const params = new URLSearchParams(window.location.search)
  const category = params.get("category")
  if (!category) {
    fetchPosts(`/api/posts`)
    return
  }
  const categoryElement = document.querySelector(`[data-category="${category}"]`)
  if (categoryElement) {
    categoryElement.classList.add("active")
  } else {
    fetchPosts("/api/posts")
    return
  }
  fetchPosts(`/api/posts?category=${category}`)
}





