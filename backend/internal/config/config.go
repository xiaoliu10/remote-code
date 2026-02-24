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

package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Auth     AuthConfig
	Security SecurityConfig
	Tmux     TmuxConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AuthConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
	AdminPassword string // 默认管理员密码
}

type SecurityConfig struct {
	MaxSessionsPerUser int
	AllowedWorkDir     string
	EnableRateLimit    bool
	RateLimitRPS       int // 每秒请求数
	RateLimitBurst     int // 突发请求数
}

type TmuxConfig struct {
	SocketPath      string // tmux socket 路径
	ScrollbackLines int    // 终端历史缓冲区行数
}

func Load() *Config {
	// 首先尝试从 config.ini 加载
	loadConfigFromFile()

	return &Config{
		Server: ServerConfig{
			Port:         getEnv("BACKEND_PORT", "9090"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Auth: AuthConfig{
			JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production-please"),
			TokenDuration: 24 * time.Hour,
			AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
		},
		Security: SecurityConfig{
			MaxSessionsPerUser: 10,
			AllowedWorkDir:     getEnv("ALLOWED_DIR", os.Getenv("HOME")),
			EnableRateLimit:    getEnvBool("RATE_LIMIT_ENABLED", true),
			RateLimitRPS:       getEnvInt("RATE_LIMIT_RPS", 10),
			RateLimitBurst:     getEnvInt("RATE_LIMIT_BURST", 20),
		},
		Tmux: TmuxConfig{
			SocketPath:      getEnv("TMUX_SOCKET", ""),
			ScrollbackLines: getEnvInt("TERMINAL_SCROLLBACK", 1000),
		},
	}
}

// loadConfigFromFile loads configuration from ~/.remote-code/config.ini
func loadConfigFromFile() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	configFile := filepath.Join(homeDir, ".remote-code", "config.ini")

	file, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Only set if not already set in environment
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		return val == "true" || val == "1"
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		if _, err := fmt.Sscanf(val, "%d", &i); err == nil {
			return i
		}
	}
	return defaultVal
}
