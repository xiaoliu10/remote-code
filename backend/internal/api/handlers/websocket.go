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
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"github.com/xiaoliu10/remote-code/internal/tmux"
	"github.com/xiaoliu10/remote-code/internal/websocket"
)

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	hub        *websocket.Hub
	tmuxManager *tmux.Manager
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler(hub *websocket.Hub, tmuxManager *tmux.Manager) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		tmuxManager: tmuxManager,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 从 JWT 获取用户信息
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sessionName := c.Param("session")
	if sessionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session name required"})
		return
	}

	// 检查会话是否存在
	session, err := h.tmuxManager.GetSession(sessionName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	// 升级到 WebSocket
	conn, err := (&gorillaws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // 生产环境应设置严格的 Origin 检查
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Failed to upgrade: %v", err)
		return
	}

	// 创建客户端
	client := &websocket.Client{
		Hub:       h.hub,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		SessionID: session.Name,
		UserID:    userID,
	}

	// 注册客户端
	h.hub.Register(client)

	// 启动读写协程
	go client.WritePump()
	go h.readPumpWithOutput(client, session)
}

// readPumpWithOutput 读取客户端消息并定期发送终端输出
func (h *WebSocketHandler) readPumpWithOutput(client *websocket.Client, session *tmux.Session) {
	defer func() {
		h.hub.Unregister(client)
		client.Conn.Close()
		// 注意：不要在这里关闭 client.Send，Hub 的 unregisterClient 会负责关闭
	}()

	// 设置读取参数
	client.Conn.SetReadLimit(512 * 1024)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 定期发送终端输出
	outputTicker := time.NewTicker(500 * time.Millisecond)
	defer outputTicker.Stop()

	// 用于检测输出的变化
	lastOutput := ""

	// 创建一个 channel 用于接收客户端消息
	messageChan := make(chan []byte, 10)
	go func() {
		for {
			_, message, err := client.Conn.ReadMessage()
			if err != nil {
				close(messageChan)
				return
			}
			messageChan <- message
		}
	}()

	for {
		select {
		case <-outputTicker.C:
			// 获取当前输出
			currentOutput, err := session.CaptureOutput()
			if err != nil {
				h.hub.SendToSession(session.Name, "error", "Failed to capture output")
				continue
			}

			// 只在输出变化时发送
			if currentOutput != lastOutput {
				// 发送增量更新（这里简化为发送全部）
				// 生产环境应该实现增量更新以减少带宽
				h.hub.SendToSession(session.Name, "output", map[string]interface{}{
					"text":      currentOutput,
					"timestamp": time.Now().Unix(),
				})
				lastOutput = currentOutput
			}

		case message, ok := <-messageChan:
			if !ok {
				// 客户端连接已关闭
				return
			}
			// 处理客户端消息
			h.handleClientMessage(client, session, message)
			// 重置读取超时
			client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}
	}
}

// handleClientMessage 处理来自客户端的消息
func (h *WebSocketHandler) handleClientMessage(client *websocket.Client, session *tmux.Session, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("[WS] Invalid message format: %v", err)
		return
	}

	msgType, _ := msg["type"].(string)

	switch msgType {
	case "command":
		// 发送命令到会话
		command, _ := msg["data"].(string)
		log.Printf("[WS] Received command: %q for session %s", command, session.Name)
		if command != "" {
			if err := session.SendCommand(command); err != nil {
				log.Printf("[WS] Failed to send command: %v", err)
				h.hub.SendToSession(session.Name, "error", "Failed to send command")
			} else {
				log.Printf("[WS] Command sent successfully")
				h.hub.SendToSession(session.Name, "status", "Command sent")
			}
		}

	case "keys":
		// 发送按键到会话（不回车）
		keys, _ := msg["data"].(string)
		log.Printf("[WS] Received keys: %q for session %s", keys, session.Name)
		if keys != "" {
			if err := session.SendKeys(keys); err != nil {
				log.Printf("[WS] Failed to send keys: %v", err)
				h.hub.SendToSession(session.Name, "error", "Failed to send keys")
			} else {
				log.Printf("[WS] Keys sent successfully")
			}
		}

	case "resize":
		// 处理终端大小调整（暂未实现）
		log.Printf("[WS] Resize requested: %+v", msg)

	case "ping":
		h.hub.SendToSession(session.Name, "pong", nil)

	default:
		log.Printf("[WS] Unknown message type: %s", msgType)
	}
}
