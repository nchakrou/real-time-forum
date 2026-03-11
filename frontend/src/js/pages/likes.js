import { createpostsContainer } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { router } from "../core/Router.js";

export async function toggleLike(id, value) {
  try {
    console.log("like :", { id, value });

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
  let main = document.querySelector("main");

  if (!main) {
    main = document.createElement("main");
    document.body.appendChild(main);
  }

  main.innerHTML = `${Header}<div id="posts-container"></div>`;

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
      postsContainer.innerHTML = "<p>No liked posts found.</p>";
      return;
    }

    createpostsContainer(posts);

  } catch (err) {
    console.error("Liked posts error:", err);
    postsContainer.innerHTML = "<p>Failed to load liked posts.</p>";
  }
}