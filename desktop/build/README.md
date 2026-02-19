# Build Directory

This directory contains build outputs and resources.

## Structure

```
build/
├── bin/              # Compiled executables
│   ├── remote-claude-code      # Main app (macOS/Linux)
│   ├── remote-claude-code.exe  # Main app (Windows)
│   └── backend                  # Backend server (embedded)
├── appicon.png       # Source icon for generating platform icons
├── icons/            # Generated platform-specific icons
│   ├── icon.icns    # macOS icon
│   ├── icon.ico     # Windows icon
│   └── png/         # Linux icons
└── darwin/           # macOS app bundle
    └── Remote Claude Code.app
```

## Icon Generation

1. Place a 512x512 PNG icon at `build/appicon.png`
2. Run: `make icons` or `wails generate icon`

## Embedding Backend

To embed the backend server:

```bash
# Build backend
cd ../../backend
go build -o ../desktop/build/bin/backend ./cmd/server

# Build desktop app with embedded backend
cd ../desktop
make build
```
