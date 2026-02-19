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
  <div class="session-sidebar">
    <!-- Header -->
    <div class="sidebar-header">
      <h3>{{ t('sidebar.sessions') }}</h3>
      <n-button quaternary size="small" @click="handleToggle">
        <template #icon>
          <n-icon><ChevronBackIcon /></n-icon>
        </template>
      </n-button>
    </div>

    <!-- Create Session Button -->
    <div class="sidebar-actions">
      <n-button type="primary" block @click="openCreateDialog">
        <template #icon>
          <n-icon><AddIcon /></n-icon>
        </template>
        {{ t('sidebar.newSession') }}
      </n-button>
    </div>

    <!-- Session List -->
    <div class="session-list">
      <n-spin :show="sessionStore.loading">
        <n-empty
          v-if="sessionStore.sessions.length === 0"
          :description="t('sidebar.noSessions')"
          size="small"
        />
        <n-menu
          v-else
          :value="currentSessionKey"
          :options="menuOptions"
          @update:value="handleSelectSession"
        />
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  NButton,
  NIcon,
  NSpin,
  NEmpty,
  NMenu,
  NForm,
  NFormItem,
  NInput,
  NTag,
  useMessage,
  useDialog,
  type MenuOption
} from 'naive-ui'
import {
  Add as AddIcon,
  ChevronBack as ChevronBackIcon,
  Terminal as TerminalIcon,
  Ellipse as StopIcon,
  TrashOutline as TrashIcon,
  SettingsOutline as SettingsIcon
} from '@vicons/ionicons5'
import { useSessionStore } from '@/stores/session'
import type { Session } from '@/api/client'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const message = useMessage()
const dialog = useDialog()
const sessionStore = useSessionStore()

// Emits
const emit = defineEmits<{
  toggle: []
}>()

// State
const creating = ref(false)
const createForm = ref({
  name: '',
  workDir: ''
})

// Edit session state
const editForm = ref({
  name: '',
  workDir: ''
})

// Current session key for menu selection
const currentSessionKey = computed(() => {
  return sessionStore.currentSession?.name || ''
})

// Convert sessions to menu options
const menuOptions = computed<MenuOption[]>(() => {
  return sessionStore.sessions.map((session: Session) => ({
    key: session.name,
    label: () =>
      h('div', { class: 'session-label' }, [
        h('div', { class: 'session-name' }, session.name)
      ]),
    icon: () =>
      h(NIcon, null, {
        default: () => h(session.is_active ? TerminalIcon : StopIcon)
      }),
    extra: () =>
      h(
        'div',
        { class: 'session-extra' },
        [
          h(
            NTag,
            {
              type: session.is_active ? 'success' : 'default',
              size: 'small',
              bordered: false
            },
            { default: () => (session.is_active ? t('dashboard.active') : t('dashboard.inactive')) }
          ),
          h(
            NButton,
            {
              text: true,
              type: 'default',
              size: 'small',
              onClick: (e: Event) => {
                e.stopPropagation()
                openEditDialog(session)
              }
            },
            { icon: () => h(NIcon, null, { default: () => h(SettingsIcon) }) }
          ),
          h(
            NButton,
            {
              text: true,
              type: 'error',
              size: 'small',
              onClick: (e: Event) => {
                e.stopPropagation()
                handleDeleteSession(session.name)
              }
            },
            { icon: () => h(NIcon, null, { default: () => h(TrashIcon) }) }
          )
        ]
      )
  }))
})

/**
 * Toggle sidebar
 */
function handleToggle() {
  emit('toggle')
}

/**
 * Select a session
 */
function handleSelectSession(key: string) {
  const session = sessionStore.selectSession(key)
  if (session) {
    router.push(`/app/sessions/${session.name}`)
  }
}

/**
 * Create new session
 */
async function handleCreateSession() {
  if (!createForm.value.name) {
    message.error(t('dashboard.enterSessionName'))
    return false
  }

  creating.value = true
  try {
    const session = await sessionStore.createSession(
      createForm.value.name,
      createForm.value.workDir || undefined
    )
    message.success(t('dashboard.sessionCreated'))
    createForm.value = { name: '', workDir: '' }
    // Navigate to the new session
    sessionStore.setCurrentSession(session)
    router.push(`/app/sessions/${session.name}`)
  } catch (error: any) {
    const errorMsg = error?.response?.data?.error || error?.message || t('dashboard.createFailed')
    message.error(errorMsg)
    return false
  } finally {
    creating.value = false
  }
  return true
}

/**
 * Open create session dialog
 */
