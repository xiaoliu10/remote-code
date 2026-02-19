# 端口修改说明

## 已修改端口：8080 → 9090

为避免与 Nacos 冲突，已将应用端口从 8080 改为 9090。

### 📝 已修改的文件

#### 后端
- ✅ `backend/.env` - PORT=9090
- ✅ `frontend/vite.config.ts` - proxy target 改为 http://localhost:9090
- ✅ `frontend/src/composables/useWebSocket.ts` - WS_BASE_URL 改为 ws://localhost:9090/api

#### Frp 配置
- ✅ `frp/frpc.ini` - localPort 和 remotePort 改为 9090
- ✅ `frp/install-server.sh` - 防火墙和提示信息改为 9090
- ✅ `frp/ALIYUN_SETUP.md` - 所有端口引用改为 9090

#### 文档
- ✅ `README.md` - 端口引用更新

---

## 📋 需要你更新的配置

### 1. 阿里云安全组

开放端口 **9090** 而不是 8080：

| 端口 | 说明 |
|------|------|
| 7000 | Frp 通信 |
| **9090** | 应用访问（新端口） |
| 7500 | 仪表板（可选） |

### 2. 服务器上（如果已经安装）

在服务器上执行：
```bash
# 如果使用 firewalld
firewall-cmd --permanent --add-port=9090/tcp
firewall-cmd --reload
```

### 3. 访问地址

**本地访问**：
- 前端：http://localhost:5173
- 后端：http://localhost:9090

**远程访问**（通过 Frp）：
- 应用：`http://YOUR_SERVER_IP:9090`
- 仪表板：`http://YOUR_SERVER_IP:7500`

---

## 🚀 启动服务

### 本地测试
```bash
./start-with-tunnel.sh local
```
访问：http://localhost:5173

### 使用 Frp 远程访问
```bash
./start-with-tunnel.sh frp
```
访问：http://YOUR_SERVER_IP:9090

---

## ✅ 验证端口

启动后验证：
```bash
# 检查本地端口
lsof -i :9090

# 或
netstat -tlnp | grep 9090
```

---

## 🔧 其他可能需要修改的地方

如果你的环境中有其他配置文件引用了 8080，请检查并修改为 9090。

---

**总结**：所有配置已更新为 9090 端口，请确保阿里云安全组也开放了 9090 端口。
