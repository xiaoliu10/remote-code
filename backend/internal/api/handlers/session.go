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

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourname/remote-code/internal/security"
	"github.com/yourname/remote-code/internal/tmux"
)

// SessionHandler 会话处理器
type SessionHandler struct {
	tmuxManager *tmux.Manager
	validator   *security.SessionValidator
}

// NewSessionHandler 创建会话处理器
func NewSessionHandler(tmuxManager *tmux.Manager, validator *security.SessionValidator) *SessionHandler {
	return &SessionHandler{
		tmuxManager: tmuxManager,
		validator:   validator,
	}
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=32"`
	WorkDir string `json:"work_dir"`
}

// SessionResponse 会话响应
type SessionResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	WorkDir   string `json:"work_dir"`
	CreatedAt string `json:"created_at"`
	IsActive  bool   `json:"is_active"`
}

// CreateSession 创建新会话
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 验证会话名称
	if err := h.validator.ValidateSessionName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 验证工作目录
	if req.WorkDir != "" {
		if err := h.validator.ValidateWorkDir(req.WorkDir); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// 创建会话
	session, err := h.tmuxManager.CreateSession(req.Name, req.WorkDir)
	if err != nil {
		if err == tmux.ErrSessionExists {
			c.JSON(http.StatusConflict, gin.H{
				"error": "session already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	c.JSON(http.StatusCreated, SessionResponse{
		ID:        session.ID,
		Name:      session.Name,
		WorkDir:   session.WorkDir,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		IsActive:  true,
	})
}

// ListSessions 列出所有会话
func (h *SessionHandler) ListSessions(c *gin.Context) {
	sessions := h.tmuxManager.ListSessions()

	response := make([]SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		response = append(response, SessionResponse{
			ID:        s.ID,
			Name:      s.Name,
			WorkDir:   s.WorkDir,
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
			IsActive:  s.IsActive(),
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetSession 获取指定会话信息
func (h *SessionHandler) GetSession(c *gin.Context) {
	name := c.Param("name")

	session, err := h.tmuxManager.GetSession(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return
	}

	c.JSON(http.StatusOK, SessionResponse{
		ID:        session.ID,
		Name:      session.Name,
		WorkDir:   session.WorkDir,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		IsActive:  session.IsActive(),
	})
}

// DeleteSession 删除会话
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	name := c.Param("name")

	if err := h.tmuxManager.DeleteSession(name); err != nil {
		if err == tmux.ErrSessionNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "session not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete session",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSessionOutput 获取会话输出
func (h *SessionHandler) GetSessionOutput(c *gin.Context) {
	name := c.Param("name")

	session, err := h.tmuxManager.GetSession(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return
	}

	output, err := session.CaptureOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to capture output",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"output": output,
	})
}

// SendCommandRequest 发送命令请求
type SendCommandRequest struct {
	Command string `json:"command" binding:"required"`
}

// SendCommand 发送命令到会话
func (h *SessionHandler) SendCommand(c *gin.Context) {
	name := c.Param("name")

	var req SendCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// 验证命令安全性
	if err := h.validator.SanitizeCommand(req.Command); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session, err := h.tmuxManager.GetSession(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return
	}

	if err := session.SendCommand(req.Command); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send command",
		})
		return
	}

	c.Status(http.StatusOK)
}

// StreamOutputRequest 流式输出请求参数
type StreamOutputRequest struct {
	Lines int `form:"lines" binding:"min=1,max=1000"`
}

// StreamOutput 获取会话流式输出
func (h *SessionHandler) StreamOutput(c *gin.Context) {
	name := c.Param("name")

	var req StreamOutputRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.Lines = 100 // 默认 100 行
	}

	session, err := h.tmuxManager.GetSession(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "session not found",
		})
		return
	}

	lines, err := session.StreamOutput(req.Lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to stream output",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lines": lines,
	})
}
