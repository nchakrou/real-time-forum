import { login } from "../pages/login.js";
import { home } from "../pages/home.js";
import { register } from "../pages/Register.js";
import { createPost } from "../pages/creatPost.js";
import { myPosts } from "../pages/myPosts.js";
import { likedPosts } from "../pages/likes.js";
import { chat } from "../pages/chat.js";
export function router(path) {
    history.pushState({}, "", path);
    if (path === "/register") {
        register()
    } else if (path === "/") {
        home();
    } else if (path === "/login") {
        login()
    } else if (path === "/createpost") {
        createPost()
    } else if (path === "/myPosts") {
        myPosts()
    } else if (path === "/likedPosts") {
        likedPosts()
    } else if (path === "/chat") {
        chat()
    }
}
