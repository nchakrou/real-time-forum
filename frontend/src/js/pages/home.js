import { headerButtons ,CategoriesListener} from "../core/Listeners/Listeners.js";
import { fetchPosts } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { ProfileDropdown } from "../components/ProfileDropdown.js";

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
  headerButtons()
  CategoriesListener()
  fetchPosts("/api/posts")
}




