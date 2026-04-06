import { router } from "./core/Router.js";
import { OpenWS } from "./core/WebSocket/initWs.js";
import { ErrorPage } from "./pages/Error.js";

export const routes = {
  "/": () => router("/"),
  "/createpost": () => router("/createpost"),
  "/myPosts": () => router("/myPosts"),
  "/likedPosts": () => router("/likedPosts"),
  "/chat": () => router("/chat"),
};

export async function init(path) {
  
  const pathname = path || window.location.pathname;
  if (!path) {
    path = window.location.pathname + window.location.search;
  }
  const [user, log] = await isLogged();

  if (log && user) {
    try {
      await OpenWS();
    } catch (e) {
      ErrorPage("connection error", "500");
      return;
    }
    if (pathname === "/register" || pathname === "/login") {
      router("/");
    } else {
      router(path);
    }
  } else {
    if (pathname === "/register" || pathname === "/login") {
      router(path);
    } else if (!routes[pathname]) {
      ErrorPage("404 Not Found", "404");
    } else {
      router("/login");
    }
  }
}
init();

export async function isLogged() {
  let req = await fetch("/api/islogged", {
    method: "GET",
    credentials: "include",
  });

  if (req.ok) {
    const data = await req.json();
    let username = data.nickname;
    return [username, true];
  }
  return ["", false];
}
