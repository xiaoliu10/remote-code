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
 * "AS IS" BASIS, WITHOUT ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 * limitations under the License.
 */

import axios from 'axios'

/**
 * Get backend port from config
 */
function getBackendPort(): string {
  return (import.meta.env.VITE_BACKEND_PORT as string) || '9090'
}

/**
 * Get API base URL
 * Priority: VITE_API_URL > VITE_BACKEND_PORT > auto-detect from current page
 */
function getApiBaseUrl(): string {
  // 1. Use explicit full URL config if provided
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL as string
  }

  // 2. Use relative path (works with reverse proxy or same-origin)
  if (window.location.port !== '5173') {
    return '/api'
  }

  // 3. Dev mode: access backend on configured port
  const { protocol, hostname } = window.location
  const backendPort = getBackendPort()
  return `${protocol}//${hostname}:${backendPort}/api`
}

const BASE_URL = getApiBaseUrl()

// API 客户端
export const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 添加 token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 如果是登录接口本身返回 401，不跳转，让登录页面显示错误信息
      const isLoginRequest = error.config?.url?.includes('/auth/login')
      if (!isLoginRequest) {
        // Token 过期，清除本地存储并跳转到登录页
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// 类型定义
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user_id: string
  username: string
}

export interface Session {
  id: string
  name: string
  work_dir: string
  created_at: string
  is_active: boolean
}

export interface CreateSessionRequest {
  name: string
  work_dir?: string
}

export interface SendCommandRequest {
  command: string
}

// API 方法
export const authApi = {
  login: (data: LoginRequest) => api.post<LoginResponse>('/auth/login', data),
  validate: () => api.get('/auth/validate')
}

export const sessionApi = {
  list: () => api.get<Session[]>('/sessions'),
  get: (name: string) => api.get<Session>(`/sessions/${name}`),
  create: (data: CreateSessionRequest) => api.post<Session>('/sessions', data),
  delete: (name: string) => api.delete(`/sessions/${name}`),
  getOutput: (name: string) => api.get<{ output: string }>(`/sessions/${name}/output`),
  sendCommand: (name: string, data: SendCommandRequest) =>
    api.post(`/sessions/${name}/command`, data),
  streamOutput: (name: string, lines = 100) =>
    api.get<{ lines: string[] }>(`/sessions/${name}/stream`, { params: { lines } })
}

export const healthApi = {
  check: () => api.get<{ status: string; tmux: boolean }>('/health')
}

// File operation types - 匹配后端响应格式
export interface FileItem {
  name: string
  path: string
  type: 'file' | 'directory'
  size: number
  modTime: string
  permission: string
}

export interface FileListResponse {
  path: string
  items: FileItem[]
  total: number
  page: number
  pageSize: number
}

export interface FileContent {
  path: string
  content: string
  size: number
  modTime: string
}

export interface CreateFileRequest {
  path: string
  type: 'file' | 'directory'
  content?: string
}

export interface RenameFileRequest {
  oldPath: string
  newPath: string
}

// File API
export const fileApi = {
  list: (path: string, page?: number, pageSize?: number) =>
    api.get<FileListResponse>('/files', { params: { path, page, pageSize } }),

  getContent: (path: string) =>
    api.get<FileContent>('/files/content', { params: { path } }),

  create: (path: string, type: 'file' | 'directory', content?: string) =>
    api.post<FileItem>('/files', { path, type, content }),

  rename: (oldPath: string, newPath: string) =>
    api.put<FileItem>('/files/rename', { oldPath, newPath }),

  delete: (path: string) =>
    api.delete('/files', { params: { path } })
}
