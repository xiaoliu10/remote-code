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
- 🎤 **语音输入**: 支持语音输入命令（需要 HTTPS）
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

### 一键启动

```bash
# 1. 克隆项目
git clone https://github.com/xiaoliu10/remote-claude-code.git
cd remote-claude-code

# 2. 启动服务（首次启动会自动创建配置）
./start.sh

# 访问 http://localhost:5173
```

首次启动时，系统会自动：
- 在 `~/.remote-claude-code/` 创建配置目录
- 生成随机管理员密码（请保存！）
- 启动后端服务（端口 9090）
- 启动前端服务（端口 5173）

### 管理命令

```bash
./start.sh          # 启动服务
./start.sh --frp    # 启动服务并启用 FRP
./start.sh --no-frp # 启动服务（禁用 FRP）
./stop.sh           # 停止服务
```

### 配置文件

配置文件位于 `~/.remote-claude-code/config.ini`：

```ini
# 后端配置
BACKEND_PORT=9090
JWT_SECRET=auto-generated
ADMIN_PASSWORD=auto-generated
ALLOWED_DIR=/home/user/projects

# FRP 配置
FRP_ENABLED=false
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-token
```

## 远程访问

### 为什么需要内网穿透？

Remote Claude Code 默认运行在本地网络中，只能在局域网内访问。当你需要**在外网（如公司、咖啡厅、移动网络）远程访问家中的 Claude Code 会话**时，就需要内网穿透。

**适用场景：**
- 家中电脑运行 Claude Code，外出时需要远程控制
- 公司网络无法暴露端口到公网
- 没有公网 IP 或运营商封锁了端口
- 需要稳定的远程访问体验

### 方式一：FRP + Nginx HTTPS（推荐）

FRP (Fast Reverse Proxy) 是一个高性能的内网穿透工具，通过公网服务器将本地服务暴露到外网。

#### 1. 配置 FRP

编辑 `~/.remote-claude-code/config.ini`：

```ini
FRP_ENABLED=true
FRP_SERVER_ADDR=your-server-ip
FRP_SERVER_PORT=7000
FRP_TOKEN=your-token
```

#### 2. 服务器端配置 Nginx

在服务器上创建 Nginx 配置：

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

    # 后端 API
    location /api {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Host $host;
    }

    # 前端
    location / {
        proxy_pass http://127.0.0.1:5173;
        proxy_set_header Host $host;
    }
}
```

#### 3. 启动服务

```bash
./start.sh --frp
```

#### 4. 访问

```
https://your-domain.com:8444
```

### 方式二：Tailscale

```bash
# 安装并登录 Tailscale
brew install tailscale
sudo tailscale up

# 启动服务
./start.sh

# 访问 http://<tailscale-ip>:5173
```

### 方式三：Cloudflare Tunnel

```bash
export CLOUDFLARE_TUNNEL_TOKEN=your-token
docker-compose --profile cloudflare up -d
```

## 使用说明

### 1. 登录系统

首次启动时会生成随机密码，查看启动输出获取密码。

### 2. 创建会话

1. 输入会话名称（如: `my-project`）
2. 选择工作目录（可选）
3. 点击 "Create" 创建会话

### 3. 终端使用

终端支持两种输入模式：
- **实时模式**: 每个字符实时发送，支持自动补全
- **命令模式**: 按 Enter 发送完整命令

移动端虚拟键盘：
- 方向键（上下左右）
- Tab、Ctrl+C、Ctrl+D
- @ 符号（文件引用）
- 🎤 语音输入

### 4. 语音输入

⚠️ **语音输入需要 HTTPS**

语音输入功能由于浏览器安全限制，需要满足以下条件之一：

| 方式 | 说明 |
|------|------|
| HTTPS | 使用 Nginx + Let's Encrypt |
| localhost | 本地访问无需配置 |
| Chrome 标志 | 仅测试用，访问 `chrome://flags/#unsafely-treat-insecure-origin-as-secure` |

### 5. 文件浏览

点击侧边栏的文件图标，浏览工作目录中的文件。

## 目录结构

```
~/.remote-claude-code/          # 配置目录
├── config.ini                  # 主配置文件
├── frpc.ini                    # FRP 配置（自动生成）
├── frpc                        # FRP 客户端（自动下载）
└── logs/                       # 日志目录
    ├── backend.log
    ├── frontend.log
    └── frp.log

remote-claude-code/             # 源码目录
├── backend/                    # Go 后端
├── frontend/                   # Vue 前端
├── desktop/                    # 桌面应用
├── nginx/                      # Nginx 配置
├── docs/                       # 文档
├── start.sh                    # 启动脚本
└── stop.sh                     # 停止脚本
```

## API 文档

### 认证

```bash
POST /api/auth/login
{"username": "admin", "password": "your-password"}
```

### 会话管理

```bash
GET    /api/sessions              # 列出会话
POST   /api/sessions              # 创建会话
GET    /api/sessions/{name}       # 获取详情
DELETE /api/sessions/{name}       # 删除会话
POST   /api/sessions/{name}/command  # 发送命令
```

### 文件操作

```bash
GET    /api/files?path=/path     # 列出目录
GET    /api/files/content?path=/path  # 获取内容
POST   /api/files                # 创建文件/目录
PUT    /api/files/rename         # 重命名
DELETE /api/files?path=/path     # 删除
```

### WebSocket

```javascript
const ws = new WebSocket('wss://your-domain:8444/api/ws/session?token=TOKEN')

// 发送命令
ws.send(JSON.stringify({type: 'command', data: 'ls -la'}))

// 发送按键
ws.send(JSON.stringify({type: 'keys', data: 'ls'}))
```

## 故障排查

### 端口被占用

```bash
./stop.sh  # 停止所有服务（包括孤儿进程）
```

### WebSocket 连接失败

- 检查 Nginx 配置中 `/api/ws` 是否正确
- 确认 HTTPS 证书有效
- 查看浏览器控制台错误

### 语音功能不可用

- 确保使用 HTTPS 访问
- 或使用 localhost 测试
- 检查浏览器是否支持 Web Speech API

## 安全建议

1. **保存初始密码**: 首次启动生成的随机密码请妥善保存
2. **使用 HTTPS**: 生产环境务必配置 HTTPS
3. **定期更新**: 保持依赖包最新
4. **限制访问**: 使用防火墙限制访问范围

## License

Apache License 2.0

## 贡献

欢迎提交 Issue 和 Pull Request！
