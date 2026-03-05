import { init } from "../main.js";

export function ErrorPage(message = "Page Not Found", code = "404") {
  document.body.innerHTML = "";

  const errorHTML = `
        <div class="error-page-container">
            <div class="error-content">
                <h1 class="error-code">${code}</h1>
                <h2 class="error-message">Oops! ${message}</h2>
                <button id="error-home-btn" class="home-btn">
                    <span>Back to Home</span>
                </button>
            </div>
        </div>
    `;

  document.body.innerHTML = errorHTML;

  document.getElementById("error-home-btn").addEventListener("click", () => {
    init("/");
  });
}
