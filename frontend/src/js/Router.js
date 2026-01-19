import { login } from "./pages/login.js";
import { home, createPost } from "./pages/home.js";
import { register } from "./pages/Register.js";
export async function router(path) {
    history.pushState({}, "", path);
    if (path === "/register") {
        register()
    } else if (path === "/") {
        const log = await isLogged()

        if (log) {
            home();
            return;
        }
        login();
    } else if (path === "/login") {
        login()
    } else if (path === "/createpost") {
        createPost()
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