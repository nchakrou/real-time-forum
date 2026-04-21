import { toggleLike } from "../../pages/likes.js";
import { Popup } from "../../components/Popup.js";
import { router } from "../Router.js";

export const postButtons = {
  like_button: (e) => likeListener(e),
  dislike_button: (e) => dislikeListener(e),
  comment_button: (e) => commentListener(e),
  Submit_comment: (e) => submitCommentListener(e),
  "load-more-comments": (e) => loadMoreComments(e),
};

export const states = {
  path: "",
};

export async function fetchPosts(path) {
  console.log("ok");
  
  try {
    const response = await fetch(path, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (response.status === 401) {
      router("/login");
      return;
    }

    if (response.ok) {
      const posts = await response.json();
      createpostsContainer(posts.Posts, true);
    } else {
      throw new Error(response);
    }
  } catch (error) {
    Popup.show("Connection error");
  }
}

export function createpostsContainer(posts, isFirstLoad = false) {
  const postsContainer = document.getElementById("posts-container");

  if (isFirstLoad) {
    postsContainer.innerHTML = "";
  }

  if (!posts || posts.length === 0) {
    if (isFirstLoad) {
      const noPosts = document.createElement("p");
      noPosts.textContent = "No posts yet";
      postsContainer.appendChild(noPosts);
    }
  } else {
    posts.forEach((post) => {
      const postElement = document.createElement("div");
      postElement.classList.add("post");
      postElement.dataset.postId = post.id;

      const h3 = document.createElement("h3");
      const p = document.createElement("p");

      h3.textContent = post.title;
      p.textContent = post.content;

      postElement.appendChild(h3);
      postElement.appendChild(p);

      postElement.insertAdjacentHTML(
        "beforeend",
        `
        <div class="post-categories">
          ${(post.categories || []).map((cat) => `<span class="category-tag">${cat}</span>`).join("")}
        </div>
      
        <div class="post-buttons">
          <button type="submit" class="like_button">
            <img src="/frontend/src/assets/like.svg">
            <span>${post.likes}</span>
          </button>
      
          <button type="submit" class="dislike_button">
            <img src="/frontend/src/assets/dislike.svg">
            <span>${post.dislikes}</span>
          </button>
      
          <button type="submit" class="comment_button">
            <img src="/frontend/src/assets/comment.svg">
            <span>${post.comments}</span>
          </button>
        </div>
      
       <div class="comments-container hidden">

  <div class="comments">
  </div>

  <div class="comments-section">
    <input type="text" class="comment-input" placeholder="Add a comment..." maxlength="200">
    <button type="submit" class="Submit_comment">
      <img src="/frontend/src/assets/send.svg">
    </button>
  </div>

</div>
        `,
      );

      postsContainer.appendChild(postElement);
      const likeBtn = postElement.querySelector(".like_button");
      const dislikeBtn = postElement.querySelector(".dislike_button");

      likeBtn.classList.toggle("active", post.userValue === 1);
      dislikeBtn.classList.toggle("active", post.userValue === -1);

      loadComments(postElement);
    });

    PostButtonsListener();
  }
}

export function PostButtonsListener() {
  const postsContainer = document.getElementById("posts-container");

  const clone = postsContainer.cloneNode(true);
  postsContainer.parentNode.replaceChild(clone, postsContainer);

  clone.addEventListener("click", (e) => {
    const button = e.target.closest("button");
    if (!button) return;

    const buttonClass = button.classList[0];
    if (postButtons[buttonClass]) postButtons[buttonClass](e);
  });
}

async function submitCommentListener(e) {
  const post = e.target.closest(".post");
  const commentInput = post.querySelector(".comment-input");

  const commentText = commentInput.value.trim();
  if (!commentText) return;

  const postId = post.dataset.postId;

  try {
    const formData = new FormData();
    formData.append("post_id", postId);
    formData.append("comment", commentText);

    const res = await fetch("/api/add-comment", {
      method: "POST",
      credentials: "include",
      body: formData,
    });

    if (res.status === 401) {
      router("/login");
      return;
    }

    if (!res.ok) {
      const text = await res.text();
      console.error("Server error:", text);
      return;
    }

    const data = await res.json();

    if (data.status === "ok") {
      commentInput.value = "";

      const commentBtn = post.querySelector(".comment_button span");

      commentBtn.textContent = data.comments;

      const commentsContainer = post.querySelector(".comments");
      commentsContainer.innerHTML = "";

      await loadComments(post);
      commentsContainer.scrollTop = commentsContainer.scrollHeight;
    }
  } catch (err) {
    console.error("Error adding comment:", err);
  }
}

async function likeListener(e) {
  const post = e.target.closest(".post");
  const postId = post.dataset.postId;

  await toggleLike(postId, 1);
}

async function dislikeListener(e) {
  const post = e.target.closest(".post");
  const postId = post.dataset.postId;

  await toggleLike(postId, -1);
}

async function commentListener(e) {
  const button = e.target.closest(".comment_button");
  if (!button) return;

  const post = button.closest(".post");
  if (!post) return;

  const section = post.querySelector(".comments-container");
  if (!section) return;

  const isHidden = section.classList.toggle("hidden");

  if (!isHidden) {
    post.querySelector(".comments").innerHTML = "";
    await loadComments(post);
  }
}

async function loadComments(post) {
  const commentsContainer = post.querySelector(".comments");
  const postId = post.dataset.postId;

  try {
    const res = await fetch(`/api/comments?post_id=${postId}`);
    if (!res.ok) throw new Error("Failed to fetch comments");
    const comments = await res.json();
    if (!comments || comments.length === 0) {
      commentsContainer.innerHTML = "<p>No comments yet</p>";
      return;
    }
    comments.forEach((c) => {
      const newComment = document.createElement("div");
      newComment.classList.add("comment");
      newComment.innerHTML = `
        <span class="comment-user">${c.username}</span>
        <span class="comment-text">${c.content}</span>
      `;
      commentsContainer.appendChild(newComment);
    });
  } catch (err) {
    console.error("Error loading comments:", err);
  }
}
async function loadMoreComments(e) {
  const post = e.target.closest(".post");
  const commentsContainer = post.querySelector(".comments");
  const previousScrollHeight = commentsContainer.scrollHeight;

  await loadComments(post);

  const newScrollHeight = commentsContainer.scrollHeight;
  commentsContainer.scrollTop += newScrollHeight - previousScrollHeight;
}
