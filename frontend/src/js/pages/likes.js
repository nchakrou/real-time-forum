import { createpostsContainer } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { router } from "../core/Router.js";
import { pagesInit } from "../components/pagesInit.js";

export async function toggleLike(id, value) {
  try {
    const res = await fetch(`/api/like?id=${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ value }),
      credentials: "include",
    });

    if (res.status === 401) {
      router("/login");
      return;
    }

    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const data = await res.json();

    const post = document.querySelector(`.post[data-post-id="${id}"]`);
    if (!post) return console.error("Post not found for id:", id);

    const likeBtn = post.querySelector(".like_button");
    const dislikeBtn = post.querySelector(".dislike_button");

    likeBtn.querySelector("span").textContent = data.likes;
    dislikeBtn.querySelector("span").textContent = data.dislikes;

    likeBtn.classList.toggle("active", data.userValue === 1);
    dislikeBtn.classList.toggle("active", data.userValue === -1);
  } catch (err) {
    console.error("Like toggle failed:", err);
  }
}

export async function likedPosts() {
  const likedPage = `
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
        <div class="app-home">
            <div class="categories">
                <h2>Categories</h2>
                <ul id="categories" class="list-categories">
                    <li data-category="FPS">FPS</li>
                    <li data-category="Battle Royale">Battle Royale</li>
                    <li data-category="MOBA">MOBA</li>
                    <li data-category="Esports">Esports</li>
                    <li data-category="RPG">RPG</li>
                    <li data-category="Strategy">Strategy</li>
                    <li data-category="Simulation">Simulation</li>
                </ul>
            </div>
            <div class="posts">
                <h2>Liked Posts</h2>
                <div id="posts-container"></div>
            </div>
            <div class="online-users">
                <h2>Users</h2>
            </div>
        </div>
    `;
  document.body.innerHTML = likedPage;
  pagesInit("/liked");

  const postsContainer = document.getElementById("posts-container");
  try {
    const response = await fetch("/api/liked-posts", {
      method: "GET",
      credentials: "include",
    });

    if (response.status === 401) {
      router("/login");
      return;
    }
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const posts = await response.json();
    if (!Array.isArray(posts) || posts.length === 0) {
      postsContainer.innerHTML = `
                <div class="empty-chat-state">
                    <h3>No liked posts yet</h3>
                    <p>Posts you like will appear here.</p>
                </div>
            `;
      return;
    }
    postsContainer.innerHTML = "";
    createpostsContainer(posts);
  } catch (err) {
    console.error("Liked posts error:", err);
    postsContainer.innerHTML = "<p>Failed to load liked posts.</p>";
  }
}
