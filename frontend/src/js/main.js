import { router } from "./core/Router.js";
import { OpenWS } from "./core/WebSocket/initWs.js";
export const routes = {
    "/": () => router("/"),
    "/createpost": () => router("/createpost"),
    "/myPosts": () => router("/myPosts"),
    "/likedPosts": () => router("/likedPosts"),
    "/chat": () => router("/chat"),
}
const path = window.location.pathname


const log = await isLogged()
if (log) {
    OpenWS()
    if (path === "/register" || path === "/login") {
        router("/");
    } else {
        router(path);
    }
} else {
    if (path === "/register" || path === "/login") {
        router(path);
    } else if (!routes[path]) {
        alert("404 Not Found")
    } else {
        router("/login");
    }
}

async function isLogged() {
    let req = await fetch("/api/islogged", {
        method: "GET",
        credentials: "include"
    })
    if (req.ok) {
        return true
    }
    return false
}


