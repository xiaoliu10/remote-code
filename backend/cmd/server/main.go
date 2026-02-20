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

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yourname/remote-code/internal/api"
	"github.com/yourname/remote-code/internal/api/middleware"
	"github.com/yourname/remote-code/internal/auth"
	"github.com/yourname/remote-code/internal/config"
	"github.com/yourname/remote-code/internal/security"
	"github.com/yourname/remote-code/internal/tmux"
	"github.com/yourname/remote-code/internal/websocket"
	"golang.org/x/time/rate"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 加载配置
	cfg := config.Load()

	// 检查 tmux 是否可用
	if !tmux.IsTmuxAvailable() {
		log.Fatal("tmux is not available. Please install tmux first.")
	}

	// 初始化组件
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.TokenDuration)
	tmuxManager := tmux.NewManager()
	validator := security.NewSessionValidator(cfg.Security.AllowedWorkDir)
	wsHub := websocket.NewHub()

	// 启动 WebSocket Hub
	go wsHub.Run()

	// 创建速率限制器
	var rateLimitMiddleware gin.HandlerFunc
	if cfg.Security.EnableRateLimit {
		limiter := middleware.NewIPRateLimiter(
			rate.Limit(cfg.Security.RateLimitRPS),
			cfg.Security.RateLimitBurst,
		)
		// 定期清理 IP 记录
		limiter.CleanupStaleIPs(5 * time.Minute)
		rateLimitMiddleware = limiter.Middleware()
	}

	// 设置路由
	routerConfig := &api.RouterConfig{
		JWTManager:    jwtManager,
		TmuxManager:   tmuxManager,
		Validator:     validator,
		Hub:           wsHub,
		AdminPassword: cfg.Auth.AdminPassword,
		Config:        cfg,
	}
	router := api.SetupRouter(routerConfig)

	// 应用速率限制
	if rateLimitMiddleware != nil {
		router.Use(rateLimitMiddleware)
	}

	// 创建服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		log.Printf("Starting Remote Claude Code Server on port %s", cfg.Server.Port)
		log.Printf("Allowed working directory: %s", cfg.Security.AllowedWorkDir)
		log.Printf("Rate limiting: %v (%d RPS, burst %d)",
			cfg.Security.EnableRateLimit,
			cfg.Security.RateLimitRPS,
			cfg.Security.RateLimitBurst,
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
