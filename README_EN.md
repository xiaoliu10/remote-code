# Remote Claude Code

English | [简体中文](./README.md)

A remote Claude Code management tool built with Go + Vue 3, allowing you to monitor and control Claude Code sessions remotely.

## Features

- 🔐 **Secure Authentication**: JWT Token authentication + bcrypt password encryption
- 🎨 **Beautiful UI**: Glassmorphism login page, modern UI design
- 🖥️ **Session Management**: Create, delete, and view Claude Code sessions
- 📡 **Real-time Terminal**: WebSocket streaming for terminal output
- ⌨️ **Remote Control**: Send commands to remote sessions with realtime mode and command mode
- 📱 **Mobile Optimized**: Custom virtual keyboard with arrow keys, Tab, Ctrl+C shortcuts
- 📂 **File Browser**: Built-in file explorer with file content viewing
- 🔒 **File Reference**: Support @ symbol for file reference functionality
- 🛡️ **Security**: Rate limiting, input validation, path whitelist
- 🐳 **Docker Deployment**: One-click deployment, ready to use
- 🌐 **Network Tunnel**: Integrated tunnel solutions (Frp, Tailscale, Cloudflare Tunnel)
- 🌍 **i18n**: Chinese and English language support
- 📜 **License**: Apache License 2.0

## Tech Stack

### Backend
- Go 1.21+
- Gin Web Framework
- gorilla/websocket
- JWT Authentication
- tmux Session Management

### Frontend
- Vue 3 + TypeScript
- Vite
- Naive UI
- xterm.js
- Pinia
- Vue Router

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 20+
- tmux
- (Optional) Docker & Docker Compose

### Option 1: One-Command Start (Recommended)

```bash
# 1. Clone the project
git clone https://github.com/yourname/remote-claude-code.git
cd remote-claude-code

# 2. Create config file
cp config.ini.example config.ini

# 3. Edit config (optional, defaults work for quick start)
# vim config.ini

# 4. Start all services
./start.sh

# Visit http://localhost:5173
```

The startup script automatically:
- Checks dependencies (tmux, go, node)
- Generates `.env` files for backend and frontend
- Starts backend service (default port 9090)
- Starts frontend service (default port 5173)
- Starts FRP tunnel if configured

**Management Commands:**
```bash
./start.sh          # Start services
./start.sh --no-frp # Start services (disable FRP)
./start.sh --frp    # Start services (force enable FRP)
./stop.sh           # Stop services
```

### Option 2: Docker Compose

```bash
# Configure environment
cp config.ini.example config.ini

# Start services
docker-compose up -d

# Visit http://localhost
```

### Option 3: Enable Network Tunnel

Edit `config.ini` and set:

```ini
# Enable FRP
FRP_ENABLED=true
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-secure-token
```

Then run `./start.sh`. The system will automatically:
- Download FRP client (if not installed)
- Generate FRP configuration
- Start network tunnel

For detailed tunnel configuration: [docs/NETWORK_TUNNEL_GUIDE.md](./docs/NETWORK_TUNNEL_GUIDE.md)

## Usage

### 1. Login

Default admin account:
- Username: `admin`
- Password: `admin123` (please change in production)

### 2. Create Session

On the Dashboard page:
1. Enter session name (e.g., `my-project`)
2. Select working directory (optional)
3. Click "Create" to create session

### 3. Open Terminal

Click the "Open" button in the session list to open the real-time terminal.

### 4. Send Commands

Terminal supports two input modes:
- **Realtime Mode**: Each character sent immediately, supports autocomplete
- **Command Mode**: Press Enter to send complete command

Mobile virtual keyboard includes:
- Arrow keys (up/down/left/right)
- Tab key
- Ctrl+C (interrupt)
- Ctrl+D (exit)
- @ symbol (file reference)

### 5. File Browser

Click the file icon in the sidebar to browse files in the working directory.

## Port Configuration

Default ports used by the project:

