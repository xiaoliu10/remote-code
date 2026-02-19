# Remote Claude Code

[English](./README_EN.md) | 简体中文

基于 Go + Vue 3 的远程 Claude Code 管理工具，可以远程监控和控制家中的 Claude Code 会话。

## 功能特性

- 🔐 **安全认证**: JWT Token 认证 + bcrypt 密码加密
- 🎨 **精美界面**: 玻璃拟态风格登录页，现代化 UI 设计
- 🖥️ **会话管理**: 创建、删除、查看 Claude Code 会话
- 📡 **实时终端**: WebSocket 实时流式传输终端输出
- ⌨️ **远程控制**: 发送命令到远程会话，支持实时模式和命令模式
- 📱 **移动端优化**: 自定义虚拟键盘，包含方向键、Tab、Ctrl+C 等快捷键
- 📂 **文件浏览**: 内置文件浏览器，支持浏览、查看文件内容
- 🔒 **文件引用**: 支持 @ 符号触发文件引用功能
- 🛡️ **安全防护**: 速率限制、输入验证、路径白名单
- 🐳 **Docker 部署**: 一键部署，开箱即用
- 🌐 **内网穿透**: 集成多种内网穿透方案（Frp、Tailscale、Cloudflare Tunnel）
- 🌍 **国际化**: 支持中英文切换
- 📜 **开源协议**: Apache License 2.0

## 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- gorilla/websocket
- JWT 认证
- tmux 会话管理

### 前端
- Vue 3 + TypeScript
- Vite
- Naive UI
- xterm.js
- Pinia
- Vue Router

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 20+
- tmux
- （可选）Docker & Docker Compose

### 方案一：一键启动（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/yourname/remote-claude-code.git
cd remote-claude-code

# 2. 创建配置文件
cp config.ini.example config.ini

# 3. 编辑配置（可选，使用默认配置即可体验）
# vim config.ini

# 4. 启动所有服务
./start.sh

# 访问 http://localhost:5173
```

启动脚本会自动：
- 检查依赖（tmux, go, node）
- 生成后端和前端的 `.env` 文件
- 启动后端服务（默认端口 9090）
- 启动前端服务（默认端口 5173）
- 如果配置了 FRP，自动启动内网穿透

**管理命令：**
```bash
./start.sh          # 启动服务
./start.sh --no-frp # 启动服务（禁用 FRP）
./start.sh --frp    # 启动服务（强制启用 FRP）
./stop.sh           # 停止服务
```

### 方案二：Docker Compose

```bash
# 配置环境变量
cp config.ini.example config.ini

# 启动服务
docker-compose up -d

# 访问 http://localhost
```

### 方案三：启用内网穿透

编辑 `config.ini`，设置以下选项：

```ini
# 启用 FRP
FRP_ENABLED=true
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-secure-token
```

然后运行 `./start.sh`，系统会自动：
- 下载 FRP 客户端（如果未安装）
- 生成 FRP 配置文件
- 启动内网穿透

详细内网穿透配置请查看：[docs/NETWORK_TUNNEL_GUIDE.md](./docs/NETWORK_TUNNEL_GUIDE.md)

#### 使用 Tailscale（最简单）

```bash
# 1. 安装 Tailscale
brew install tailscale  # macOS
# 或 curl -fsSL https://tailscale.com/install.sh | sh  # Linux

# 2. 登录
sudo tailscale up

# 3. 启动服务
./start-with-tunnel.sh tailscale

# 4. 在任何设备上访问（需安装 Tailscale）
# http://<tailscale-ip>:5173
```

#### 使用 Frp

```bash
# 1. 编辑 frp/frpc.ini，填入服务器信息
# 2. 启动服务
docker-compose --profile frp up -d

# 或使用启动脚本
./start-with-tunnel.sh frp
```

#### 使用 Cloudflare Tunnel

```bash
# 1. 获取 Tunnel Token
# https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/

# 2. 启动服务
export CLOUDFLARE_TUNNEL_TOKEN=your-token
docker-compose --profile cloudflare up -d

