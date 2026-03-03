export const Popup = {
    show: (message, duration = 3000) => {
        const existingPopup = document.getElementById('popup');
        if (existingPopup) {
            existingPopup.remove();
        }

        const popup = document.createElement('div');
        popup.id = 'popup';
        popup.className = 'popup';

        const content = document.createElement('div');
        content.id = 'popup-content';

        const text = document.createElement('p');
        text.textContent = message;

        content.appendChild(text);
        popup.appendChild(content);
        document.body.appendChild(popup);

        setTimeout(() => {
            popup.classList.add('show');
        }, 10);
        setTimeout(() => {
            popup.classList.remove('show');
            setTimeout(() => {
                popup.remove();
            }, 500);
        }, duration);

        popup.onclick = () => {
            popup.classList.remove('show');
            setTimeout(() => {
                popup.remove();
            }, 500);
        };
    }
};
