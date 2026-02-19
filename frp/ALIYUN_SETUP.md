# 使用阿里云服务器配置 Frp 内网穿透

## 📋 总览

- **服务器**：阿里云 ECS
- **系统**：CentOS/RHEL
- **方案**：Frp 内网穿透
- **需要开放的端口**：7000（Frp通信）、9090（应用）、7500（仪表板，可选）

---

## 第一步：配置阿里云安全组

在配置 Frp 之前，需要在阿里云控制台开放端口：

1. 登录 [阿里云控制台](https://ecs.console.aliyun.com/)
2. 进入 **云服务器 ECS** → 选择你的实例
3. 点击 **安全组** → **配置规则**
4. **添加入方向规则**：

| 端口范围 | 授权对象 | 说明 |
|---------|---------|------|
| 7000 | 0.0.0.0/0 | Frp 服务端口 |
| 9090 | 0.0.0.0/0 | Remote Claude Code 应用 |
| 7500 | 0.0.0.0/0 | Frp 仪表板（可选） |

---

## 第二步：在服务器上安装 Frp

### 方式 A：使用一键脚本（推荐）

1. **上传脚本到服务器**：

在本地执行：
```bash
scp frp/install-server.sh root@你的服务器IP:/root/
scp frp/frps.toml root@你的服务器IP:/opt/frp/
```

2. **SSH 连接到服务器**：
```bash
ssh root@你的服务器IP
```

3. **运行安装脚本**：
```bash
cd /root
chmod +x install-server.sh
./install-server.sh
```

### 方式 B：手动安装

```bash
# 1. 下载 Frp
cd /opt
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_linux_amd64.tar.gz
tar -xzf frp_0.52.3_linux_amd64.tar.gz
cd frp_0.52.3_linux_amd64

# 2. 安装
cp frps /opt/frp/
cp frps.toml /opt/frp/

# 3. 创建 systemd 服务
vim /etc/systemd/system/frps.service
```

粘贴以下内容：
```ini
[Unit]
Description=Frp Server Service
After=network.target

[Service]
Type=simple
ExecStart=/opt/frp/frps -c /opt/frp/frps.toml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

---

## 第三步：配置 Frp 服务器

1. **编辑配置文件**：
```bash
vim /opt/frp/frps.toml
```

2. **修改认证令牌（重要！）**：
```toml
auth.token = "你的复杂密码-请修改这里"
```

3. **可选：修改仪表板密码**：
```toml
webServer.password = "你的仪表板密码"
```

---

## 第四步：启动 Frp 服务器

```bash
# 启动服务
systemctl start frps

# 设置开机自启
systemctl enable frps

# 查看状态
systemctl status frps

# 查看日志
tail -f /var/log/frps.log
```

成功后你会看到类似：
```
[I] [service.go:XXX] frps started successfully
```

---

## 第五步：配置本地 Frp 客户端

在你的 Mac 上：

1. **编辑配置文件**：
```bash
vim frp/frpc.ini
```

2. **修改服务器地址**：
```ini
serverAddr = "你的服务器公网IP"  # 例如：123.45.67.89
auth.token = "change-this-to-your-secret-token-2024"  # 与服务器端一致
```

3. **下载 Frp 客户端（如果还没有）**：
```bash
cd frp
wget https://github.com/fatedier/frp/releases/download/v0.52.3/frp_0.52.3_darwin_arm64.tar.gz
tar -xzf frp_0.52.3_darwin_arm64.tar.gz
mv frp_0.52.3_darwin_arm64/frpc .
chmod +x frpc
```

4. **测试连接**：
```bash
./frpc -c frpc.ini
```

---

## 第六步：启动 Remote Claude Code

```bash
# 方式 1：使用启动脚本
./start-with-tunnel.sh frp

# 方式 2：手动启动
# 终端 1：启动 Frp 客户端
cd frp && ./frpc -c frpc.ini

# 终端 2：启动后端
cd backend && go run cmd/server/main.go

# 终端 3：启动前端
cd frontend && npm run dev
```

---

## 第七步：访问服务

从任何地方访问：
```
http://你的服务器IP:9090
```

例如：`http://123.45.67.89:9090`

---

## 🎯 快速检查清单

- [ ] 阿里云安全组已开放 7000、9090 端口
- [ ] 服务器上 Frp 服务正在运行（`systemctl status frps`）
- [ ] 本地 frpc.ini 中的 serverAddr 已修改
- [ ] 本地 frpc.ini 中的 auth.token 与服务器一致
- [ ] 本地 Frp 客户端连接成功
- [ ] Remote Claude Code 后端和前端已启动

---

## 🔧 故障排查

### 问题 1：连接失败

**检查**：
```bash
# 在服务器上查看 Frp 日志
tail -f /var/log/frps.log

# 检查端口是否监听
netstat -tlnp | grep 7000

# 测试端口连通性（在本地执行）
telnet 你的服务器IP 7000
```

### 问题 2：认证失败

**原因**：auth.token 不一致

**解决**：确保服务器和客户端的 auth.token 完全相同

### 问题 3：无法访问应用

**检查**：
```bash
# 确认本地服务正在运行
curl http://localhost:9090/health

# 检查 Frp 代理状态
curl http://服务器IP:7500
```

### 问题 4：防火墙问题

**阿里云安全组检查**：
- 确认规则已添加
- 授权对象是 0.0.0.0/0
- 方向是入方向

**服务器防火墙**：
```bash
# CentOS/RHEL
firewall-cmd --list-ports

# 如果没有输出，添加规则
firewall-cmd --permanent --add-port=7000/tcp
firewall-cmd --permanent --add-port=9090/tcp
firewall-cmd --reload
```

---

## 📊 性能优化

### 提高 Frp 性能

编辑 `/opt/frp/frps.toml`：

```toml
# 增加传输缓冲区大小
transport.tcpMux = true
transport.tcpMuxKeepaliveInterval = 60

# 最大连接池大小
transport.maxPoolCount = 5

# 心跳配置
transport.heartbeatTimeout = 90
```

---

## 🔒 安全建议

1. **使用强密码**：
   - auth.token 使用随机生成的复杂密码
   - 至少 32 位，包含大小写字母、数字、特殊字符

2. **限制仪表板访问**：
   ```toml
   webServer.addr = "127.0.0.1"  # 只允许本地访问
   ```

3. **使用 HTTPS**（推荐）：
   - 配置 Nginx 反向代理
   - 使用 Let's Encrypt 免费证书

4. **定期更新**：
   ```bash
   # 检查新版本
   https://github.com/fatedier/frp/releases
   ```

---

## 📱 移动端访问

配置完成后，在手机上：
1. 打开浏览器（Safari、Chrome 等）
2. 访问：`http://你的服务器IP:9090`
3. 登录使用

无需安装任何 App！

---

## 🆘 需要帮助？

1. 查看服务器日志：`tail -f /var/log/frps.log`
2. 查看客户端日志：运行 `./frpc -c frpc.ini` 时会显示
3. 检查网络连通性：`telnet 服务器IP 7000`
4. 查看项目文档：[frp/README.md](./frp/README.md)

---

**下一步**：现在开始配置你的服务器吧！从第一步配置阿里云安全组开始。
