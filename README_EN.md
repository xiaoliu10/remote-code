# Remote Code

English | [简体中文](./README.md)

A remote terminal management tool built with Go + Vue 3, allowing you to monitor and control CLI programming tools (Claude Code, Cursor, Aider, etc.) remotely.

## Features

- 🔐 **Secure Authentication**: JWT Token + bcrypt password encryption
- 🎨 **Beautiful UI**: Glassmorphism login page, modern design
- 🖥️ **Session Management**: Create, delete, and view terminal sessions
- 📡 **Real-time Terminal**: WebSocket streaming for terminal output
- ⌨️ **Remote Control**: Send commands with realtime mode and command mode
- 📱 **Mobile Optimized**: Custom virtual keyboard with arrow keys, Tab, Ctrl+C, Enter
- 🎤 **Voice Input**: Voice input support (requires HTTPS)
- 📂 **File Browser**: Built-in file explorer
- 🔒 **File Reference**: @ symbol for file references
- 🛡️ **Security**: Rate limiting, input validation, path whitelist
- 📦 **Cross-Platform Build**: Compile to standalone executables, no Go runtime needed
- 🌐 **Network Tunnel**: FRP, Tailscale, Cloudflare Tunnel
- 🌍 **i18n**: Chinese and English support
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

- Go 1.21+ (only needed for source code, not required for compiled version)
- Node.js 20+ (needed for frontend development)
- tmux

### Option 1: Run from Source

```bash
# 1. Clone the project
git clone https://github.com/xiaoliu10/remote-code.git
cd remote-code

# 2. Start services (auto-creates config on first run)
./start.sh

# Visit http://localhost:5173
```

### Option 2: Build and Run (Recommended for Production)

```bash
# 1. Clone the project
git clone https://github.com/xiaoliu10/remote-code.git
cd remote-code

# 2. Build backend (generates multi-platform executables)
./build.sh

# 3. Start services (will use compiled binary automatically)
./start.sh --frp

# Visit http://localhost:5173
```

**Build Outputs:**

| Platform | Filename | Architecture |
|----------|----------|--------------|
| macOS Intel | `remote-code-macos-intel` | x86_64 |
| macOS Apple Silicon | `remote-code-macos-apple` | arm64 |
| Linux x64 | `remote-code-linux-x64` | x86_64 |
| Linux ARM64 | `remote-code-linux-arm64` | arm64 |
| Windows x64 | `remote-code-windows-x64.exe` | x86_64 |

On first run, the system will automatically:
- Create config directory at `~/.remote-code/`
- Generate random admin password (save it!)
- Start backend service (port 9090)
- Start frontend service (port 5173)

### Management Commands

```bash
./build.sh          # Build multi-platform executables
./start.sh          # Start services
./start.sh --frp    # Start with FRP enabled
./start.sh --no-frp # Start without FRP
./stop.sh           # Stop services
```

### Configuration

Config file located at `~/.remote-code/config.ini`:

```ini
# Backend config
BACKEND_PORT=9090
JWT_SECRET=auto-generated
ADMIN_PASSWORD=auto-generated
ALLOWED_DIR=/home/user/projects

# FRP config
FRP_ENABLED=false
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-token
```

## Remote Access

### Why Do You Need Tunneling?

Remote Code runs on your local network by default and is only accessible within the LAN. When you need to **access your home Claude Code sessions remotely from outside (office, cafe, mobile network)**, you need a tunneling solution.

**Use Cases:**
- Running Claude Code on home PC, need remote control when away
- Office network blocks external ports
- No public IP or ISP blocks ports
- Need stable remote access experience

### Option 1: FRP + Nginx HTTPS (Recommended)

FRP (Fast Reverse Proxy) is a high-performance tunneling tool that exposes local services through a public server.

#### 1. Configure FRP

Edit `~/.remote-code/config.ini`:

```ini
FRP_ENABLED=true
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-token
```

#### 2. Configure Nginx on Server

Create Nginx config on your server:

```nginx
server {
    listen 8444 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # WebSocket
    location /api/ws {
        proxy_pass http://127.0.0.1:9090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400s;
    }

    # Backend API
    location /api {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
    }

    # Frontend
    location / {
        proxy_pass http://127.0.0.1:5173;
        proxy_set_header Host $host;
    }
}
```

#### 3. Start Services

```bash
./start.sh --frp
```

#### 4. Access

```
https://your-domain.com:8444
```

### Option 2: Tailscale

```bash
# Install and login to Tailscale
brew install tailscale
sudo tailscale up

# Start services
./start.sh

# Visit http://<tailscale-ip>:5173
```

### Option 3: Cloudflare Tunnel

```bash
export CLOUDFLARE_TUNNEL_TOKEN=your-token
docker-compose --profile cloudflare up -d
```

## Usage

### 1. Login

Random password is generated on first run. Check startup output for the password.

### 2. Create Session

1. Enter session name (e.g., `my-project`)
2. Select working directory (optional)
3. Click "Create"

### 3. Terminal

Two input modes:
- **Realtime Mode**: Each character sent immediately, supports autocomplete
- **Command Mode**: Press Enter to send complete command

Mobile virtual keyboard:
- Arrow keys (up/down/left/right)
- Tab, Ctrl+C, Ctrl+D
- @ symbol (file reference)
- 🎤 Voice input

### 4. Voice Input

⚠️ **Voice input requires HTTPS**

Voice input has browser security restrictions:

| Method | Description |
|--------|-------------|
| HTTPS | Use Nginx + Let's Encrypt |
| localhost | No config needed for local access |
| Chrome flag | Testing only, visit `chrome://flags/#unsafely-treat-insecure-origin-as-secure` |

### 5. File Browser

Click the file icon in sidebar to browse files.

## Directory Structure

```
~/.remote-code/          # Config directory
├── config.ini                  # Main config
├── frpc.ini                    # FRP config (auto-generated)
├── frpc                        # FRP client (auto-downloaded)
└── logs/                       # Logs
    ├── backend.log
    ├── frontend.log
    └── frp.log

remote-code/             # Source directory
├── backend/                    # Go backend
├── frontend/                   # Vue frontend
├── desktop/                    # Desktop app
├── nginx/                      # Nginx config
├── docs/                       # Documentation
├── start.sh                    # Start script
└── stop.sh                     # Stop script
```

## API Documentation

### Authentication

```bash
POST /api/auth/login
{"username": "admin", "password": "your-password"}
```

### Session Management

```bash
GET    /api/sessions              # List sessions
POST   /api/sessions              # Create session
GET    /api/sessions/{name}       # Get details
DELETE /api/sessions/{name}       # Delete session
POST   /api/sessions/{name}/command  # Send command
```

### File Operations

```bash
GET    /api/files?path=/path     # List directory
GET    /api/files/content?path=/path  # Get content
POST   /api/files                # Create file/folder
PUT    /api/files/rename         # Rename
DELETE /api/files?path=/path     # Delete
```

### WebSocket

```javascript
const ws = new WebSocket('wss://your-domain:8444/api/ws/session?token=TOKEN')

// Send command
ws.send(JSON.stringify({type: 'command', data: 'ls -la'}))

// Send keys
ws.send(JSON.stringify({type: 'keys', data: 'ls'}))
```

## Troubleshooting

### Port in Use

```bash
./stop.sh  # Stops all services including orphan processes
```

### WebSocket Connection Failed

- Check Nginx `/api/ws` configuration
- Verify HTTPS certificate is valid
- Check browser console for errors

### Voice Input Not Working

- Ensure HTTPS access
- Or use localhost for testing
- Check browser Web Speech API support

## Security Recommendations

1. **Save Initial Password**: Save the randomly generated password on first run
2. **Use HTTPS**: Always use HTTPS in production
3. **Regular Updates**: Keep dependencies up to date
4. **Restrict Access**: Use firewall to limit access

## License

Apache License 2.0

## Contributing

Issues and Pull Requests are welcome!
