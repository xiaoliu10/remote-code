# 常见问题解答 (FAQ)

## Claude Code 相关问题

### 1. 在 Remote Code 中运行 Claude Code 报错：Nested sessions are not allowed

**错误信息：**
```
Error: Claude Code cannot be launched inside another Claude Code session.
Nested sessions share runtime resources and will crash all active sessions.
To bypass this check, unset the CLAUDECODE environment variable.
```

**原因：**
Claude Code 默认禁止在另一个 Claude Code 会话中运行，以避免资源冲突。

**解决方案：**

取消设置 `CLAUDECODE` 环境变量即可绕过此限制。

**临时设置（仅当前会话有效）：**
```bash
unset CLAUDECODE
claude
```

**永久设置（推荐）：**

在 `~/.zshrc`（zsh）或 `~/.bashrc`（bash）中添加：
```bash
# Allow Claude Code to run in nested sessions (e.g., inside Remote Code)
unset CLAUDECODE
```

然后执行以下命令使其生效：
```bash
# zsh
source ~/.zshrc

# bash
source ~/.bashrc
```

**对于已存在的 tmux 会话：**

如果 tmux 会话是在添加配置之前创建的，需要在会话中重新加载配置：
```bash
source ~/.zshrc
```

或者重新创建 tmux 会话以继承新的环境变量。

---

## 语音输入相关问题

### 2. 语音输入功能不可用

**原因：**
浏览器安全策略要求语音识别 API 必须在安全上下文（HTTPS）中使用。

**解决方案：**

| 方式 | 说明 |
|------|------|
| HTTPS | 使用 Nginx + Let's Encrypt 配置 HTTPS |
| localhost | 本地访问无需配置，直接支持 |
| Chrome 标志 | 仅测试用，访问 `chrome://flags/#unsafely-treat-insecure-origin-as-secure` |

---

## 连接相关问题

### 3. WebSocket 连接失败

**排查步骤：**
1. 检查 Nginx 配置中 `/api/ws` 是否正确配置了 WebSocket 代理
2. 确认 HTTPS 证书有效
3. 查看浏览器控制台错误信息

**Nginx WebSocket 代理配置示例：**
```nginx
location /api/ws {
    proxy_pass http://127.0.0.1:9090;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_read_timeout 86400s;
}
```

### 4. 多设备连接冲突

**现象：**
同一个会话被另一个设备连接后，当前设备收到 "Connection replaced by another device" 提示。

**说明：**
这是正常的安全机制，同一会话同时只允许一个设备连接，防止终端状态混乱。

**解决方案：**
刷新页面重新连接即可。

---

## 其他问题

### 5. 端口被占用

**解决方案：**
```bash
./stop.sh  # 停止所有服务（包括孤儿进程）
```

### 6. 如何查看日志

```bash
# 查看所有日志
tail -f ~/.remote-code/logs/*.log

# 仅查看后端日志
tail -f ~/.remote-code/logs/backend.log

# 仅查看前端日志
tail -f ~/.remote-code/logs/frontend.log
```

---

## 更多帮助

如果以上内容没有解决你的问题，请提交 Issue：
https://github.com/xiaoliu10/remote-code/issues
