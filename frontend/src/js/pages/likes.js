import { createpostsContainer } from "../core/Listeners/postListners.js";
import { Header } from "../components/Header.js";
import { pagesInit } from "../components/pagesInit.js";
export async function toggleLike(id, value) {
  try {
    console.log("like :", { id, value });
    const res = await fetch(`/api/like?id=${id}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ value }),
      credentials: "include"
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
${Header}
`;

export async function likedPosts() {
  document.body.innerHTML = likedPostsPage;
  pagesInit("/api/liked-posts")
  try {
    const response = await fetch("/api/liked-posts", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
      credentials: "include"
    });

    if (!response.ok) throw new Error(`HTTP ${response.status}`);

    const posts = await response.json();
    createpostsContainer(posts);
  } catch (err) {
    console.log(err);

    document.getElementById("posts-container").innerHTML = "<p>Failed to load liked posts.</p>";
  }
}