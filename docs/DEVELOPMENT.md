# 开发文档

## 重要提醒

### ⚠️ 代码更新后必须重新编译二进制文件

**问题：**
修改后端Go代码后，`./build.sh` 可能不会立即更新二进制文件。

**解决方案：**

```bash
# 方法1：手动编译（推荐）
cd backend
go build -o ../build/remote-code-macos-apple ./cmd/server

# 方法2：使用build.sh（如果不起作用，使用方法1）
./build.sh
```

**验证编译：**
```bash
# 检查二进制文件时间戳
ls -lh build/remote-code-macos-apple

# 检查是否包含新功能
strings build/remote-code-macos-apple | grep "scroll_up"
```

**重启服务：**
```bash
./stop.sh
./start.sh --frp
```

## 当前实现状态

### Tmux Copy Mode 滚动功能（方案1）

**实现方式：**
- 使用 `tmux copy-mode` 命令进入copy mode
- 使用 `tmux send-keys -X cursor-up/down` 滚动
- 使用 `tmux send-keys -X cancel` 退出

**后端API：**
- `enter_copy_mode` - 进入copy mode
- `scroll_up` - 向上滚动
- `scroll_down` - 向下滚动
- `exit_copy_mode` - 退出copy mode

**前端使用：**
1. 虚拟键盘"滚动"按钮
2. 鼠标滚轮（远程模式）
3. Ctrl+B 快捷键

## 常见问题

### 1. 后端日志显示 "Unknown message type"

**原因：** 二进制文件未包含新代码

**解决：** 手动编译后端
```bash
cd backend
go build -o ../build/remote-code-macos-apple ./cmd/server
./stop.sh && ./start.sh --frp
```

### 2. scroll_up/scroll_down 不工作

**检查步骤：**
1. 查看后端日志：`tail -f ~/.remote-code/logs/backend.log`
2. 确认没有 "Unknown message type" 错误
3. 如果有，说明需要重新编译

### 3. 二进制文件时间戳未更新

**原因：** go build 可能因为缓存而不重新编译

**解决：**
```bash
# 清理并重新编译
cd backend
go clean
go build -o ../build/remote-code-macos-apple ./cmd/server
```