| Service | Default Port | Environment Variable |
|---------|--------------|---------------------|
| Backend API | 9090 | `PORT` |
| Frontend Dev Server | 5173 | - |
| Docker Frontend | 80 | - |

To change the backend port, set `PORT` in `backend/.env` and `VITE_BACKEND_PORT` in `frontend/.env`.

## API Documentation

### Authentication

```bash
# Login
POST /api/auth/login
{
  "username": "admin",
  "password": "admin123"
}
```

### Session Management

```bash
# List all sessions
GET /api/sessions

# Create session
POST /api/sessions
{
  "name": "my-project",
  "work_dir": "/home/user/projects"
}

# Get session details
GET /api/sessions/{name}

# Delete session
DELETE /api/sessions/{name}

# Send command
POST /api/sessions/{name}/command
{
  "command": "ls -la"
}

# Get output
GET /api/sessions/{name}/output
```

### File Operations

```bash
# List directory contents
GET /api/files?path=/home/user/projects

# Get file content
GET /api/files/content?path=/home/user/projects/file.txt

# Create file/directory
POST /api/files
{
  "path": "/home/user/projects/newfile.txt",
  "type": "file",
  "content": "Hello World"
}

# Rename
PUT /api/files/rename
{
  "oldPath": "/home/user/projects/old.txt",
  "newPath": "/home/user/projects/new.txt"
}

# Delete
DELETE /api/files?path=/home/user/projects/file.txt
```

### WebSocket

```javascript
// Connect to session
const ws = new WebSocket('ws://localhost:9090/ws/session-name?token=YOUR_JWT_TOKEN')

// Send command
ws.send(JSON.stringify({
  type: 'command',
  data: 'your-command'
}))

// Send keys (realtime mode)
ws.send(JSON.stringify({
  type: 'keys',
  data: 'ls'
}))

// Receive output
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  if (message.type === 'output') {
    console.log(message.data.text)
  }
}
```

## Security Recommendations

1. **Change Default Password**: Modify admin password before first use (`ADMIN_PASSWORD` in `backend/.env`)
2. **Use Strong JWT Secret**: Generate a random string for `JWT_SECRET`
3. **Configure HTTPS**: Use HTTPS in production environments
4. **Restrict Access**: Use firewall or VPN to limit access
5. **Regular Updates**: Keep dependencies up to date

## Directory Structure

```
remote-claude-code/
├── backend/              # Go backend
│   ├── cmd/
│   │   └── server/       # Entry point
│   ├── internal/         # Internal packages
│   │   ├── api/          # API handlers
│   │   ├── auth/         # Authentication
│   │   ├── config/       # Configuration
│   │   ├── security/     # Security
│   │   ├── tmux/         # tmux management
│   │   └── websocket/    # WebSocket
│   └── Dockerfile
├── frontend/             # Vue frontend
│   ├── src/
│   │   ├── api/          # API client
│   │   ├── components/   # Components
│   │   ├── composables/  # Composables
│   │   ├── router/       # Router
│   │   ├── stores/       # State management
│   │   └── views/        # Pages
│   ├── Dockerfile
│   └── nginx.conf
├── frp/                  # FRP tunnel configuration
├── docker-compose.yml
└── README.md
```

## Troubleshooting

### tmux Not Available
```bash
# macOS
brew install tmux

# Ubuntu/Debian
sudo apt-get install tmux

# CentOS/RHEL
sudo yum install tmux
```

### WebSocket Connection Failed
- Check if JWT Token is valid
- Confirm session name is correct
- Check browser console for errors
- Check if port is occupied or firewall settings

### Empty Session Output
- Confirm session has started claude-code
- Check tmux session status: `tmux list-sessions`

### Login Failed via Domain
- Ensure Vite config `allowedHosts` includes your domain
- Check if backend API responds correctly

## License

Apache License 2.0

## Contributing

Issues and Pull Requests are welcome!
