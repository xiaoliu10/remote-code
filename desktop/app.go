package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// App application struct
type App struct {
	ctx        context.Context
	ConfigDir  string
	ConfigPath string
	config     *Config
	backendCmd *exec.Cmd
	isRunning  bool
	mu         sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		config: GetDefaultConfig(),
	}
}

// Config represents the application configuration
type Config struct {
	// Server settings
	BackendPort  int    `json:"backendPort"`
	FrontendPort int    `json:"frontendPort"`

	// Security
	JWTSecret     string `json:"jwtSecret"`
	AdminPassword string `json:"adminPassword"`

	// Working directory
	AllowedDir string `json:"allowedDir"`

	// FRP settings
	FRPEnabled   bool   `json:"frpEnabled"`
	FRPServer    string `json:"frpServer"`
	FRPPort      int    `json:"frpPort"`
	FRPToken     string `json:"frpToken"`

	// Setup completed flag
	SetupCompleted bool `json:"setupCompleted"`

	// Platform info (read-only)
	Platform string `json:"platform"`
	IsLocalSupported bool `json:"isLocalSupported"`
}

// GetDefaultConfig returns default configuration
func GetDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	platform := runtime.GOOS
	isLocalSupported := platform != "windows" // Windows doesn't have tmux

	return &Config{
		BackendPort:     9090,
		FrontendPort:    5173,
		JWTSecret:       generateRandomString(32),
		AdminPassword:   generateRandomString(12),
		AllowedDir:      filepath.Join(homeDir, "projects"),
		FRPEnabled:      false,
		FRPServer:       "",
		FRPPort:         7000,
		FRPToken:        "",
		SetupCompleted:  false,
		Platform:        platform,
		IsLocalSupported: isLocalSupported,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadConfig()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	a.StopServices()
}

// loadConfig loads configuration from file
func (a *App) loadConfig() {
	a.config = GetDefaultConfig()

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(a.ConfigDir, 0755); err != nil {
		logMessage("Failed to create config directory: " + err.Error())
		return
	}

	data, err := os.ReadFile(a.ConfigPath)
	if err != nil {
		// Config file doesn't exist, will use defaults
		return
	}

	json.Unmarshal(data, a.config)
}

// saveConfig saves configuration to file
func (a *App) saveConfig() error {
	dir := filepath.Dir(a.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(a.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(a.ConfigPath, data, 0644)
}

// GetConfig returns current configuration
func (a *App) GetConfig() *Config {
	return a.config
}

// SaveConfiguration saves the configuration and marks setup as complete
func (a *App) SaveConfiguration(config Config) error {
	a.config = &config
	a.config.SetupCompleted = true
	return a.saveConfig()
}

// IsSetupComplete returns whether setup has been completed
func (a *App) IsSetupComplete() bool {
	return a.config.SetupCompleted
}

// StartServices starts the backend and optionally FRP
func (a *App) StartServices() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isRunning {
		return nil
	}

	// Get the executable directory
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Check platform support
	if !a.config.IsLocalSupported {
		// On Windows, we can only act as a remote client
		return a.startRemoteMode()
	}

	// Start backend
	backendPath := findBackend(exeDir)
	if backendPath == "" {
		return fmt.Errorf("backend executable not found")
	}

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("PORT=%d", a.config.BackendPort))
	env = append(env, fmt.Sprintf("JWT_SECRET=%s", a.config.JWTSecret))
	env = append(env, fmt.Sprintf("ADMIN_PASSWORD=%s", a.config.AdminPassword))
	env = append(env, fmt.Sprintf("ALLOWED_DIR=%s", a.config.AllowedDir))

	a.backendCmd = exec.Command(backendPath)
	a.backendCmd.Env = env
	a.backendCmd.Stdout = os.Stdout
	a.backendCmd.Stderr = os.Stderr

	if err := a.backendCmd.Start(); err != nil {
		return fmt.Errorf("failed to start backend: %v", err)
	}

	// Start FRP if enabled
	if a.config.FRPEnabled && a.config.FRPServer != "" {
		if err := a.startFRP(); err != nil {
			logMessage("FRP start error: " + err.Error())
		}
	}

	a.isRunning = true

	// Wait for backend to be ready
	go a.waitForBackend()

	return nil
}

// startRemoteMode starts in remote-only mode (for Windows)
func (a *App) startRemoteMode() error {
	a.isRunning = true
	return nil
}

// waitForBackend waits for the backend to be ready
func (a *App) waitForBackend() {
	url := fmt.Sprintf("http://localhost:%d/api/health", a.config.BackendPort)
	for i := 0; i < 30; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				logMessage("Backend is ready")
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	logMessage("Backend failed to start within timeout")
}

// StopServices stops all running services
func (a *App) StopServices() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.backendCmd != nil && a.backendCmd.Process != nil {
		a.backendCmd.Process.Kill()
		a.backendCmd = nil
	}

	// Kill FRP process if running
	if runtime.GOOS != "windows" {
		exec.Command("pkill", "-f", "frpc.*remote-claude").Run()
	} else {
		exec.Command("taskkill", "/F", "/IM", "frpc.exe").Run()
	}

	a.isRunning = false
}