function openCreateDialog() {
  console.log('openCreateDialog called')
  createForm.value = { name: '', workDir: '' }
  dialog.create({
    title: t('sidebar.createSession'),
    positiveText: t('common.create'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      if (!createForm.value.name) {
        message.error(t('dashboard.enterSessionName'))
        return false
      }
      handleCreateSession()
    },
    content: () => h('div', { style: 'padding: 16px 0;' }, [
      h('div', { style: 'margin-bottom: 16px;' }, [
        h('label', { style: 'display: block; margin-bottom: 8px; color: #A9B7C6; font-size: 13px;' }, t('sidebar.sessionName') + ' *'),
        h(NInput, {
          value: createForm.value.name,
          onUpdateValue: (v: string) => { createForm.value.name = v },
          placeholder: t('sidebar.sessionNamePlaceholder')
        })
      ]),
      h('div', {}, [
        h('label', { style: 'display: block; margin-bottom: 8px; color: #A9B7C6; font-size: 13px;' }, t('sidebar.workDirectory')),
        h(NInput, {
          value: createForm.value.workDir,
          onUpdateValue: (v: string) => { createForm.value.workDir = v },
          placeholder: t('sidebar.workDirectoryPlaceholder')
        })
      ])
    ])
  })
}

/**
 * Open edit session dialog
 */
function openEditDialog(session: Session) {
  editForm.value = {
    name: session.name,
    workDir: session.work_dir || ''
  }

  dialog.create({
    title: t('sidebar.sessionConfig'),
    positiveText: t('common.save'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      handleUpdateSession()
    },
    content: () => h('div', { style: 'padding: 16px 0;' }, [
      h('div', { style: 'margin-bottom: 16px;' }, [
        h('label', { style: 'display: block; margin-bottom: 8px; color: #A9B7C6; font-size: 13px;' }, t('sidebar.sessionName')),
        h(NInput, {
          value: editForm.value.name,
          disabled: true,
          placeholder: t('sidebar.sessionNamePlaceholder')
        }),
        h('span', { style: 'font-size: 11px; color: #808080; margin-top: 4px; display: block;' }, t('sidebar.sessionNameReadOnly'))
      ]),
      h('div', {}, [
        h('label', { style: 'display: block; margin-bottom: 8px; color: #A9B7C6; font-size: 13px;' }, t('sidebar.workDirectory')),
        h(NInput, {
          value: editForm.value.workDir,
          onUpdateValue: (v: string) => { editForm.value.workDir = v },
          placeholder: t('sidebar.workDirectoryPlaceholder')
        }),
        h('span', { style: 'font-size: 11px; color: #808080; margin-top: 4px; display: block;' }, t('sidebar.workDirectoryNote'))
      ])
    ])
  })
}

/**
 * Update session (currently only work_dir display, actual change requires session restart)
 */
function handleUpdateSession() {
  // For now, just show a message that this is display-only
  // In the future, we could implement session restart with new work_dir
  message.info(t('sidebar.sessionInfoUpdated'))
}

/**
 * Delete a session
 */
async function handleDeleteSession(name: string) {
  try {
    await sessionStore.deleteSession(name)
    message.success(t('dashboard.sessionDeleted'))
    // If we're deleting the current session, navigate away
    if (route.params.name === name) {
      router.push('/app/sessions')
    }
  } catch {
    message.error(t('dashboard.deleteFailed'))
  }
}

// Load sessions on mount
onMounted(() => {
  sessionStore.fetchSessions()
})
</script>

<style scoped>
.session-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #2B2B2B;
  border-right: 1px solid #4E5052;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #4E5052;
  background: #3C3F41;
}

.sidebar-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #D4D4D4;
}

.sidebar-actions {
  padding: 12px 16px;
  border-bottom: 1px solid #4E5052;
  background: #3C3F41;
}

.session-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
  background: #2B2B2B;
}

.session-extra {
  display: flex;
  align-items: center;
  gap: 4px;
}

.session-label {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.session-name {
  font-weight: 500;
  color: #A9B7C6;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Force visible text colors for menu items */
:deep(.n-menu-item) {
  padding-right: 8px;
}

:deep(.n-menu-item-content) {
  padding: 0 12px;
}

:deep(.n-menu-item-content-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .session-sidebar {
    border-right: none;
    border-bottom: 1px solid #4E5052;
  }

  .sidebar-header {
    padding: 12px;
  }

  .sidebar-header h3 {
    font-size: 13px;
  }

  .sidebar-actions {
    padding: 10px 12px;
  }

  .session-list {
    padding: 6px 0;
    max-height: 120px;
  }

  :deep(.n-menu-item) {
    padding: 10px 12px;
    min-height: 44px;
  }

  :deep(.n-button) {
    min-height: 36px;
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .sidebar-header {
    padding: 10px;
  }

  .sidebar-header h3 {
    font-size: 12px;
  }

  .sidebar-actions {
    padding: 8px 10px;
  }

  :deep(.n-menu-item) {
    padding: 8px 10px;
    min-height: 40px;
  }

  :deep(.n-button) {
    min-height: 32px;
    font-size: 12px;
    padding: 0 10px;
  }

  :deep(.n-tag) {
    font-size: 11px;
    padding: 2px 6px;
  }

  .session-list {
    max-height: 100px;
  }
}
</style>
