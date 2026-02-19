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

package security

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	// 会话名称只允许字母、数字、下划线、短横线
	sessionNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,32}$`)
	// 危险命令检测
	dangerousCommands = []string{
		"rm -rf /",
		"rm -rf /*",
		"mkfs",
		":(){:|:&};:", // fork bomb
		"dd if=/dev/zero",
		"> /dev/sda",
		":w !sudo tee %",
	}
)

type SessionValidator struct {
	allowedWorkDir string
	mu             sync.RWMutex
}

func NewSessionValidator(allowedDir string) *SessionValidator {
	return &SessionValidator{allowedWorkDir: allowedDir}
}

func (v *SessionValidator) ValidateSessionName(name string) error {
	if !sessionNameRegex.MatchString(name) {
		return errors.New("invalid session name: only alphanumeric, underscore and hyphen allowed (1-32 chars)")
	}
	return nil
}

func (v *SessionValidator) ValidateWorkDir(path string) error {
	v.mu.RLock()
	allowedDir := v.allowedWorkDir
	v.mu.RUnlock()

	// 清理路径，解析 .. 和 .
	cleanPath := filepath.Clean(path)

	// 转换为绝对路径用于比较
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	absAllowed, err := filepath.Abs(allowedDir)
	if err != nil {
		return fmt.Errorf("invalid allowed directory: %w", err)
	}

	// 检查是否在允许的目录内
	if !strings.HasPrefix(absPath, absAllowed) {
		return fmt.Errorf("path '%s' is outside allowed directory '%s'", path, allowedDir)
	}

	return nil
}

func (v *SessionValidator) SanitizeCommand(cmd string) error {
	cmdLower := strings.ToLower(cmd)

	for _, dangerous := range dangerousCommands {
		if strings.Contains(cmdLower, strings.ToLower(dangerous)) {
			return fmt.Errorf("command contains dangerous pattern: %s", dangerous)
		}
	}

	return nil
}

// 验证命令参数是否安全
func (v *SessionValidator) ValidateCommandArgs(args []string) error {
	for _, arg := range args {
		// 检查是否包含管道或重定向
		if strings.ContainsAny(arg, "|&;<>$`") {
			return errors.New("command arguments contain special characters")
		}
	}
	return nil
}