// GetServiceStatus returns the current status of services
func (a *App) GetServiceStatus() map[string]interface{} {
	status := map[string]interface{}{
		"isRunning":       a.isRunning,
		"backendRunning":  a.isBackendRunning(),
		"backendPort":     a.config.BackendPort,
		"frontendPort":    a.config.FrontendPort,
		"frpEnabled":      a.config.FRPEnabled,
		"platform":        a.config.Platform,
		"isLocalSupported": a.config.IsLocalSupported,
	}

	if a.config.FRPEnabled {
		status["frpServer"] = a.config.FRPServer
	}

	return status
}

// isBackendRunning checks if backend is responding
func (a *App) isBackendRunning() bool {
	url := fmt.Sprintf("http://localhost:%d/api/health", a.config.BackendPort)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// GetAppURL returns the URL to open in the webview
func (a *App) GetAppURL() string {
	return fmt.Sprintf("http://localhost:%d", a.config.BackendPort)
}

// OpenInBrowser opens a URL in the default browser
func (a *App) OpenInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	// This would normally use native dialogs
	// For simplicity, we'll return the current allowed dir
	return a.config.AllowedDir, nil
}

// startFRP starts the FRP client
func (a *App) startFRP() error {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	frpcPath := filepath.Join(exeDir, "frpc")
	if runtime.GOOS == "windows" {
		frpcPath += ".exe"
	}

	if _, err := os.Stat(frpcPath); os.IsNotExist(err) {
		return fmt.Errorf("FRP client not found at %s", frpcPath)
	}

	// Generate FRP config
	frpConfig := fmt.Sprintf(`
serverAddr = "%s"
serverPort = %d
auth.token = "%s"

[[proxies]]
name = "remote-claude-frontend"
type = "tcp"
localIP = "127.0.0.1"
localPort = %d
remotePort = %d

[[proxies]]
name = "remote-claude-backend"
type = "tcp"
localIP = "127.0.0.1"
localPort = %d
remotePort = %d
`, a.config.FRPServer, a.config.FRPPort, a.config.FRPToken,
		a.config.FrontendPort, a.config.FrontendPort,
		a.config.BackendPort, a.config.BackendPort)

	configPath := filepath.Join(os.TempDir(), "frpc-config.ini")
	if err := os.WriteFile(configPath, []byte(frpConfig), 0644); err != nil {
		return err
	}

	cmd := exec.Command(frpcPath, "-c", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}

// Helper functions

func findBackend(exeDir string) string {
	// Check common locations
	locations := []string{
		filepath.Join(exeDir, "backend"),
		filepath.Join(exeDir, "..", "backend", "server"),
		filepath.Join(exeDir, "resources", "backend"),
	}

	for _, loc := range locations {
		if runtime.GOOS == "windows" {
			loc += ".exe"
		}
		if _, err := os.Stat(loc); err == nil {
			return loc
		}
	}

	// Try to find in PATH
	path, err := exec.LookPath("remote-claude-backend")
	if err == nil {
		return path
	}

	return ""
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}

func logMessage(msg string) {
	fmt.Printf("[Desktop App] %s\n", msg)
}

// GetEmbeddedFrontend returns the embedded frontend URL for Wails
func (a *App) GetEmbeddedFrontend() string {
	// In production, the frontend is embedded
	// Return empty to use embedded assets
	return ""
}

// ReadFile reads a file and returns its contents
func (a *App) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file
func (a *App) WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// FileExists checks if a file exists
func (a *App) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListDirectory lists contents of a directory
func (a *App) ListDirectory(path string) ([]map[string]interface{}, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"name":  entry.Name(),
			"isDir": entry.IsDir(),
			"size":  info.Size(),
			"mode":  info.Mode().String(),
		})
	}

	return result, nil
}

// ConnectToRemote connects to a remote server (for Windows/remote mode)
func (a *App) ConnectToRemote(serverURL string) error {
	// Store remote server URL for the frontend
	a.config.AllowedDir = serverURL
	a.saveConfig()
	return nil
}

// Copy copies text to clipboard
func (a *App) Copy(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("cmd", "/c", "echo", text, "|", "clip")
	default:
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		io.WriteString(stdin, text)
		stdin.Close()
	}()
	return cmd.Run()
}

// GetHomeDir returns the user's home directory
func (a *App) GetHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// GetPlatformInfo returns platform information
func (a *App) GetPlatformInfo() map[string]string {
	return map[string]string{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"hasTmux": fmt.Sprintf("%v", hasTmux()),
	}
}

func hasTmux() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// ValidatePath checks if a path is valid and accessible
func (a *App) ValidatePath(path string) map[string]interface{} {
	info, err := os.Stat(path)
	result := map[string]interface{}{
		"valid": false,
	}

	if err != nil {
		result["error"] = err.Error()
		return result
	}

	result["valid"] = true
	result["isDir"] = info.IsDir()
	result["absolute"] = path

	abs, err := filepath.Abs(path)
	if err == nil {
		result["absolute"] = abs
	}

	// Check if path is within allowed dir
	if a.config.AllowedDir != "" {
		result["withinAllowed"] = strings.HasPrefix(abs, a.config.AllowedDir)
	}

	return result
}
