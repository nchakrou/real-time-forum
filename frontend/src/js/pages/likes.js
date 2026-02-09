import { headerButtons } from "../core/Listeners/Listeners.js";
import { createpostsContainer } from "../core/Listeners/postListners.js";

export async function toggleLike(id, value) {
  try {
    console.log("like :", { id, value });
    const res = await fetch(`/api/like?id=${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ value })
    });

    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const data = await res.json();

    const post = document.querySelector(`.post[data-post-id="${id}"]`);
    if (!post) {
      console.error("Post not found for id:", id);
      return;
    }

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

const likedPostsPage = `
<header>
  <h2>Forum</h2>
  <div id="header-buttons">
    <button id="home">Home</button>
    <button id="createpost">Create Post</button>
    <button id="myPosts">My Posts</button>
    <button id="likedPosts">Liked Posts</button>
    <button id="logout">Logout</button>
  </div>
</header>
<div id="posts-container"></div>
`;

export async function likedPosts() {
  document.body.innerHTML = likedPostsPage;
  headerButtons();

  try {
    const response = await fetch("/api/liked-posts", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });

    if (!response.ok) throw new Error(`HTTP ${response.status}`);

    const posts = await response.json();
    createpostsContainer(posts);
  } catch (err) {
    console.error("Failed to fetch liked posts:", err);
    document.getElementById("posts-container").innerHTML = "<p>Failed to load liked posts.</p>";
  }
}