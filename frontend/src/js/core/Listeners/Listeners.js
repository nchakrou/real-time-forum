import { router } from "../Router.js";
import {fetchPosts} from "./postListners.js"
export const routes = {
    home: () => router("/"),
    createpost: () => router("/createpost"),
    myPosts: () => router("/myPosts"),
    likedPosts: () => router("/likedPosts"),
    logout: Logout,

}


export function headerButtons() {
    document.getElementById("header-buttons").addEventListener("click", (e) => {
        const action = e.target.closest("button")?.id;
        
        if (routes[action]) routes[action]();
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

