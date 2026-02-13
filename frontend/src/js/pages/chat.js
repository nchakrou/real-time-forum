import { headerButtons } from "../core/Listeners/Listeners.js";
import { ProfileDropdown } from "../components/ProfileDropdown.js";
import { Header } from "../components/Header.js";

const chatPage = `
${Header}
<div class = "app-chat">
</div>
`

export function chat() {
    document.body.innerHTML = chatPage
    ProfileDropdown()
    headerButtons()
}