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

import { ref, onUnmounted } from 'vue'

interface WSMessage {
  type: string
  data: unknown
  session?: string
}

type MessageHandler = (message: WSMessage) => void
type ConnectionHandler = () => void

/**
 * Get backend port from config
 */
function getBackendPort(): string {
  return (import.meta.env.VITE_BACKEND_PORT as string) || '9090'
}

/**
 * Get WebSocket base URL
 * Priority: VITE_WS_URL > VITE_BACKEND_PORT > auto-detect from current page
 */
function getWebSocketBaseUrl(): string {
  // 1. Use explicit full URL config if provided
  if (import.meta.env.VITE_WS_URL) {
    return import.meta.env.VITE_WS_URL as string
  }

  // 2. Auto-detect from current page location
  const { protocol, hostname, port } = window.location
  const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:'

  // If accessing via dev server (5173), use configured backend port
  // If accessing via same port, use the same port (production/reverse proxy)
  const backendPort = port === '5173' ? getBackendPort() : port

  return `${wsProtocol}//${hostname}:${backendPort}/api`
}

export function useWebSocket(sessionName: string) {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const error = ref<string | null>(null)

  const messageHandlers: MessageHandler[] = []
  const connectHandlers: ConnectionHandler[] = []
  const disconnectHandlers: ConnectionHandler[] = []

  const token = localStorage.getItem('token')
  const wsBaseUrl = getWebSocketBaseUrl()
  const wsUrl = `${wsBaseUrl}/ws/${sessionName}?token=${token}`

  const connect = () => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      return
    }

    try {
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        connected.value = true
        error.value = null
        connectHandlers.forEach((h) => h())
      }

      ws.value.onclose = () => {
        connected.value = false
        disconnectHandlers.forEach((h) => h())
      }

      ws.value.onerror = (event) => {
        error.value = 'WebSocket connection error'
        console.error('WebSocket error:', event)
      }

      ws.value.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data)
          messageHandlers.forEach((h) => h(message))
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e)
        }
      }
    } catch (e) {
      error.value = 'Failed to connect to WebSocket'
      console.error('WebSocket connection error:', e)
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    connected.value = false
  }

  const send = (type: string, data: unknown) => {
    if (ws.value?.readyState === WebSocket.OPEN) {
      ws.value.send(JSON.stringify({ type, data }))
    }
  }

  const sendCommand = (command: string) => {
    send('command', command)
  }

  const sendKeys = (keys: string) => {
    send('keys', keys)
  }

  const onMessage = (handler: MessageHandler) => {
    messageHandlers.push(handler)
  }

  const onConnect = (handler: ConnectionHandler) => {
    connectHandlers.push(handler)
  }

  const onDisconnect = (handler: ConnectionHandler) => {
    disconnectHandlers.push(handler)
  }

  // 清理
  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    error,
    connect,
    disconnect,
    send,
    sendCommand,
    sendKeys,
    onMessage,
    onConnect,
    onDisconnect
  }
}
