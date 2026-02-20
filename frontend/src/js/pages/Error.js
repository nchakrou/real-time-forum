import { router } from "../core/Router.js";

export function ErrorPage(message = "Page Not Found", code = "404") {
    const errorHTML = `
        <div class="error-page-container">
            <div class="error-content">
                <h1 class="error-code">${code}</h1>
                <h2 class="error-message">Oops! ${message}</h2>
                <p class="error-description">
                    The link you followed might be broken, or the page may have been removed. 
                    Let's get you back on track.
                </p>
                <button id="error-home-btn" class="home-btn">
                    <span>Back to Home</span>
                </button>
            </div>
        </div>
    `;

    document.body.innerHTML = errorHTML;

    // Load CSS dynamically if not already present
    if (!document.getElementById('error-css')) {
        const link = document.createElement('link');
        link.id = 'error-css';
        link.rel = 'stylesheet';
        link.href = '/src/css/error.css';
        document.head.appendChild(link);
    }

    // Event Listener for the button
    document.getElementById('error-home-btn').addEventListener('click', () => {
        router("/");
    });
}
