package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Get user home directory for config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	app.ConfigDir = filepath.Join(homeDir, ".remote-claude-code")
	app.ConfigPath = filepath.Join(app.ConfigDir, "config.json")

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Remote Claude Code",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Println("Error:", err.Error())
	}
}
