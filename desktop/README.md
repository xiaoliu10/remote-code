# Remote Claude Code Desktop Application

桌面版 Remote Claude Code 应用程序，支持 macOS、Windows 和 Linux。

## 功能特性

- 🖥️ **桌面原生应用** - 原生桌面体验，无需浏览器
- ⚙️ **首次配置向导** - 图形化配置界面，轻松设置
- 🔄 **自动启动服务** - 自动管理后端和 FRP 服务
- 📊 **状态监控** - 实时显示服务运行状态
- 🔧 **设置管理** - 随时修改配置

## 平台支持

| 平台 | 本地会话 | 远程连接 |
|------|----------|----------|
| macOS | ✅ 支持 | ✅ 支持 |
| Linux | ✅ 支持 | ✅ 支持 |
| Windows | ❌ 不支持* | ✅ 支持 |

*Windows 没有 tmux，因此不支持本地会话管理，但可以作为远程客户端连接到 Mac/Linux 服务器。

## 系统要求

### macOS
- macOS 10.15 (Catalina) 或更高版本
- Xcode Command Line Tools

### Windows
- Windows 10 或更高版本
- WebView2 Runtime (Windows 10+ 通常已内置)

### Linux
- GTK3
- WebKit2GTK
- gcc, pkg-config

## 开发环境设置

### 1. 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. 安装依赖

```bash
cd desktop
make deps
```

### 3. 运行开发模式

```bash
make dev
```

## 构建

### 当前平台

```bash
make build
```

### 所有平台

```bash
make build-all
```

### 特定平台

```bash
# macOS
make build-mac

# Windows (需要在 Windows 或交叉编译环境)
make build-windows

# Linux
make build-linux
```

## 项目结构

```
desktop/
├── main.go          # Wails 入口点
├── app.go           # 应用逻辑
├── wails.json       # Wails 配置
├── Makefile         # 构建脚本
├── go.mod           # Go 模块
├── frontend/        # 前端代码 (Vue)
│   ├── src/
│   │   ├── main.js
│   │   └── App.vue
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
└── build/           # 构建输出
    └── bin/
```

## 使用流程

1. **首次启动**
   - 显示配置向导
   - 设置管理员密码
   - 选择工作目录 (Mac/Linux)
   - 配置 FRP (可选)

2. **正常使用**
   - 自动启动后端服务
   - 显示服务状态
   - 点击 "Open in Browser" 访问完整界面

3. **修改配置**
   - 点击设置图标
   - 修改配置后保存

## 配置存储

配置文件存储位置：

- **macOS**: `~/Library/Application Support/remote-claude-code/config.json`
- **Windows**: `%APPDATA%/remote-claude-code/config.json`
- **Linux**: `~/.config/remote-claude-code/config.json`

## 打包分发

### macOS

```bash
# 生成 .app 文件
make build

# 创建 DMG (需要 create-dmg)
brew install create-dmg
make package-mac
```

### Windows

```bash
# 生成 .exe 文件
make build

# 使用 Inno Setup 或 NSIS 创建安装程序
```

### Linux

```bash
# 生成二进制文件
make build

# 创建 .deb, .rpm 或 AppImage
```

## 开发说明

### 前端开发

前端使用 Vue 3 + Vite，位于 `frontend/` 目录。

```bash
cd frontend
npm run dev    # 开发服务器
npm run build  # 构建生产版本
```

### 后端集成

后端 Go 服务需要单独编译并放置在 `build/bin/` 目录。

```bash
# 从项目根目录
cd backend
go build -o ../desktop/build/bin/backend ./cmd/server
```

### API 调用

前端通过 Wails 绑定调用 Go 函数：

```javascript
// 调用 Go 函数
const config = await window.go.main.App.GetConfig()
await window.go.main.App.SaveConfiguration(config)
```

## 故障排查

### Wails 命令未找到

确保 Go bin 目录在 PATH 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 前端构建失败

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### macOS 权限问题

```bash
xcode-select --install
```

## License

Apache License 2.0
