# Nginx 反向代理配置指南

本指南帮助你在远程服务器上配置 Nginx 反向代理，实现 HTTPS 访问。

## 前提条件

- 远程服务器已安装 Nginx
- 已有 HTTPS 证书（Let's Encrypt 或其他）
- FRP 已配置并运行（本地 5173/9090 端口已转发到服务器）

## 配置步骤

### 1. 复制配置文件

```bash
# 在服务器上
sudo cp nginx/remote-claude-code.conf /etc/nginx/sites-available/
sudo ln -s /etc/nginx/sites-available/remote-claude-code.conf /etc/nginx/sites-enabled/
```

### 2. 修改配置

编辑配置文件，修改以下内容：

```bash
sudo vim /etc/nginx/sites-available/remote-claude-code.conf
```

需要修改：
- `server_name` - 改为你的域名
- `ssl_certificate` - 改为你的证书路径
- `ssl_certificate_key` - 改为你的私钥路径

### 3. 测试配置

```bash
sudo nginx -t
```

### 4. 重载 Nginx

```bash
sudo systemctl reload nginx
```

### 5. 防火墙配置

确保开放 80 和 443 端口：

```bash
# 使用 firewalld
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload

# 或使用 ufw
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

## 域名配置

### 方案一：使用子域名

1. **前端**: `claude.yourdomain.com` → 代理到 5173
2. **后端**: `api.yourdomain.com` → 代理到 9090

需要配置两个 DNS A 记录指向服务器 IP。

### 方案二：使用路径

只使用一个域名，通过路径区分：

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    # 前端
    location / {
        proxy_pass http://127.0.0.1:5173;
    }

    # 后端 API
    location /api {
        proxy_pass http://127.0.0.1:9090;
    }

    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:9090;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Let's Encrypt 证书

如果没有证书，可以使用 certbot 获取免费证书：

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d claude.yourdomain.com -d api.yourdomain.com

# 自动续期
sudo systemctl enable certbot.timer
```

## 验证

配置完成后，访问以下地址验证：

- 前端: `https://claude.yourdomain.com`
- 后端: `https://api.yourdomain.com/api/health`

## 故障排查

### 502 Bad Gateway
- 检查 FRP 是否正常运行
- 检查本地服务是否启动
- 检查端口是否正确

### SSL 证书错误
- 检查证书路径是否正确
- 检查证书是否过期

### WebSocket 连接失败
- 确保 Nginx 配置了 `Upgrade` 和 `Connection` 头
- 检查超时设置

## 完整架构

```
用户 → HTTPS (443) → Nginx → 本地回环 (5173/9090) → FRP Client → FRP Server → 本地服务
```

## 安全建议

1. 启用 HSTS
2. 配置 CSP 头
3. 限制请求大小
4. 启用访问日志
