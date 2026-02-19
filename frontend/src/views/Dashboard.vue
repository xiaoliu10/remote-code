<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->

<template>
  <div class="dashboard-container">
    <!-- 头部 -->
    <header class="header">
      <div class="header-left">
        <h1>{{ t('dashboard.title') }}</h1>
      </div>
      <div class="header-right">
        <LanguageSwitcher />
        <n-button quaternary @click="handleRefresh">
          <template #icon>
            <n-icon><RefreshIcon /></n-icon>
          </template>
          {{ t('common.refresh') }}
        </n-button>
        <n-button quaternary @click="handleLogout">
          <template #icon>
            <n-icon><LogoutIcon /></n-icon>
          </template>
          {{ t('common.logout') }}
        </n-button>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <!-- 创建会话 -->
      <n-card :title="t('dashboard.createSession')" class="create-card">
        <n-form inline :model="createForm" label-placement="left">
          <n-form-item :label="t('dashboard.sessionName')">
            <n-input
              v-model:value="createForm.name"
              :placeholder="t('dashboard.sessionNamePlaceholder')"
              style="width: 200px"
            />
          </n-form-item>
          <n-form-item :label="t('dashboard.workDirectory')">
            <n-input
              v-model:value="createForm.workDir"
              placeholder="/Users/jason/projects (optional)"
              style="width: 300px"
            />
          </n-form-item>
          <n-form-item>
            <n-button type="primary" :loading="sessionStore.loading" @click="handleCreateSession">
              {{ t('common.create') }}
            </n-button>
          </n-form-item>
        </n-form>
      </n-card>

      <!-- 会话列表 -->
      <n-card :title="t('dashboard.sessions')" class="sessions-card">
        <n-spin :show="sessionStore.loading">
          <n-empty v-if="sessionStore.sessions.length === 0" :description="t('dashboard.noSessions')" />
          <n-list v-else hoverable clickable>
            <n-list-item v-for="session in sessionStore.sessions" :key="session.name">
              <template #prefix>
                <n-icon :component="session.is_active ? TerminalIcon : StopIcon" />
              </template>
              <n-thing :title="session.name">
                <template #description>
                  <n-space align="center">
                    <n-tag :type="session.is_active ? 'success' : 'default'" size="small">
                      {{ session.is_active ? t('dashboard.active') : t('dashboard.inactive') }}
                    </n-tag>
                    <span class="work-dir">{{ session.work_dir || t('dashboard.default') }}</span>
                    <span class="created-at">{{ formatTime(session.created_at) }}</span>
                  </n-space>
                </template>
              </n-thing>
              <template #suffix>
                <n-space>
                  <n-button
                    type="primary"
                    size="small"
                    @click="handleOpenSession(session)"
                  >
                    {{ t('common.open') }}
                  </n-button>
                  <n-popconfirm @positive-click="handleDeleteSession(session.name)">
                    <template #trigger>
                      <n-button size="small" type="error" tertiary>
                        {{ t('common.delete') }}
                      </n-button>
                    </template>
                    <span>{{ t('dashboard.deleteConfirm') }}</span>
                  </n-popconfirm>
                </n-space>
              </template>
            </n-list-item>
          </n-list>
        </n-spin>
      </n-card>
    </main>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NList,
  NListItem,
  NThing,
  NTag,
  NSpace,
  NSpin,
  NEmpty,
  NIcon,
  NPopconfirm,
  useMessage
} from 'naive-ui'
import { Terminal as TerminalIcon, LogOut as LogoutIcon, Refresh as RefreshIcon, Ellipse as StopIcon } from '@vicons/ionicons5'
import { useSessionStore } from '@/stores/session'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const sessionStore = useSessionStore()

const createForm = reactive({
  name: '',
  workDir: ''
})

onMounted(() => {
  sessionStore.fetchSessions()
})

async function handleCreateSession() {
  if (!createForm.name) {
    message.error(t('dashboard.enterSessionName'))
    return
  }

  try {
    await sessionStore.createSession(createForm.name, createForm.workDir || undefined)
    message.success(t('dashboard.sessionCreated'))
    createForm.name = ''
    createForm.workDir = ''
  } catch (error: any) {
    const errorMsg = error?.response?.data?.error || error?.message || t('dashboard.createFailed')
    message.error(errorMsg)
  }
}

function handleOpenSession(session: any) {
  sessionStore.setCurrentSession(session)
  router.push(`/session/${session.name}`)
}

async function handleDeleteSession(name: string) {
  try {
    await sessionStore.deleteSession(name)
    message.success(t('dashboard.sessionDeleted'))
  } catch {
    message.error(t('dashboard.deleteFailed'))
  }
}

function handleRefresh() {
  sessionStore.fetchSessions()
  message.success(t('dashboard.refreshSuccess'))
}

function handleLogout() {
  localStorage.removeItem('token')
  router.push('/login')
}

function formatTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleString()
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
}

.header h1 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.create-card {
  margin-bottom: 24px;
}

.work-dir {
  color: #666;
  font-size: 13px;
}

.created-at {
  color: #999;
  font-size: 12px;
}
</style>
