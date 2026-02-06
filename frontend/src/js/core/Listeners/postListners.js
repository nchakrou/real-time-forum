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
