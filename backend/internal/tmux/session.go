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

package tmux

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExists   = errors.New("session already exists")
)

// Session 表示一个 tmux 会话
type Session struct {
	ID        string
	Name      string
	WorkDir   string
	CreatedAt time.Time
	mu        sync.RWMutex
}

// Manager 管理所有 tmux 会话
type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewManager 创建新的会话管理器
func NewManager() *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
	}
	// 启动时加载现有会话
	m.loadExistingSessions()
	return m
}

// loadExistingSessions 加载已存在的 tmux 会话
func (m *Manager) loadExistingSessions() {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	sessions := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, name := range sessions {
		if name == "" {
			continue
		}
		m.sessions[name] = &Session{
			ID:        generateSessionID(),
			Name:      name,
			CreatedAt: time.Now(),
		}
	}
}

// CreateSession 创建新的 tmux 会话
func (m *Manager) CreateSession(name, workDir string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[name]; exists {
		return nil, ErrSessionExists
	}

	// 使用 tmux 创建新会话
	args := []string{"new-session", "-d", "-s", name}
	if workDir != "" {
		args = append(args, "-c", workDir)
	}

	cmd := exec.Command("tmux", args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to create tmux session: %w", err)
	}

	session := &Session{
		ID:        generateSessionID(),
		Name:      name,
		WorkDir:   workDir,
		CreatedAt: time.Now(),
	}

	m.sessions[name] = session
	return session, nil
}

// GetSession 获取指定名称的会话
func (m *Manager) GetSession(name string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[name]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

// ListSessions 列出所有会话
func (m *Manager) ListSessions() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessions := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		sessions = append(sessions, s)
	}
	return sessions
}

// DeleteSession 删除指定会话
func (m *Manager) DeleteSession(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[name]; !exists {
		return ErrSessionNotFound
	}

	cmd := exec.Command("tmux", "kill-session", "-t", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill session: %w", err)
	}

	delete(m.sessions, name)
	return nil
}

// CaptureOutput 捕获会话输出
func (s *Session) CaptureOutput() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cmd := exec.Command("tmux", "capture-pane", "-t", s.Name, "-p", "-e")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to capture output: %w", err)
	}

	return string(output), nil
}

// SendKeys 发送按键到会话
func (s *Session) SendKeys(keys string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cmd := exec.Command("tmux", "send-keys", "-t", s.Name, keys)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send keys: %w", err)
	}

	return nil
}

// EnterCopyMode 进入 tmux copy mode
func (s *Session) EnterCopyMode() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use tmux command to enter copy mode directly
	cmd := exec.Command("tmux", "copy-mode", "-t", s.Name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enter copy mode: %w", err)
	}

	log.Printf("[Tmux] Entered copy mode for session %s", s.Name)
	return nil
}

// ScrollUp 在 copy mode 中向上滚动
func (s *Session) ScrollUp(lines int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < lines; i++ {
		cmd := exec.Command("tmux", "send-keys", "-t", s.Name, "-X", "scroll-up")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to scroll up: %w", err)
		}
	}

	return nil
}

// ScrollDown 在 copy mode 中向下滚动
func (s *Session) ScrollDown(lines int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < lines; i++ {
		cmd := exec.Command("tmux", "send-keys", "-t", s.Name, "-X", "scroll-down")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to scroll down: %w", err)
		}
	}

	return nil
}

// ExitCopyMode 退出 tmux copy mode
func (s *Session) ExitCopyMode() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Use tmux command to exit copy mode
	cmd := exec.Command("tmux", "send-keys", "-t", s.Name, "-X", "cancel")
	if err := cmd.Run(); err != nil {
		// Fallback: send q key
		cmd = exec.Command("tmux", "send-keys", "-t", s.Name, "q")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to exit copy mode: %w", err)
		}
	}

	return nil
}

// SendCommand 发送命令（带 Enter）
func (s *Session) SendCommand(cmd string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 发送命令，然后发送回车键
	// 先发送命令文本
	cmdExec := exec.Command("tmux", "send-keys", "-t", s.Name, cmd)
	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	// 然后发送回车键
	enterCmd := exec.Command("tmux", "send-keys", "-t", s.Name, "Enter")
	if err := enterCmd.Run(); err != nil {
		return fmt.Errorf("failed to send enter key: %w", err)
	}

	return nil
}

// IsActive 检查会话是否仍然活跃
func (s *Session) IsActive() bool {
	cmd := exec.Command("tmux", "has-session", "-t", s.Name)
	return cmd.Run() == nil
}

// GetPaneCount 获取会话中的 pane 数量
func (s *Session) GetPaneCount() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cmd := exec.Command("tmux", "display-message", "-t", s.Name, "-p", "#{window_panes}")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 1, nil // 默认返回 1
	}

	return count, nil
}

// StreamOutput 流式获取会话输出
func (s *Session) StreamOutput(lineCount int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cmd := exec.Command("tmux", "capture-pane", "-t", s.Name, "-p", "-e", "-S", fmt.Sprintf("-%d", lineCount))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to capture output: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	// 过滤空行
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

// generateSessionID 生成唯一的会话 ID
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

// IsTmuxAvailable 检查 tmux 是否可用
func IsTmuxAvailable() bool {
	cmd := exec.Command("which", "tmux")
	return cmd.Run() == nil
}
