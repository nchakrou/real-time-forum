import { fetchPosts, states } from "../core/Listeners/postListners.js";
import { throttle } from "../utils/throttle.js";
import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
import { Popup } from "../components/Popup.js";

const homePage = `
${Header}
<div class="mobile-toggles">
  <button id="toggle-categories" class="mobile-toggle-btn">
    <img src="/frontend/src/assets/plus.svg" alt="Categories">
    <span>Categories</span>
  </button>
  <button id="toggle-users" class="mobile-toggle-btn">
    <img src="/frontend/src/assets/plus.svg" alt="Users">
    <span>Users</span>
  </button>
</div>
<div class = "app-home">
<div class = "categories" id="mobile-categories">
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
<div class = "online-users" id="mobile-users">
  <h2>Users</h2>
</div>
</div>
<div id ="popup" class="popup">
<div id ="popup-content">
<p>test</p>
</div>
</div>
`;

export function home() {
  document.body.innerHTML = homePage;
  pagesInit();
  const params = new URLSearchParams(window.location.search);
  const category = params.get("category");
  const categoryElement = document.querySelector(
    `[data-category="${category}"]`,
  );
  if (categoryElement) {
    categoryElement.classList.add("active");
    fetchPosts(`/api/posts?category=${category}`);
  } else {
    fetchPosts("/api/posts");
  }

  // Infinite Scroll Listener
  const postsContainer = document.querySelector(".posts");
  if (postsContainer) {
    postsContainer.addEventListener(
      "scroll",
      throttle(() => {
        const { scrollTop, scrollHeight, clientHeight } = postsContainer;
        console.log(scrollTop + clientHeight - scrollHeight - 200);

        if (scrollTop + clientHeight >= scrollHeight - 200) {
          const path = category
            ? `/api/posts?category=${category}`
            : "/api/posts";
          fetchPosts(path);
        }
      }, 200),
    );
  }
}
