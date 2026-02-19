# 内网穿透方案选择指南

## 快速选择

| 场景 | 推荐方案 | 理由 |
|------|---------|------|
| 个人临时使用 | **Tailscale** | 5分钟搞定，无需服务器 |
| 长期稳定使用 | **Frp + VPS** | 完全控制，速度快 |
| 有域名 + 生产环境 | **Cloudflare Tunnel** | 免费，专业，自带 HTTPS |
| 临时测试/演示 | **免费 Frp 服务** | 无需任何资源 |

## 详细对比

### 1. Tailscale ⭐️ 推荐新手

**优点**：
- ✅ 零配置，安装即用
- ✅ P2P 直连，速度快
- ✅ 免费且无流量限制
- ✅ 支持所有平台（包括手机）

**缺点**：
- ❌ 需要在访问端也安装 Tailscale
- ❌ 某些企业网络可能被屏蔽

**成本**：免费

**适合**：个人使用，快速搭建

**使用步骤**：
```bash
# 1. 安装
brew install tailscale  # macOS
# 或访问 https://tailscale.com/download

# 2. 登录（需要 Google/GitHub/Microsoft 账号）
sudo tailscale up

# 3. 启动服务
./start-with-tunnel.sh tailscale

# 4. 在其他设备上同样安装登录，即可访问
```

---

### 2. Frp + 自己的服务器 ⭐️ 推荐长期使用

**优点**：
- ✅ 完全控制
- ✅ 速度快，稳定
- ✅ 支持多端口、多服务
- ✅ 可自定义域名

**缺点**：
- ❌ 需要购买 VPS（$3-10/月）
- ❌ 需要一定配置

**成本**：$3-10/月（VPS）

**适合**：长期稳定使用，多服务

**推荐 VPS**：
- [Vultr](https://www.vultr.com/) - $2.5/月起
- [DigitalOcean](https://www.digitalocean.com/) - $4/月
- [搬瓦工](https://bandwagonhost.com/) - $49.99/年

**使用步骤**：
```bash
# 服务器端（VPS）
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64
./frps -c frps.ini &

# 本地端
# 1. 编辑 frp/frpc.ini
# 2. 启动
./start-with-tunnel.sh frp
```

---

### 3. Cloudflare Tunnel ⭐️ 推荐生产环境

**优点**：
- ✅ 免费
- ✅ 自带 HTTPS
- ✅ 支持 DDoS 防护
- ✅ 无需开放端口
- ✅ 稳定可靠

**缺点**：
- ❌ 需要域名
- ❌ 需要配置 DNS

**成本**：免费（Cloudflare 免费套餐）

**适合**：有域名，生产环境

**使用步骤**：
```bash
# 1. 在 Cloudflare Dashboard 创建 Tunnel
# https://one.dash.cloudflare.com/

# 2. 获取 Token

# 3. 启动
export CLOUDFLARE_TUNNEL_TOKEN=your-token
./start-with-tunnel.sh cloudflare
```

---

### 4. 免费 Frp 服务

**优点**：
- ✅ 完全免费
- ✅ 无需服务器

**缺点**：
- ❌ 可能不稳定
- ❌ 有流量限制
- ❌ 速度可能较慢
- ❌ 安全性未知

**成本**：免费

**适合**：临时测试

**推荐服务**：
- [OpenFrp](https://openfrp.net/) - 免费 2 个隧道
- [Sakura Frp](https://www.natfrp.com/) - 免费 5 个隧道
- [Chmlfrp](https://www.chmlfrp.cn/) - 免费使用

---

### 5. ZeroTier（类似 Tailscale）

**优点**：
- ✅ 免费最多 25 设备
- ✅ 支持 P2P
- ✅ 支持所有平台

**缺点**：
- ❌ 配置比 Tailscale 复杂
- ❌ 连接速度可能较慢

**成本**：免费（25 设备内）

---

## 实测对比

| 指标 | Tailscale | Frp+VPS | Cloudflare Tunnel | 免费 Frp |
|------|-----------|---------|-------------------|----------|
| 延迟 | 20-50ms | 30-100ms | 50-200ms | 100-500ms |
| 稳定性 | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️ |
| 速度 | 快 | 很快 | 中等 | 慢 |
| 配置难度 | 简单 | 中等 | 中等 | 简单 |
| 推荐度 | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️ |

## 我的推荐

### 🎯 快速开始：Tailscale
```bash
# 1 分钟开始使用
brew install tailscale
sudo tailscale up
./start-with-tunnel.sh tailscale
```

### 🏆 长期使用：Frp + VPS
```bash
# 买一个便宜 VPS，配置一次，稳定使用
./start-with-tunnel.sh frp
```

### 🌐 有域名：Cloudflare Tunnel
```bash
# 免费 + 专业 + 安全
export CLOUDFLARE_TUNNEL_TOKEN=xxx
./start-with-tunnel.sh cloudflare
```

## 常见问题

### Q: 哪个最快？
A: Tailscale（P2P 直连）或 Frp + 本地 VPS

### Q: 哪个最稳定？
A: Frp + VPS 或 Cloudflare Tunnel

### Q: 哪个完全免费？
A: Tailscale、Cloudflare Tunnel、免费 Frp 服务

### Q: 哪个最简单？
A: Tailscale - 安装即用，无需配置

### Q: 需要公网 IP 吗？
A: 都不需要！这正是内网穿透的意义

### Q: 安全吗？
A:
- Tailscale: 端到端加密，非常安全
- Frp: 需要自己配置 TLS
- Cloudflare Tunnel: 企业级安全
- 建议：无论用哪个，都要设置强密码和 JWT Secret

## 下一步

1. 选择一个方案
2. 按照上面的步骤配置
3. 运行 `./start-with-tunnel.sh <方案>`
4. 开始使用！

有问题？查看 [frp/README.md](../frp/README.md) 获取详细配置说明。
