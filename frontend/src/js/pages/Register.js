import { router } from "../core/Router.js";
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
  <a id="login">login</a>
    </div>
    <button class = "auth-button"type="submit">Register</button>
  </form>
</div>

</div>`;
export function register() {
  document.body.innerHTML = registerPage
  loginListener()
}
function loginListener() {
  document.getElementById("login").addEventListener("click", () => {
    router("/login");
  });
}