# 或
./start-with-tunnel.sh cloudflare
```

**详细内网穿透配置请查看**: [frp/README.md](./frp/README.md)

## 使用说明

### 1. 登录系统

默认管理员账号:
- 用户名: `admin`
- 密码: `admin123`（生产环境请务必修改）

### 2. 创建会话

在 Dashboard 页面:
1. 输入会话名称（如: `my-project`）
2. 选择工作目录（可选）
3. 点击 "Create" 创建会话

### 3. 打开终端

点击会话列表中的 "Open" 按钮，即可打开实时终端。

### 4. 发送命令

终端支持两种输入模式：
- **实时模式**: 每个字符实时发送，支持自动补全
- **命令模式**: 按 Enter 发送完整命令

移动端提供虚拟键盘，包含：
- 方向键（上下左右）
- Tab 键
- Ctrl+C（中断）
- Ctrl+D（退出）
- @ 符号（文件引用）

### 5. 文件浏览

点击侧边栏的文件图标，可以浏览工作目录中的文件。

## 端口配置

项目使用以下默认端口：

| 服务 | 默认端口 | 环境变量 |
|------|----------|----------|
| 后端 API | 9090 | `PORT` |
| 前端开发服务器 | 5173 | - |
| Docker 前端 | 80 | - |

如需修改后端端口，在 `backend/.env` 中设置 `PORT`，同时在 `frontend/.env` 中设置 `VITE_BACKEND_PORT`。

## API 文档

### 认证

```bash
# 登录
POST /api/auth/login
{
  "username": "admin",
  "password": "admin123"
}
```

### 会话管理

```bash
# 列出所有会话
GET /api/sessions

# 创建会话
POST /api/sessions
{
  "name": "my-project",
  "work_dir": "/home/user/projects"
}

# 获取会话详情
GET /api/sessions/{name}

# 删除会话
DELETE /api/sessions/{name}

# 发送命令
POST /api/sessions/{name}/command
{
  "command": "ls -la"
}

# 获取输出
GET /api/sessions/{name}/output
```

### 文件操作

```bash
# 列出目录内容
GET /api/files?path=/home/user/projects

# 获取文件内容
GET /api/files/content?path=/home/user/projects/file.txt

# 创建文件/目录
POST /api/files
{
  "path": "/home/user/projects/newfile.txt",
  "type": "file",
  "content": "Hello World"
}

# 重命名
PUT /api/files/rename
{
  "oldPath": "/home/user/projects/old.txt",
  "newPath": "/home/user/projects/new.txt"
}

# 删除
DELETE /api/files?path=/home/user/projects/file.txt
```

### WebSocket

```javascript
// 连接到会话
const ws = new WebSocket('ws://localhost:9090/ws/session-name?token=YOUR_JWT_TOKEN')

// 发送命令
ws.send(JSON.stringify({
  type: 'command',
  data: 'your-command'
}))

// 发送按键（实时模式）
ws.send(JSON.stringify({
  type: 'keys',
  data: 'ls'
}))

// 接收输出
ws.onmessage = (event) => {
  const message = JSON.parse(event.data)
  if (message.type === 'output') {
    console.log(message.data.text)
  }
}
```

## 安全建议

1. **修改默认密码**: 首次使用前请修改管理员密码（`backend/.env` 中的 `ADMIN_PASSWORD`）
2. **使用强 JWT Secret**: 生成随机字符串作为 `JWT_SECRET`
3. **配置 HTTPS**: 生产环境建议使用 HTTPS
4. **限制访问**: 使用防火墙或 VPN 限制访问
5. **定期更新**: 保持依赖包最新

## 目录结构

```
remote-claude-code/
├── backend/              # Go 后端
│   ├── cmd/
│   │   └── server/       # 入口文件
│   ├── internal/         # 内部包
│   │   ├── api/          # API 处理器
│   │   ├── auth/         # 认证
│   │   ├── config/       # 配置
│   │   ├── security/     # 安全
│   │   ├── tmux/         # tmux 管理
│   │   └── websocket/    # WebSocket
│   └── Dockerfile
├── frontend/             # Vue 前端
│   ├── src/
│   │   ├── api/          # API 客户端
│   │   ├── components/   # 组件
│   │   ├── composables/  # Composables
│   │   ├── router/       # 路由
│   │   ├── stores/       # 状态管理
│   │   └── views/        # 页面
│   ├── Dockerfile
│   └── nginx.conf
├── frp/                  # FRP 内网穿透配置
├── docker-compose.yml
└── README.md
```

## 故障排查

### tmux 不可用
```bash
# macOS
brew install tmux

# Ubuntu/Debian
sudo apt-get install tmux

# CentOS/RHEL
sudo yum install tmux
```

### WebSocket 连接失败
- 检查 JWT Token 是否有效
- 确认会话名称正确
- 查看浏览器控制台错误信息
- 检查端口是否被占用或防火墙设置

### 会话输出为空
- 确认会话已启动 claude-code
- 检查 tmux 会话状态: `tmux list-sessions`

### 通过域名访问时登录失败
- 确保 Vite 配置中 `allowedHosts` 包含你的域名
- 检查后端 API 是否正常响应

## License

Apache License 2.0

## 贡献

欢迎提交 Issue 和 Pull Request！
