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
 * Unless required by law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package tmux

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SessionMetadata 会话元数据，用于持久化存储
type SessionMetadata struct {
	Name      string    `json:"name"`
	WorkDir   string    `json:"work_dir"`
	CreatedAt time.Time `json:"created_at"`
}

// Persistence 会话持久化管理器
type Persistence struct {
	dataDir  string
	filePath string
	mu       sync.RWMutex
}

// NewPersistence 创建持久化管理器
func NewPersistence(dataDir string) *Persistence {
	p := &Persistence{
		dataDir:  dataDir,
		filePath: filepath.Join(dataDir, "sessions.json"),
	}
	// 确保数据目录存在
	os.MkdirAll(dataDir, 0755)
	return p
}

// LoadSessions 从文件加载会话元数据
func (p *Persistence) LoadSessions() ([]SessionMetadata, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	data, err := os.ReadFile(p.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []SessionMetadata{}, nil
		}
		return nil, err
	}

	var sessions []SessionMetadata
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// SaveSessions 保存所有会话元数据到文件
func (p *Persistence) SaveSessions(sessions []SessionMetadata) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	data, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p.filePath, data, 0644)
}

// AddSession 添加单个会话元数据
func (p *Persistence) AddSession(meta SessionMetadata) error {
	sessions, err := p.LoadSessions()
	if err != nil {
		return err
	}

	// 检查是否已存在
	for i, s := range sessions {
		if s.Name == meta.Name {
			// 更新现有会话
			sessions[i] = meta
			return p.SaveSessions(sessions)
		}
	}

	// 添加新会话
	sessions = append(sessions, meta)
	return p.SaveSessions(sessions)
}

// RemoveSession 移除会话元数据
func (p *Persistence) RemoveSession(name string) error {
	sessions, err := p.LoadSessions()
	if err != nil {
		return err
	}

	// 过滤掉要删除的会话
	newSessions := make([]SessionMetadata, 0, len(sessions))
	for _, s := range sessions {
		if s.Name != name {
			newSessions = append(newSessions, s)
		}
	}

	return p.SaveSessions(newSessions)
}
