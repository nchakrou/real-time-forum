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
  if (!path) {
    path = window.location.pathname + window.location.search;
  }
  const [user, log] = await isLogged();
  console.log(log, user, path);

  if (log && user) {
    try {
      await OpenWS();
    } catch (e) {
      ErrorPage("connection error", "500");
      return;
    }
    if (path === "/register" || path === "/login") {
      router("/");
    } else {
      router(path);
    }
  } else {
    if (path === "/register" || path === "/login") {
      router(path);
    } else if (!routes[path]) {
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
  console.log(req);

  if (req.ok) {
    const data = await req.json();
    let username = data.nickname;
    return [username, true];
  }
  return ["", false];
}
