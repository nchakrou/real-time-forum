
import { headerButtonsRoutes as routes } from "../core/Listeners/Listeners.js";
export function ProfileDropdown() {
    const profileTrigger = document.getElementById('user-profile');
    const userDropdown = document.getElementById('user-dropdown');
    const notificationBtn = document.getElementById('notification-btn');
    const notificationDropdown = document.getElementById('notification-dropdown');
    if (!profileTrigger || !userDropdown || !notificationBtn || !notificationDropdown) return;

    // Toggle dropdown
    profileTrigger.addEventListener('click', (e) => {
        closeDropdown(notificationBtn, notificationDropdown);
        e.stopPropagation();
        profileTrigger.classList.toggle('active');
        userDropdown.classList.toggle('hidden');
    });

    userDropdown.addEventListener("click", (e) => {

        const action = e.target.closest("button")?.id;

        if (routes[action]) {
            routes[action]();
            closeDropdown(profileTrigger, userDropdown);
            closeDropdown(notificationBtn, notificationDropdown);
        }
    });
    document.addEventListener('click', (e) => {



        if (!profileTrigger.contains(e.target) && !userDropdown.contains(e.target) && !notificationBtn.contains(e.target) && !notificationDropdown.contains(e.target)) {
            closeDropdown(profileTrigger, userDropdown);
            closeDropdown(notificationBtn, notificationDropdown);
        }
    });

    notificationBtn.addEventListener('click', (e) => {
        closeDropdown(profileTrigger, userDropdown);
        e.stopPropagation();
        notificationBtn.classList.toggle('active');
        notificationDropdown.classList.toggle('hidden');
    });

    function closeDropdown(rm, add) {
        rm.classList.remove('active');
        add.classList.add('hidden');
    }
}
