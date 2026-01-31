import { router } from "./Router.js";


const routes = {
    home: () => router("/"),
    createpost: () => router("/createpost"),
    myPosts: () => router("/myPosts"),
    likedPosts: () => router("/likedPosts"),
    logout: Logout,
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
        console.log(e.target);

       e.target.classList.add("active")
    });
}
