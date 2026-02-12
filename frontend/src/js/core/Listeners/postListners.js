import { toggleLike } from "../../pages/likes.js";

const postButtons = {
    like_button: (e) => likeListener(e),
    dislike_button: (e) => dislikeListener(e),
    comment_button: (e) => commentListener(e),
    Submit_comment: (e) => submitCommentListener(e)
}

export async function fetchPosts(path){   
    try{     
        const response = await fetch(path, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            
            if (response.ok) {
                const posts = await response.json();
                createpostsContainer(posts);
            } else {
               throw new Error(response);
            }

    }catch(error){
        alert("ok",error);
}
}
export function createpostsContainer(posts) {
    const postsContainer = document.getElementById("posts-container")
    postsContainer.innerHTML = ""
    if (!posts) {

        const noPosts = document.createElement("p")
        noPosts.textContent = "No posts yet"
        postsContainer.appendChild(noPosts)
    } else {
        console.log(posts);

        posts.forEach((post) => {
            const postElement = document.createElement("div")
            postElement.classList.add("post")
            postElement.dataset.postId = post.id
            postElement.innerHTML = `
          <h3>${post.title}</h3>
          <p>${post.content}</p>
          <p><strong>Category:</strong> ${post.categories}</p>
          <div class= "post-buttons">
          <button type="submit" class = "like_button">👍 <span>${post.likes}</span> Like</button>
          <button type="submit" class = "dislike_button">👎 <span>${post.dislikes}</span> Dislike</button>
          <button type="submit" class = "comment_button">💬<span> ${post.comments}</span> comment</button>
          </div>
          <div class = "comments-section hidden">
          <input class = "comment-input" type="text" placeholder="Write a comment...">
          <button type="submit" class = "Submit_comment">Submit</button>
          </div>
        `
            postsContainer.appendChild(postElement)


        })
        PostButtonsListener()
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
    if (!commentText) return alert("Comment cannot be empty");

    const postId = post.dataset.postId;

    try {
        const formData = new FormData();
        formData.append("post_id", postId);
        formData.append("comment", commentText);

        const res = await fetch("/add-comment", {
            method: "POST",
            credentials: "include",
            body: formData
        });

        if (!res.ok) return console.error("Failed to add comment");

        const data = await res.json();
        if (data.status === "ok") {
            commentInput.value = "";
            loadComments(post);
        }

    } catch (err) {
        console.error("Error adding comment:", err);
    }
}


async function likeListener(e) {
    const post = e.target.closest(".post");
    const postId = post.dataset.postId; 
    console.log("like click, postid =", postId);
    await toggleLike(postId, 1);        
}

async function dislikeListener(e) {
    const post = e.target.closest(".post");
    const postId = post.dataset.postId;
    await toggleLike(postId, -1);       
}
async function loadComments(post) {
    const commentsContainer = post.querySelector(".comments-section .comments-list");
    const postId = post.dataset.postId;
    
    try {
        const res = await fetch(`/comments?post_id=${postId}`);
        if (!res.ok) throw new Error("Failed to fetch comments");

        const comments = await res.json();

        commentsContainer.innerHTML = ""; 

        comments.forEach(c => {
            const newComment = document.createElement("p");
            newComment.innerHTML = `<strong>${c.username}:</strong> ${c.content}`;
            commentsContainer.appendChild(newComment);
        });

    } catch (err) {
        console.error("Error loading comments:", err);
    }
}


function commentListener(e) {
    const post = e.target.closest(".post");
    const section = post.querySelector(".comments-section");
    section.classList.toggle("hidden");

    if (!section.classList.contains("hidden") && section.querySelector(".comments-list").children.length === 0) {
        loadComments(post);
    }
}
