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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/remote-code/internal/api/handlers"
	"github.com/yourname/remote-code/internal/api/middleware"
	"github.com/yourname/remote-code/internal/auth"
	"github.com/yourname/remote-code/internal/config"
	"github.com/yourname/remote-code/internal/security"
	"github.com/yourname/remote-code/internal/tmux"
	"github.com/yourname/remote-code/internal/websocket"
)

// RouterConfig 路由配置
type RouterConfig struct {
	JWTManager    *auth.JWTManager
	TmuxManager   *tmux.Manager
	Validator     *security.SessionValidator
	Hub           *websocket.Hub
	AdminPassword string
	Config        *config.Config
}

// SetupRouter 设置路由
func SetupRouter(cfg *RouterConfig) *gin.Engine {
	router := gin.Default()

	// 中间件
	router.Use(middleware.CORS())

	// 创建 handlers
	authHandler := handlers.NewAuthHandler(cfg.JWTManager, cfg.AdminPassword)
	sessionHandler := handlers.NewSessionHandler(cfg.TmuxManager, cfg.Validator)
	wsHandler := handlers.NewWebSocketHandler(cfg.Hub, cfg.TmuxManager)

	// 创建文件处理器
	pathValidator, err := handlers.NewPathValidator(cfg.Config.Security.AllowedWorkDir)
	if err != nil {
		panic("failed to create path validator: " + err.Error())
	}
	fileHandler := handlers.NewFileHandler(pathValidator)

	// 公开路由
	public := router.Group("/api")
	{
		public.POST("/auth/login", authHandler.Login)
	}

	// 需要 JWT 认证的路由
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTManager))
	{
		// 认证相关
		protected.GET("/auth/validate", authHandler.ValidateToken)

		// 会话管理
		protected.POST("/sessions", sessionHandler.CreateSession)
		protected.GET("/sessions", sessionHandler.ListSessions)
		protected.GET("/sessions/:name", sessionHandler.GetSession)
		protected.DELETE("/sessions/:name", sessionHandler.DeleteSession)
		protected.GET("/sessions/:name/output", sessionHandler.GetSessionOutput)
		protected.POST("/sessions/:name/command", sessionHandler.SendCommand)
		protected.GET("/sessions/:name/stream", sessionHandler.StreamOutput)

		// WebSocket
		protected.GET("/ws/:session", wsHandler.HandleWebSocket)

		// 文件系统操作
		files := protected.Group("/files")
		files.GET("", fileHandler.ListDirectory)
		files.GET("/content", fileHandler.GetFileContent)
		files.POST("", fileHandler.CreateFileFolder)
		files.PUT("/rename", fileHandler.RenameFileFolder)
		files.DELETE("", fileHandler.DeleteFileFolder)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"tmux":   tmux.IsTmuxAvailable(),
		})
	})

	return router
}
