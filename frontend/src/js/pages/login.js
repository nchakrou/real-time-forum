import { router } from "../core/Router.js"
const loginPage = `
<div class ="app-auth">
<div class="login">
<img src="/src/assets/logo.png" class = "logo">
  <h2>Login</h2>
  <p>Welcome back, please enter your details to login</p>
  <form id="login-form">
    <input type="text" id="username" placeholder = "Nickname or Email" required />
    <input type="password" id="password" placeholder = "Password" required />
    <button type="submit">Login</button>
    <p id = "error" style = "display:none"></p>
    <label for="register">Don't have an account?</label>
    <a id="register">Register</a>
  </form>
</div>
</div>`;

export function login() {
    document.body.innerHTML = loginPage;

    registerListener();
    handleLogin();
}
function handleLogin() {
    document.getElementById("login-form").addEventListener("submit", async (event) => {
        event.preventDefault();
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        const request = {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password }),
        };
        const response = await fetch("/api/login", request);

        if (response.ok) {
            init()
        } else {
            const errorElement = document.getElementById("error");
            errorElement.style.display = "block";
            errorElement.textContent = "Invalid username or password";
        }
    });
}
function registerListener() {
    document.getElementById("register").addEventListener("click", () => {
        router("/register");
    });
}


