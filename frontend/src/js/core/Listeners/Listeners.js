import { router } from "../Router.js";
import {fetchPosts} from "./postListners.js"
const routes = {
    home: () => router("/"),
    createpost: () => router("/createpost"),
    myPosts: () => router("/myPosts"),
    likedPosts: () => router("/likedPosts"),
    logout: Logout,

}


export function headerButtons() {
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
export function CategoriesListener(path = "/api/posts") {
    document.getElementById("categories").addEventListener("click", (e) => {
        if (e.target.tagName !== "LI") return;
        const categories = document.querySelectorAll(".list-categories li");
        window.history.pushState({}, "", `?category=${e.target.textContent}`);
        if (e.target.classList.contains("active")) {
            e.target.classList.remove("active");
            window.history.pushState({}, "", "/");
            router("/");
            return;
        }
        categories.forEach((cat) => cat.classList.remove("active"));

        e.target.classList.add("active")
        console.log(`${path}?category=${e.target.textContent}`);
        
       fetchPosts(`${path}?category=${e.target.textContent}`)
    });
}

