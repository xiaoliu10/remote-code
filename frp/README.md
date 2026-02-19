# 使用 Frp 进行内网穿透

## 方案一：使用自己的服务器（推荐）

### 1. 在有公网 IP 的服务器上安装 Frp Server

```bash
# 下载 frp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 配置 frps.toml
cat > frps.toml << EOF
bindPort = 7000
vhostHTTPPort = 80

# 认证令牌（重要！请修改为你自己的复杂密码）
auth.token = "your-secure-token-here"

# 仪表板配置（可选）
webServer.addr = "0.0.0.0"
webServer.port = 7500
webServer.user = "admin"
webServer.password = "your-dashboard-password"
EOF

# 启动
./frps -c frps.toml
```

### 2. 在本地机器上配置 Frp Client

编辑 `frp/frpc.ini`，填入你的服务器信息：

```ini
serverAddr = "your-server-ip"
serverPort = 7000

# 认证令牌（必须与服务器端一致）
auth.token = "your-secure-token-here"

# 前端服务 (Vite 开发服务器)
[[proxies]]
name = "remote-claude-frontend"
type = "tcp"
localIP = "127.0.0.1"
localPort = 5173
remotePort = 5173

# 后端 API + WebSocket 服务
[[proxies]]
name = "remote-claude-backend"
type = "tcp"
localIP = "127.0.0.1"
localPort = 9090
remotePort = 9090
```

### 3. 启动服务

```bash
# 方式 1: 使用 Docker Compose（包含 frp）
docker-compose --profile frp up -d

# 方式 2: 手动启动
# 终端 1: 启动后端
cd backend && go run cmd/server/main.go

# 终端 2: 启动前端
cd frontend && npm run dev

# 终端 3: 启动 frp 客户端
./frpc -c frp/frpc.ini
```

### 4. 访问

- 前端：`http://your-server-ip:5173`
- 后端：`http://your-server-ip:9090`
- 仪表板：`http://your-server-ip:7500`

### 5. 服务器防火墙配置

确保开放以下端口：
```bash
# 使用 firewalld
firewall-cmd --permanent --add-port=7000/tcp  # Frp 通信
firewall-cmd --permanent --add-port=5173/tcp  # 前端
firewall-cmd --permanent --add-port=9090/tcp  # 后端
firewall-cmd --permanent --add-port=7500/tcp  # 仪表板（可选）
firewall-cmd --reload
```

## 方案二：使用免费的 Frp 服务

### 1. 使用 Sakura Frp（免费）

1. 访问 https://openfrp.net/
2. 注册账号
3. 创建隧道，选择 TCP 隧道
4. 本地 IP: `127.0.0.1`
5. 本地端口: `9090`
6. 下载启动器并运行

### 2. 使用 Cloudflare Tunnel（免费）

```bash
# 安装 cloudflared
brew install cloudflare/cloudflare/cloudflared

# 登录
cloudflared tunnel login

# 创建隧道
cloudflared tunnel create remote-claude

# 运行隧道
cloudflared tunnel route dns remote-claude remote-claude.your-domain.com
cloudflared tunnel run --url http://localhost:9090 remote-claude
```

## 方案三：使用 Tailscale（最简单）

### 1. 安装 Tailscale

```bash
# macOS
brew install tailscale

# Linux
curl -fsSL https://tailscale.com/install.sh | sh
```

### 2. 在两台机器上登录

```bash
# 家里的电脑
sudo tailscale up

# 外出的电脑/手机
sudo tailscale up
```

### 3. 直接访问

使用 Tailscale 分配的内网 IP 访问，例如：
`http://100.x.y.z:5173`

## 方案四：使用 ZeroTier（类似 Tailscale）

```bash
# 安装
brew install zerotier-one

# 加入网络（在 https://my.zerotier.com 创建）
sudo zerotier-cli join your-network-id

# 授权后在任何地方访问
```

## 推荐方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| Frp + 自己服务器 | 完全控制、稳定、快速 | 需要服务器 | 长期使用 |
| 免费Frp服务 | 免费、无需服务器 | 可能不稳定、有流量限制 | 临时测试 |
| Cloudflare Tunnel | 免费、稳定、自带HTTPS | 需要域名 | 生产环境 |
| Tailscale | 超级简单、P2P连接 | 速度取决于网络 | 个人使用 |
| ZeroTier | 类似 Tailscale | 配置稍复杂 | 个人使用 |

## 快速开始

对于本项目，推荐：

1. **测试阶段**: 使用 Tailscale（5 分钟搞定）
2. **正式使用**: Frp + 便宜 VPS（$3-5/月）
3. **有域名**: Cloudflare Tunnel（免费且专业）

## 端口说明

| 端口 | 用途 |
|------|------|
| 7000 | Frp 服务器通信端口 |
| 5173 | 前端服务（开发模式） |
| 9090 | 后端 API + WebSocket |
| 7500 | Frp 仪表板（可选） |
