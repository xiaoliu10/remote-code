# Remote Code

English | [简体中文](./README.md)

> **🤖 This project is 100% written by AI (Claude Code)**

A remote terminal management tool built with Go + Vue 3, allowing you to monitor and control CLI programming tools remotely.

### Supported CLI Tools

Remote Code helps you remotely manage and control various AI programming assistants and terminal tools:

- **Claude Code**: Anthropic's official AI programming assistant CLI, providing intelligent code generation, refactoring, and debugging
- **Aider**: Open-source AI pair programming tool supporting GPT-4, Claude, and other models
- **Open Code**: Code assistant based on OpenAI API, providing code generation and optimization suggestions
- **Qwen Code**: Alibaba's Tongyi Qianwen code assistant, supporting Chinese and multi-language programming
- **Continue**: Open-source AI programming assistant supporting multiple LLM models
- **Other CLI Tools**: Any terminal-based command-line tools (e.g., vim, htop, irssi) can be remotely managed through Remote Code

### Why Remote Code?

When running these AI programming tools locally, you might encounter these scenarios:
- Need to remotely monitor and control programming sessions from mobile devices
- Want to continue programming tasks on your home computer while away
- Need to access and manage multiple programming sessions anytime, anywhere
- Want to use AI programming assistants on phones or tablets

Remote Code provides remote access and control of these CLI tools through a web interface and network tunneling technology, offering a smooth mobile experience and real-time terminal interaction.

## Features

- 🔐 **Secure Authentication**: JWT Token + bcrypt password encryption
- 🎨 **Beautiful UI**: Glassmorphism login page, modern design
- 🖥️ **Session Management**: Create, delete, view terminal sessions with persistence support
- 📡 **Real-time Terminal**: WebSocket streaming for terminal output
- 📜 **Terminal Scrolling**: tmux copy mode scrolling, view up to 5000 lines of history with visual scrollbar
- 📋 **Terminal Copy**: Select mode for copying terminal content, works on both desktop and mobile
- ⌨️ **Remote Control**: Send commands with realtime mode and command mode
- ⌨️ **Keyboard Shortcuts**: Ctrl+B support
- 📱 **Mobile Optimized**: Custom virtual keyboard with arrow keys, Tab, Ctrl+C, Enter, scroll mode toggle, long-press repeat support
- 🎤 **Voice Input**: Voice input support (requires HTTPS)
- 📂 **File Browser**: Built-in file explorer
- 🔒 **File Reference**: @ symbol for file references
- 🛡️ **Security**: Rate limiting, input validation, path whitelist
- 📦 **Cross-Platform Build**: Compile to standalone executables, no Go runtime needed
- 🌐 **Network Tunnel**: FRP, Tailscale, Cloudflare Tunnel
- 🌍 **i18n**: Chinese and English support
- 📜 **License**: Apache License 2.0

## What's New

### v0.0.4 - Terminal Copy Feature

**New Features**
- ✨ New select mode for copying terminal content
  - Click copy icon in toolbar to enter select mode
  - Click lines to select content to copy
  - Support select all and clear selection
  - Desktop: Shift+click for range selection
- ✨ Virtual keyboard keys support long-press repeat (arrow keys, backspace)

**Improvements**
- 🎯 Removed Ctrl+D button (prevents accidental tmux session closure)
- 🎯 Removed scroll up/down buttons (simplified UI, use scroll mode instead)
- 🎯 Optimized select mode UI for both desktop and mobile

### v0.0.3 - Voice Input & FAQ

**New Features**
- ✨ Added voice input documentation
- ✨ Added FAQ (Frequently Asked Questions)

**Improvements**
- 🎯 Fixed voice input focus issue
- 🎯 Optimized session persistence

### v0.0.2 - Terminal Scrolling & Experience Improvements

**New Features**
- ✨ Added terminal scroll mode with tmux copy mode support (default 1000 lines, configurable)
- ✨ Added visual scrollbar with mouse wheel support
- ✨ New scroll mode toggle button in virtual keyboard
- ✨ Keyboard shortcuts: Ctrl+B to toggle scroll mode
- ✨ tmux session persistence - auto-restore sessions after restart
- ✨ Configuration wizard on first run

**Improvements**
- 🎯 Improved terminal focus management with visual feedback
- 🎯 Fixed mobile viewport zoom issue after keyboard closes
- 🎯 Fixed 502 Bad Gateway issue in Docker deployment
- 🎯 Fixed multi-device connection conflict handling
- 🎯 Fixed session deletion error handling

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

**Scroll Mode**
- Click the scroll button on virtual keyboard or press **Ctrl+B** (Mac: **⌘+B**) to enter scroll mode
- Use mouse wheel or scrollbar to view history output (up to 5000 lines)
- Press **q** or click scroll button again to exit

**Copy Terminal Content**
- Click the **copy icon** in toolbar to enter select mode
- In the popup window, **click lines** to select content to copy
  - Desktop: **Shift+click** to quickly select a range
  - Mobile: Tap multiple lines to select
- Click **"Select All"** to select all content
- Click **"Copy"** button to copy selected content
- Click **"Close"** to exit select mode

Mobile virtual keyboard:
- Arrow keys (up/down/left/right), support long-press repeat
- Tab, Ctrl+C, Ctrl+L
- 📜 Scroll mode toggle
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

**Two voice input modes:**

1. **Realtime Mode** (default): Voice recognition results sent to terminal in real-time, like typing characters
2. **Command Mode**: Voice recognition results filled in command input box, edit before sending

**How to use:**

1. Click the **⚡ lightning button** to switch modes:
   - Blue highlight = realtime mode
   - Gray = command mode

2. **Command mode usage** (recommended):
   - Make sure the lightning button is gray (command mode)
   - Click the 🎤 microphone button to start voice input
   - After speaking, recognized text will be automatically filled in the command input
   - The input box will automatically focus, allowing you to edit any inaccurate recognition
   - Click "Send" or press Enter to send

**Use cases:**
- Command mode is suitable for scenarios where you need to edit the recognition result before sending
- Realtime mode is suitable for quickly inputting simple commands

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

> For more issues, see [FAQ](./docs/FAQ.md)

### Claude Code Nested Session Error

If you encounter "Nested sessions are not allowed" error when running Claude Code inside Remote Code, unset the environment variable:

```bash
# Add to ~/.zshrc or ~/.bashrc
unset CLAUDECODE
```

See [FAQ - Claude Code Nested Sessions](./docs/FAQ.md#1-在-remote-code-中运行-claude-code-报错nested-sessions-are-not-allowed) for details.

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
