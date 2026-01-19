const loginPage = `
<div class ="app-auth">
<div class="login">
<img src="/src/assets/logo.png" 
 style="width: 100px; margin: 0 auto;">
  <h2>Login</h2>
  <p>Welcome back, please enter your details to login</p>
  <form id="login-form">
    <input type="text" id="username" placeholder = "Nickname or Email" required />
    <input type="password" id="password" placeholder = "Password" required />
    <button type="submit">Login</button>
    <p id = "error" style = "display:none"></p>
    <label for="register">Don't have an account?</label>
    <a href="/register">Register</a>
  </form>
</div>
</div>`;
const registerPage = `
  <div class ="app-auth">
<div class="register">
<img src="/src/assets/logo.png" 
 style="width: 100px; margin: 0 auto;">
  <h2>Create Account</h2>
  
  <p>Sing Up to get started</p>
  <form id="register-form">
  <div id = fullname class = "fullname-auth">
    <input type="text" id="firstname" name="firstname" placeholder="First Name" required maxlength="20" />
    <input type="text" id="lastname" name="lastname" placeholder="Last Name" required maxlength="20"/>
  </div>
    <input type="text" id="nickname" name="nickname" placeholder="Nickname" required />
    <input type="text" id="age" name="age" placeholder="Age" required />

    <input type="email" id="email" name="email" placeholder="Email" required />
    <input type="password" id="password" name="password" placeholder="Password" required />
    <select id="gender" name="gender" aria-placeholder="Gender" required>
      <option value="male">Male</option>
      <option value="female">Female</option>
    </select>
    <div class="message">
    Already have an account?
  <a href="/login">login</a>
    </div>
    <button class = "auth-button"type="submit">Register</button>
  </form>
</div>

</div>`;
const homePage = `
<header>
<h2>Forum</h2>
<div>
<button>Create Post</button>
  <button>created posts</button>
  <button>liked posts</button>
  <button>logout</button>
</div>
</header>
<div class = "app-home">
<div class = "categories">
  <h2>Categories</h2>
  <ul class="list-categories">
  <li>FPS</li>
  <li>Battle Royale</li>
  <li>MOBA</li>
  <li>Esports</li>
  <li>RPG</li>
  <li>Strategy</li>
  <li>Simulation</li>
  </ul>
</div>
<div class = "posts">
  <h2>Posts</h2>
  <p>No posts yet</p>
</div>
<div class = "users">
  <h2>Users</h2>
  <p>No users yet</p>
</div>
</div>
`
function navigateTo(path) {
  if (path === "/register") {
    document.body.innerHTML = registerPage;
  } else if (path === "/") {
    document.body.innerHTML = loginPage;
  } else if (path === "/home") {
    document.body.innerHTML = homePage;
  }
}
navigateTo(window.location.pathname);
document
  .getElementById("login-form")
  ?.addEventListener("submit", async (event) => {
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
    const response = await fetch("/check", request);

    if (response.ok) {
      history.pushState({}, "", "/home");
      navigateTo("/home");
    } else {
      const errorElement = document.getElementById("error");
      errorElement.style.display = "block";
      errorElement.textContent = "Invalid username or password";
    }
  });
