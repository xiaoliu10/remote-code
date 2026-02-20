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

package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

// Message WebSocket 消息类型
type Message struct {
	Type    string      `json:"type"`    // output, command, status, error
	Data    interface{} `json:"data"`
	Session string      `json:"session,omitempty"`
}

// Client WebSocket 客户端连接
type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	SessionID string
	UserID    string
	closed    bool
	closeMu   sync.Mutex
}

// SafeClose 安全关闭客户端连接
func (c *Client) SafeClose() {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	if c.closed {
		return
	}
	c.closed = true
	close(c.Send)
}

// Hub WebSocket 连接池管理器
type Hub struct {
	clients    map[string]*Client // sessionID -> Client
	userClients map[string]map[string]*Client // userID -> sessionID -> Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		userClients: make(map[string]map[string]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan Message),
	}
}

// Run 运行 Hub 主循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient 注册新客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 如果同一会话已有连接，发送踢出消息给旧连接
	if existing, exists := h.clients[client.SessionID]; exists {
		// 发送踢出消息
		kickMsg := Message{
			Type: "kicked",
			Data: "Your connection has been replaced by a new connection from another device",
		}
		if data, err := json.Marshal(kickMsg); err == nil {
			select {
			case existing.Send <- data:
			default:
			}
		}
		// 给旧连接一点时间接收消息
		time.Sleep(100 * time.Millisecond)
		existing.Conn.Close()
		existing.SafeClose()
	}

	h.clients[client.SessionID] = client

	// 按用户索引
	if h.userClients[client.UserID] == nil {
		h.userClients[client.UserID] = make(map[string]*Client)
	}
	h.userClients[client.UserID][client.SessionID] = client

	log.Printf("[Hub] Client registered: session=%s user=%s", client.SessionID, client.UserID)
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.SessionID]; ok {
		delete(h.clients, client.SessionID)
		client.SafeClose()
		client.Conn.Close()
	}

	if userMap, ok := h.userClients[client.UserID]; ok {
		delete(userMap, client.SessionID)
		if len(userMap) == 0 {
			delete(h.userClients, client.UserID)
		}
	}

	log.Printf("[Hub] Client unregistered: session=%s user=%s", client.SessionID, client.UserID)
}

// broadcastMessage 广播消息
func (h *Hub) broadcastMessage(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("[Hub] Failed to marshal message: %v", err)
		return
	}

	// 发送给目标会话的客户端
	if client, ok := h.clients[message.Session]; ok {
		select {
		case client.Send <- data:
		default:
			log.Printf("[Hub] Send channel full for session=%s", message.Session)
		}
	}
}

// SendToSession 发送消息到指定会话
func (h *Hub) SendToSession(sessionID string, msgType string, data interface{}) {
	h.broadcast <- Message{
		Type:    msgType,
		Data:    data,
		Session: sessionID,
	}
}

// GetUserSessions 获取用户的所有活跃会话
func (h *Hub) GetUserSessions(userID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessions := make([]string, 0)
	if userMap, ok := h.userClients[userID]; ok {
		for sessionID := range userMap {
			sessions = append(sessions, sessionID)
		}
	}
	return sessions
}

// IsSessionActive 检查会话是否有活跃的 WebSocket 连接
func (h *Hub) IsSessionActive(sessionID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[sessionID]
	return ok
}

// ReadPump 从 WebSocket 连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Unexpected close: %v", err)
			}
			break
		}

		// 处理接收到的消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err == nil {
			c.handleMessage(msg)
		}
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("[WS] Write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "command":
		// 命令消息，转发给处理器
		if cmd, ok := msg.Data.(string); ok {
			log.Printf("[WS] Command received: %s", cmd)
			// 这里可以触发命令执行逻辑
		}
	case "ping":
		c.Hub.SendToSession(c.SessionID, "pong", nil)
	}
}

// Upgrader WebSocket 升级器
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应设置严格的 Origin 检查
	},
}

// Register 注册客户端（导出方法）
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端（导出方法）
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
