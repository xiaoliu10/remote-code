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
  <div class="session-container">
    <!-- 头部 -->
    <header class="header">
      <div class="header-left">
        <n-button text @click="handleBack">
          <template #icon>
            <n-icon><ArrowBackIcon /></n-icon>
          </template>
        </n-button>
        <h2>{{ sessionName }}</h2>
        <n-tag :type="connected ? 'success' : 'error'" size="small">
          {{ connected ? t('common.connected') : t('common.disconnected') }}
        </n-tag>
      </div>
      <div class="header-right">
        <n-input
          v-model:value="command"
          :placeholder="t('session.typeCommand')"
          @keydown.enter="handleSendCommand"
          style="width: 400px"
        >
          <template #prefix>
            <n-icon><TerminalIcon /></n-icon>
          </template>
        </n-input>
        <n-button type="primary" :disabled="!connected" @click="handleSendCommand">
          {{ t('common.send') }}
        </n-button>
      </div>
    </header>

    <!-- 终端 -->
    <div class="terminal-wrapper">
      <div ref="terminalRef" class="terminal" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import {
  NButton,
  NIcon,
  NInput,
  NTag,
  useMessage,
  useDialog
} from 'naive-ui'
import { ArrowBack as ArrowBackIcon, Terminal as TerminalIcon } from '@vicons/ionicons5'
import { useWebSocket } from '@/composables/useWebSocket'
import { useSessionStore } from '@/stores/session'

const { t } = useI18n()
const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const sessionStore = useSessionStore()

const props = defineProps<{
  name: string
}>()

const sessionName = ref(props.name)
const terminalRef = ref<HTMLElement>()
const command = ref('')

let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null

// WebSocket 连接
const { connected, connect, disconnect, sendCommand, onMessage } = useWebSocket(sessionName.value)

onMounted(() => {
  initTerminal()
  connect()
  setupWebSocketHandlers()
})

onUnmounted(() => {
  if (terminal) {
    terminal.dispose()
  }
  disconnect()
})

function initTerminal() {
  if (!terminalRef.value) return

  terminal = new Terminal({
    cursorBlink: true,
    fontSize: 14,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#d4d4d4',
      selection: '#264f78'
    },
    scrollback: 1000
  })

  fitAddon = new FitAddon()
  terminal.loadAddon(fitAddon)

  terminal.open(terminalRef.value)
  fitAddon.fit()

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    fitAddon?.fit()
  })

  // 欢迎消息
  terminal.writeln(t('session.welcomeMessage'))
  terminal.writeln(`${t('session.sessionLabel')}: ${sessionName.value}`)
  terminal.writeln(t('session.connecting'))
  terminal.writeln('')
}

function setupWebSocketHandlers() {
  onMessage((msg) => {
    if (!terminal) return

    switch (msg.type) {
      case 'output':
        const data = msg.data as { text: string; timestamp: number }
        // 清屏并重新输出
        terminal.clear()
        terminal.write(data.text.replace(/\n/g, '\r\n'))
        break
      case 'error':
        message.error(String(msg.data))
        break
      case 'status':
        // 状态更新
        break
      case 'pong':
        // 心跳响应
        break
    }
  })
}

function handleSendCommand() {
  if (!command.value.trim() || !connected.value) return

  sendCommand(command.value)
  command.value = ''
}

function handleBack() {
  dialog.warning({
    title: t('session.leaveSession'),
    content: t('session.leaveConfirm'),
    positiveText: t('session.leave'),
    negativeText: t('session.stay'),
    onPositiveClick: () => {
      router.push('/dashboard')
    }
  })
}
</script>

<style scoped>
.session-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1e1e1e;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #3e3e3e;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header h2 {
  margin: 0;
  color: #d4d4d4;
  font-size: 16px;
  font-weight: 600;
}

.terminal-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 8px;
}

.terminal {
  height: 100%;
}

:deep(.n-input) {
  background: #3e3e3e;
}

:deep(.n-input__input) {
  color: #d4d4d4;
}
</style>
