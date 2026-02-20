# 中国用户内网穿透方案指南

## 🇨🇳 推荐方案（适合中国用户）

### 方案一：Frp + 国内云服务器（⭐️ 最推荐）

#### 推荐的国内 VPS
- **阿里云 ECS** - 新用户有优惠
- **腾讯云 CVM** - 价格便宜
- **华为云** - 稳定性好
- **UCloud** - 便宜好用

#### 快速配置

**1. 在云服务器上安装 Frp Server**

```bash
# 下载 Frp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 配置 frps.toml
cat > frps.toml << EOF
bindPort = 7000
vhostHTTPPort = 80
auth.token = "your-secret-token"
EOF

# 启动
nohup ./frps -c frps.toml > frps.log 2>&1 &
```

**2. 在本地配置 Frp Client**

编辑 `frp/frpc.ini`:

```ini
serverAddr = "你的服务器公网IP"
serverPort = 7000
auth.token = "your-secret-token"

[[proxies]]
name = "remote-claude"
type = "tcp"
localIP = "127.0.0.1"
localPort = 8080
remotePort = 8080
```

**3. 启动服务**

```bash
./start-with-tunnel.sh frp
```

**4. 访问**

在任何地方访问：`http://你的服务器IP:8080`

---

### 方案二：使用国内内网穿透服务（最简单）

#### 1. **cpolar**（推荐）

- 官网：https://www.cpolar.com/
- 免费 1 个隧道
- 国内服务器，速度快
- 无需自己有服务器

**使用步骤**：

```bash
# 安装 cpolar
curl -sL https://git.io/cpolar | sudo bash

# 启动（会自动打开浏览器登录）
cpolar http 8080

# 会得到一个公网地址，例如：
# https://r4e8k2.cpolar.cn
```

#### 2. **花生壳**

- 官网：https://hsk.oray.com/
- 老牌稳定
- 免费版够用

#### 3. **NATAPP**

- 官网：https://natapp.cn/
- 免费隧道
- 支持自定义域名

#### 4. **Sakura Frp（免费）**

- 官网：https://www.natfrp.com/
- 免费 5 个隧道
- 国内节点多

---

### 方案三：ZeroTier（替代 Tailscale）

ZeroTier 可能在某些应用商店可用，且配置类似。

```bash
# 安装
brew install zerotier-one

# 启动服务
sudo brew services start zerotier-one

# 加入网络（在 https://my.zerotier.com 创建网络）
sudo zerotier-cli join your-network-id

# 获取 IP
sudo zerotier-cli listnetworks
```

---

## 📱 移动端访问方案

### 方案 A：使用 Web 浏览器（最简单）

无论什么方案，都可以直接用手机浏览器访问：
- Safari
- Chrome
- Edge

只需要访问分配的公网地址即可。

### 方案 B：使用 SSH 客户端

iOS 上可以使用：
- **Termius**（免费）
- **Blink Shell**
- **Prompt 2**

### 方案 C：使用 Tailscale 的 Web 管理界面

虽然 iOS 无法安装 Tailscale app，但可以：
1. 在电脑上使用 Tailscale
2. 在手机浏览器访问 Tailscale 管理页面
3. 通过其他方式访问服务

---

## 🎯 推荐组合

### 组合一：个人长期使用
- **方案**：Frp + 腾讯云/阿里云（$5/月）
- **优点**：稳定、快速、完全控制
- **适合**：经常需要远程访问

### 组合二：临时测试
- **方案**：cpolar（免费）
- **优点**：无需服务器，即开即用
- **适合**：偶尔使用

### 组合三：多设备访问
- **方案**：Frp + 国内 VPS
- **优点**：手机、平板、电脑都能访问
- **适合**：多设备用户

---

## 🚀 快速开始（推荐 cpolar）

最简单的方式，无需服务器：

```bash
# 1. 安装 cpolar
curl -sL https://git.io/cpolar | sudo bash

# 2. 启动 Remote Code 本地服务
./start-with-tunnel.sh local

# 3. 在另一个终端启动 cpolar
cpolar http 8080

# 4. 访问显示的公网地址（例如 https://r4e8k2.cpolar.cn）
```

---

## 📊 方案对比（中国用户）

| 方案 | 成本 | 速度 | 稳定性 | 配置难度 | 推荐度 |
|------|------|------|--------|----------|--------|
| Frp+国内VPS | ¥30-50/月 | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️⭐️ |
| cpolar | 免费/付费 | ⭐️⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️ | ⭐️ | ⭐️⭐️⭐️⭐️⭐️ |
| 花生壳 | 免费/付费 | ⭐️⭐️⭐️ | ⭐️⭐️⭐️⭐️ | ⭐️ | ⭐️⭐️⭐️⭐️ |
| NATAPP | 免费/付费 | ⭐️⭐️⭐️ | ⭐️⭐️⭐️ | ⭐️ | ⭐️⭐️⭐️ |
| ZeroTier | 免费 | ⭐️⭐️⭐️ | ⭐️⭐️⭐️ | ⭐️⭐️ | ⭐️⭐️⭐️ |

---

## 💡 小贴士

1. **安全第一**：无论用哪个方案，都要：
   - 修改默认密码（admin123）
   - 使用强 JWT_SECRET
   - 考虑添加 HTTPS

2. **速度优化**：
   - 优先选择国内服务器
   - 选择离你最近的节点
   - Frp 可以优化配置提高速度

3. **成本控制**：
   - 测试阶段用免费方案
   - 长期使用建议买个便宜 VPS
   - 阿里云/腾讯云新用户有很大优惠

---

## 🆘 需要帮助？

- **Frp 配置问题**：查看 [frp/README.md](../frp/README.md)
- **cpolar 使用**：访问 https://www.cpolar.com/docs
- **其他问题**：在项目提 Issue

---

**推荐**：如果是第一次使用，建议先用 **cpolar**（免费、简单），熟悉后再考虑购买 VPS 使用 Frp。
