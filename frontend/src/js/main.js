import { router } from "./core/Router.js";

const path =window.location.pathname

    
const log = await isLogged() 
if(log) {
    if (path !== "/register" && path !== "/login") {
        router(path);
    }else {
        router("/");
    }
}else {

    if (path === "/regstier" || path === "/login") {
        router(path);
    }else {
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

