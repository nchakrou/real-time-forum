import { router } from "./Router.js";


const routes = {
    home: () => router("/"),
    createpost: () => router("/createpost"),
    myPosts: () => router("/myPosts"),
    likedPosts: () => router("/likedPosts"),
    logout: Logout,

}

const postButtons = {
    like_button: (e) => likeListener(e),
    dislike_button: (e) => dislikeListener(e),
    comment_button: (e) => commentListener(e),
    Submit_comment: (e) => submitCommentListener(e)
}
export function headerButtons() {
    console.log("ok");

    document.getElementById("header-buttons").addEventListener("click", (e) => {
        console.log(e.target.id);

        if (routes[e.target.id]) routes[e.target.id]();
    });
}
async function Logout() {
    await fetch("/api/logout", {
        method: "GET"
    })
    router("/login")
}
export function CategoriesListener() {
    document.getElementById("categories").addEventListener("click", (e) => {
        if (e.target.tagName !== "LI") return;
        const categories = document.querySelectorAll(".list-categories li");
        window.history.pushState({}, "", `?category=${e.target.textContent}`);
        if (e.target.classList.contains("active")) {
            e.target.classList.remove("active");
            window.history.pushState({}, "", "/");
            return;
        }
        categories.forEach((cat) => cat.classList.remove("active"));
        
       e.target.classList.add("active")
    });
}
export function PostButtonsListener() {
   const postsContainer = document.getElementById("posts-container");
    postsContainer.addEventListener("click", (e) => {
        const button = e.target.closest("button");
        if (!button) return;
        const buttonClass = button.classList[0];
        if (postButtons[buttonClass]) postButtons[buttonClass](e);
        
        

        
        

    })
}
function submitCommentListener(e) {
    const post = e.target.closest(".post");
    const commentInput = post.querySelector(".comment-input");
    
    const commentText = commentInput.value.trim();
    if (commentText === "") {
        alert("Comment cannot be empty");
    }
    const commentsSection = post.querySelector(".comments-section");
    const newComment = document.createElement("p");
    newComment.textContent = commentText;
    commentsSection.appendChild(newComment);
    commentInput.value = "";
}
function likeListener(e) {
    const button = e.target.closest("button");
    button.querySelector("span").textContent++
    
}
function dislikeListener(e) {
    const button = e.target.closest("button");
    console.log(button.querySelector("span").textContent);
    button.querySelector("span").textContent++
    
}
function commentListener(e) {
    const post = e.target.closest(".post");
    post.querySelector(".comments-section").classList.toggle("hidden")
}
  
