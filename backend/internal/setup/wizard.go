/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package setup

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the setup configuration
type Config struct {
	AdminPassword string
	AllowedDir    string
	BackendPort   string
	EnableFRP     bool
	FRPServerAddr string
	FRPServerPort string
	FRPToken      string
}

// Wizard handles interactive setup
type Wizard struct {
	configDir string
	configFile string
	reader    *bufio.Reader
}

// NewWizard creates a new setup wizard
func NewWizard() *Wizard {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".remote-code")

	return &Wizard{
		configDir:  configDir,
		configFile: filepath.Join(configDir, "config.ini"),
		reader:     bufio.NewReader(os.Stdin),
	}
}

// NeedsSetup checks if setup is needed
func (w *Wizard) NeedsSetup() bool {
	_, err := os.Stat(w.configFile)
	return os.IsNotExist(err)
}

// Run executes the setup wizard
func (w *Wizard) Run() (*Config, error) {
	fmt.Println("")
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║          Welcome to Remote Code Setup Wizard            ║")
	fmt.Println("║                                                          ║")
	fmt.Println("║  This wizard will help you configure Remote Code        ║")
	fmt.Println("║  for the first time.                                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println("")

	// Create config directory
	if err := os.MkdirAll(w.configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	cfg := &Config{
		BackendPort: "9090",
	}

	// Admin password
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📝 Step 1: Admin Password")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Println("Choose an admin password for logging into the web interface.")
	fmt.Println("Press Enter to generate a random password.")
	fmt.Println("")

	password := w.promptWithDefault("Admin password", "(generate random)")
	if password == "" || password == "(generate random)" {
		cfg.AdminPassword = w.generateRandomPassword()
		fmt.Printf("✅ Generated password: %s\n", cfg.AdminPassword)
		fmt.Println("⚠️  Please save this password!")
	} else {
		cfg.AdminPassword = password
	}
	fmt.Println("")

	// Allowed directory
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📁 Step 2: Working Directory")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Println("Choose the directory where sessions can be created.")
	fmt.Println("Sessions will only be able to access files within this directory.")
	fmt.Println("")

	homeDir, _ := os.UserHomeDir()
	defaultDir := filepath.Join(homeDir, "projects")

	cfg.AllowedDir = w.promptWithDefault("Working directory", defaultDir)
	fmt.Println("")

	// FRP Configuration
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🌐 Step 3: Remote Access (FRP)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Println("FRP allows you to access Remote Code from outside your network.")
	fmt.Println("If you have a public server, you can configure FRP now or skip.")
	fmt.Println("")
	fmt.Println("Do you want to configure FRP for remote access?")

	enableFRP := w.promptYesNo("Enable FRP", false)
	cfg.EnableFRP = enableFRP

	if enableFRP {
		fmt.Println("")
		fmt.Println("Enter your FRP server details:")
		cfg.FRPServerAddr = w.promptRequired("FRP server address")
		cfg.FRPServerPort = w.promptWithDefault("FRP server port", "7000")
		cfg.FRPToken = w.promptRequired("FRP token")
	}
	fmt.Println("")

	// Summary
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Configuration Summary")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("")
	fmt.Printf("  Admin Password:    %s\n", cfg.AdminPassword)
	fmt.Printf("  Working Directory: %s\n", cfg.AllowedDir)
	fmt.Printf("  Backend Port:      %s\n", cfg.BackendPort)
	fmt.Printf("  FRP Enabled:       %v\n", cfg.EnableFRP)
	if cfg.EnableFRP {
		fmt.Printf("  FRP Server:        %s:%s\n", cfg.FRPServerAddr, cfg.FRPServerPort)
	}
	fmt.Println("")

	// Confirm
	fmt.Println("Save this configuration?")
	if !w.promptYesNo("Save configuration", true) {
		fmt.Println("❌ Setup cancelled.")
		os.Exit(0)
	}

	// Save configuration
	if err := w.saveConfig(cfg); err != nil {
		return nil, err
	}

	fmt.Println("")
	fmt.Println("✅ Configuration saved!")
	fmt.Println("")

	return cfg, nil
}

func (w *Wizard) promptWithDefault(prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, _ := w.reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}

func (w *Wizard) promptRequired(prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		input, _ := w.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			return input
		}
		fmt.Println("⚠️  This field is required. Please enter a value.")
	}
}

func (w *Wizard) promptYesNo(prompt string, defaultYes bool) bool {
	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}

	fmt.Printf("%s [%s]: ", prompt, defaultStr)
	input, _ := w.reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}

func (w *Wizard) generateRandomPassword() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (w *Wizard) saveConfig(cfg *Config) error {
	var content strings.Builder

	content.WriteString("# Remote Code Configuration\n")
	content.WriteString("# 配置文件路径: " + w.configFile + "\n")
	content.WriteString("#\n")
	content.WriteString("# 修改此文件后重启服务生效\n")
	content.WriteString("\n")
	content.WriteString("# ==================== Backend Configuration ====================\n")
	content.WriteString("# 后端服务配置\n")
	content.WriteString("\n")
	content.WriteString("BACKEND_PORT=" + cfg.BackendPort + "\n")
	content.WriteString("JWT_SECRET=" + w.generateRandomPassword() + "\n")
	content.WriteString("ADMIN_PASSWORD=" + cfg.AdminPassword + "\n")
	content.WriteString("ALLOWED_DIR=" + cfg.AllowedDir + "\n")
	content.WriteString("\n")
	content.WriteString("# Rate limiting\n")
	content.WriteString("RATE_LIMIT_ENABLED=true\n")
	content.WriteString("RATE_LIMIT_RPS=10\n")
	content.WriteString("RATE_LIMIT_BURST=20\n")
	content.WriteString("\n")
	content.WriteString("# ==================== Frontend Configuration ====================\n")
	content.WriteString("# 前端服务配置\n")
	content.WriteString("\n")
	content.WriteString("FRONTEND_PORT=5173\n")
	content.WriteString("\n")
	content.WriteString("# ==================== FRP Configuration ====================\n")
	content.WriteString("# FRP 内网穿透配置\n")
	content.WriteString("\n")

	if cfg.EnableFRP {
		content.WriteString("FRP_ENABLED=true\n")
		content.WriteString("FRP_SERVER_ADDR=" + cfg.FRPServerAddr + "\n")
		content.WriteString("FRP_SERVER_PORT=" + cfg.FRPServerPort + "\n")
		content.WriteString("FRP_TOKEN=" + cfg.FRPToken + "\n")
	} else {
		content.WriteString("FRP_ENABLED=false\n")
		content.WriteString("FRP_SERVER_ADDR=\n")
		content.WriteString("FRP_SERVER_PORT=7000\n")
		content.WriteString("FRP_TOKEN=\n")
	}

	content.WriteString("\n")
	content.WriteString("# ==================== Advanced Settings ====================\n")
	content.WriteString("# 高级设置\n")
	content.WriteString("\n")
	content.WriteString("LOG_DIR=" + w.configDir + "/logs\n")
	content.WriteString("PID_DIR=" + w.configDir + "/logs\n")

	return os.WriteFile(w.configFile, []byte(content.String()), 0600)
}
