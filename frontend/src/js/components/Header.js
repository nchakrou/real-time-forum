export const Header = `
<header>
  
        <div class="header-text" >
            <h1 class="logo-text">Forum</h1>
        </div>
        <nav>
        <div class="header-buttons" id = "header-buttons">
            <button id="home" title="Home" class="nav-btn">
                <img src="/src/assets/home.svg">
                <span class="nav-text">Home</span>
            </button>
            <button id="createpost" title="Create Post" class="nav-btn">
                <img src="src/assets/plus.svg">
                <span class="nav-text">Create</span>
            </button>
            
           
            <button id="chat" title="Chat" class="nav-btn">
                <img src="src/assets/chat.svg">
                <span class="nav-text">Chat</span>
            </button>
       
        </div>
        </nav>
        <div class="header-actions">
            <div class="notification-wrapper">
                <button id="notification-btn" class="notification-btn" title="Notifications">
                    <img src="src/assets/bell.svg" alt="Notifications">
                    <span id="notification-badge" class="notification-badge hidden">0</span>
                </button>
                
                
                <div id="notification-dropdown" class="notification-dropdown hidden">
                    <div class="notification-dropdown-header">
                        <h3>Notifications</h3>
                    </div>
                    <div id="notification-list" class="notification-list">
                    </div>
                    <div class="notification-dropdown-footer">
                        <a href="#" id="view-all-notifications">View all notifications</a>
                    </div>
                </div>
            </div>
           
        
            
                <div class="user-profile" id="user-profile">
                    <div class="user-avatar" id="user-avatar">
                      <span id="user-initials">U</span>
                    </div>
                    
                    <div id="user-dropdown" class="user-dropdown hidden">
                    <div class="user-header">
                        <div class="user-avatar">
                            <span id="user-initials-dropdown">U</span>
                        </div>
                        <div class="user-info">
                            <span id="user-name-dropdown" class="user-name-large">Username</span>
                            
                            <span class="user-status-dropdown">Online</span>
                        </div>
                    </div>
                    <div class= "divider"></div>
                  
                    <div class="user-dropdown-menu">
                        <button id="myPosts" class="dropdown-menu-item">
                            <img src="src/assets/myposts.svg" alt="My Posts">
                            <span>My Posts</span>
                        </button>
                        <button id="likedPosts" class="dropdown-menu-item">
                            <img src="src/assets/heart.svg" alt="Liked Posts">
                            <span>Liked Posts</span>
                        </button>
                    <div class= "divider"></div>
                        <button id="logout" class="dropdown-menu-item logout-item">
                            <img src="src/assets/log-out.svg" alt="Logout">
                            <span>Logout</span>
                        </button>
                         
                    </div>
                </div>
            </div>
            
           
       
    </div>
</header>
`;
