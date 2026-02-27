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
	sessions    map[string]*Session
	mu          sync.RWMutex
	persistence *Persistence
}

// NewManager 创建新的会话管理器
func NewManager(dataDir string) *Manager {
	m := &Manager{
		sessions:    make(map[string]*Session),
		persistence: NewPersistence(dataDir),
	}
	// 启动时加载现有会话
	m.loadExistingSessions()
	// 恢复持久化的会话
	m.restoreSessions()
	// 将内存中的会话同步到持久化文件（处理服务重启前未持久化的情况）
	m.syncExistingSessionsToPersistence()
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
		log.Printf("[Tmux] Loaded existing tmux session: %s", name)
	}

	log.Printf("[Tmux] Loaded %d existing tmux session(s)", len(m.sessions))
}

// syncExistingSessionsToPersistence 将已存在的 tmux 会话同步到持久化文件
func (m *Manager) syncExistingSessionsToPersistence() {
	// 加载已持久化的会话列表
	persisted, err := m.persistence.LoadSessions()
	if err != nil {
		log.Printf("[Tmux] Failed to load persisted sessions for sync: %v", err)
		return
	}

	// 创建已持久化会话的 map
	persistedMap := make(map[string]bool)
	for _, meta := range persisted {
		persistedMap[meta.Name] = true
	}

	// 将内存中的会话同步到持久化文件
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 收集有效的会话（同时存在于内存和 tmux 中的）
	validSessions := make([]SessionMetadata, 0)

	for name, session := range m.sessions {
		// 检查 tmux 会话是否真实存在
		cmd := exec.Command("tmux", "has-session", "-t", name)
		if cmd.Run() != nil {
			log.Printf("[Tmux] Session %s no longer exists in tmux, removing from memory", name)
			delete(m.sessions, name)
			continue
		}

		validSessions = append(validSessions, SessionMetadata{
			Name:      name,
			WorkDir:   session.WorkDir,
			CreatedAt: session.CreatedAt,
		})

		if !persistedMap[name] {
			// 该会话未被持久化，记录日志
			log.Printf("[Tmux] Synced session %s to persistence", name)
		}
	}

	// 如果有效会话数量与持久化文件不同，更新持久化文件
	if len(validSessions) != len(persisted) {
		log.Printf("[Tmux] Updating persistence file: %d valid sessions, %d persisted", len(validSessions), len(persisted))
		if err := m.persistence.SaveSessions(validSessions); err != nil {
			log.Printf("[Tmux] Failed to update persistence file: %v", err)
		}
	}
}

// restoreSessions 从持久化数据恢复会话
func (m *Manager) restoreSessions() {
	metadata, err := m.persistence.LoadSessions()
	if err != nil {
		log.Printf("[Tmux] Failed to load persisted sessions: %v", err)
		return
	}

	log.Printf("[Tmux] Loaded %d persisted session(s) from file", len(metadata))

	for _, meta := range metadata {
		// 检查会话是否已存在（可能是 loadExistingSessions 加载的）
		if session, exists := m.sessions[meta.Name]; exists {
			// 会话已存在，更新其 WorkDir 信息
			session.WorkDir = meta.WorkDir
			session.CreatedAt = meta.CreatedAt
			log.Printf("[Tmux] Updated existing session %s with persisted metadata (work_dir: %s)", meta.Name, meta.WorkDir)
			continue
		}

		// 尝试重新创建 tmux 会话
		args := []string{"new-session", "-d", "-s", meta.Name}
		if meta.WorkDir != "" {
			args = append(args, "-c", meta.WorkDir)
		}

		cmd := exec.Command("tmux", args...)
		if err := cmd.Run(); err != nil {
			log.Printf("[Tmux] Failed to restore session %s: %v", meta.Name, err)
			continue
		}

		log.Printf("[Tmux] Restored session: %s (work_dir: %s)", meta.Name, meta.WorkDir)
		m.sessions[meta.Name] = &Session{
			ID:        generateSessionID(),
			Name:      meta.Name,
			WorkDir:   meta.WorkDir,
			CreatedAt: meta.CreatedAt,
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

	// 持久化会话元数据
	if err := m.persistence.AddSession(SessionMetadata{
		Name:      name,
		WorkDir:   workDir,
		CreatedAt: session.CreatedAt,
	}); err != nil {
		log.Printf("[Tmux] Failed to persist session %s: %v", name, err)
	}

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

	// 清理持久化数据
	if err := m.persistence.RemoveSession(name); err != nil {
		log.Printf("[Tmux] Failed to remove persisted session %s: %v", name, err)
	}

	return nil
}

// CaptureOutput 捕获会话输出（包括历史缓冲区）
func (s *Session) CaptureOutput() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Capture visible area plus history (up to 5000 lines back)
	cmd := exec.Command("tmux", "capture-pane", "-t", s.Name, "-p", "-e", "-S", "-5000")
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
