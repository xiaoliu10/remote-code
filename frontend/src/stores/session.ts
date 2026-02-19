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

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Session } from '@/api/client'
import * as api from '@/api/client'

export const useSessionStore = defineStore('session', () => {
  // 状态
  const sessions = ref<Session[]>([])
  const currentSession = ref<Session | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const activeSessions = computed(() => sessions.value.filter((s) => s.is_active))

  // 方法
  async function fetchSessions() {
    loading.value = true
    error.value = null
    try {
      const response = await api.sessionApi.list()
      sessions.value = response.data
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch sessions'
    } finally {
      loading.value = false
    }
  }

  async function createSession(name: string, workDir?: string) {
    loading.value = true
    error.value = null
    try {
      const response = await api.sessionApi.create({ name, work_dir: workDir })
      sessions.value.push(response.data)
      return response.data
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to create session'
      error.value = message
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteSession(name: string) {
    loading.value = true
    error.value = null
    try {
      await api.sessionApi.delete(name)
      sessions.value = sessions.value.filter((s) => s.name !== name)
      if (currentSession.value?.name === name) {
        currentSession.value = null
      }
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to delete session'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function sendCommand(name: string, command: string) {
    try {
      await api.sessionApi.sendCommand(name, { command })
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to send command'
      throw e
    }
  }

  function setCurrentSession(session: Session | null) {
    currentSession.value = session
  }

  function getSessionByName(name: string) {
    return sessions.value.find((s) => s.name === name)
  }

  /**
   * Select a session by name
   */
  function selectSession(name: string) {
    const session = getSessionByName(name)
    if (session) {
      setCurrentSession(session)
    }
    return session
  }

  return {
    sessions,
    currentSession,
    loading,
    error,
    activeSessions,
    fetchSessions,
    createSession,
    deleteSession,
    sendCommand,
    setCurrentSession,
    getSessionByName,
    selectSession
  }
})
