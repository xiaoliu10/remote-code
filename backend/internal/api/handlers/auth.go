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

	"github.com/gin-gonic/gin"
	"github.com/xiaoliu10/remote-code/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	jwtManager    *auth.JWTManager
	adminPassword string // 默认管理员密码的哈希
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(jwtManager *auth.JWTManager, adminPassword string) *AuthHandler {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	return &AuthHandler{
		jwtManager:    jwtManager,
		adminPassword: string(hashedPassword),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
}

// Login 处理登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request format",
		})
		return
	}

	// 简单验证：用户名为 "admin"，密码匹配配置的密码
	if req.Username != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(h.adminPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})
		return
	}

	// 生成 token
	token, err := h.jwtManager.Generate("admin", req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// 返回 token
	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: 0, // 前端可以从 JWT 解析
		UserID:    "admin",
		Username:  req.Username,
	})
}

// ValidateToken 验证 token 有效性
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// 如果能到达这里，说明中间件已经验证过了
	userID := c.GetString("user_id")
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"user_id":  userID,
		"username": username,
	})
}
