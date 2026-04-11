# Real-Time Forum

A modern, real-time web application for discussions, featuring a Go backend and a responsive vanilla JavaScript frontend.

## 🚀 Features

- **User Authentication**: Secure Login, Registration, and Session Management.
- **Discussion Posts**: Create and view posts across different categories.
- **Interactions**:
  - Add comments to posts.
  - Like/Dislike posts and comments.
- **Real-Time Updates**: Using WebSockets for instant notifications and updates.
- **Mobile Responsive**: Fully optimized for a seamless experience on all devices.
- **Personalized Feed**: View your own posts and posts you've liked.

## 🛠️ Technologies Used

- **Backend**: [Go](https://golang.org/) (Standard Library)
- **Real-Time**: [Gorilla WebSocket](https://github.com/gorilla/websocket)
- **Database**: [SQLite](https://www.sqlite.org/)
- **Frontend**: HTML5, CSS3 (Vanilla), JavaScript (Vanilla ES6+)
- **Architecture**: Single Page Application (SPA).

## 📋 Prerequisites

- [Go](https://golang.org/doc/install) (version 1.16 or higher)
- [GCC](https://gcc.gnu.org/) (required for SQLite driver)

## ⚙️ Installation & Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/nchakrou/real-time-forum.git
   cd real-time-forum
   ```

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Run the application**:

   ```bash
   go run main.go
   ```

4. **Access the application**:
   Open your browser and navigate to `http://localhost:8081`.

## 📂 Project Structure

- `backend/`: Go source code for handlers, database initialization, and middleware.
- `frontend/`: HTML, CSS, and JavaScript source code.
  - `src/js/`: Modular JavaScript structure (components, pages, core logic).
  - `src/css/`: Vanilla CSS styling.
- `main.go`: Application entry point and route definitions.


