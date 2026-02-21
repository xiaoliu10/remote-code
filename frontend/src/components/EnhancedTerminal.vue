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
  <div class="enhanced-terminal">
    <!-- Terminal Header -->
    <div class="terminal-header">
      <div class="header-left">
        <h3>{{ sessionName || t('terminal.title') }}</h3>
        <n-tag :type="connected ? 'success' : 'error'" size="small">
          {{ connected ? t('common.connected') : t('common.disconnected') }}
        </n-tag>
        <span v-if="workDir" class="work-dir">{{ workDir }}</span>
      </div>
      <div class="header-right">
        <n-button quaternary size="small" @click="handleClear">
          <template #icon>
            <n-icon><TrashIcon /></n-icon>
          </template>
        </n-button>
        <n-button quaternary size="small" @click="handleReconnect" :disabled="connected || kicked">
          <template #icon>
            <n-icon><RefreshIcon /></n-icon>
          </template>
        </n-button>
        <n-divider vertical />
        <!-- Scroll mode toggle -->
        <n-button
          quaternary
          size="small"
          @click="scrollMode = scrollMode === 'local' ? 'remote' : 'local'"
          :type="scrollMode === 'remote' ? 'primary' : 'default'"
          :title="scrollMode === 'local' ? t('terminal.localScroll') : t('terminal.remoteScroll')"
        >
          <template #icon>
            <n-icon><SwapVerticalIcon /></n-icon>
          </template>
        </n-button>
        <!-- Scroll controls -->
        <n-button-group size="small">
          <n-button quaternary size="small" @click="scrollToTop" title="Scroll to top">
            <template #icon>
              <n-icon><ArrowUpIcon /></n-icon>
            </template>
          </n-button>
          <n-button quaternary size="small" @click="scrollPageUp" title="Page up">
            <template #icon>
              <n-icon><ChevronUpIcon /></n-icon>
            </template>
          </n-button>
          <n-button quaternary size="small" @click="scrollPageDown" title="Page down">
            <template #icon>
              <n-icon><ChevronDownIcon /></n-icon>
            </template>
          </n-button>
          <n-button quaternary size="small" @click="scrollToBottom" title="Scroll to bottom">
            <template #icon>
              <n-icon><ArrowDownIcon /></n-icon>
            </template>
          </n-button>
        </n-button-group>
      </div>
    </div>

    <!-- Terminal Container -->
    <div ref="terminalContainer" class="terminal-container" />

    <!-- File Reference Popover (positioned over terminal) -->
    <div v-if="showFilePopover" class="file-popover-overlay">
      <div class="file-popover">
        <div class="file-popover-header">
          <span>{{ t('terminal.selectFile') }}</span>
          <n-button quaternary size="tiny" @click="showFilePopover = false">
            <template #icon>
              <n-icon><CloseIcon /></n-icon>
            </template>
          </n-button>
        </div>
        <n-input
          v-model:value="fileSearchQuery"
          :placeholder="t('terminal.searchFiles')"
          size="small"
          style="margin-bottom: 8px;"
        />
        <div class="file-reference-list">
          <n-spin :show="loadingFiles" size="small">
            <div v-if="filteredFiles.length === 0 && !loadingFiles" class="no-files">
              {{ t('terminal.noFilesFound') }}
            </div>
            <div
              v-for="file in filteredFiles"
              :key="file.path"
              class="file-item"
              @click="insertFileReference(file)"
            >
              <n-icon :color="file.type === 'directory' ? '#4A9CFF' : '#A9B7C6'">
                <FolderIcon v-if="file.type === 'directory'" />
                <DocumentIcon v-else />
              </n-icon>
              <span class="file-name">{{ file.name }}</span>
              <span class="file-path">{{ file.path }}</span>
            </div>
          </n-spin>
        </div>
      </div>
    </div>

    <!-- Command Input -->
    <div class="terminal-input-wrapper">
      <!-- Custom Keyboard -->
      <div v-if="showCustomKeyboard" class="custom-keyboard">
        <!-- Main Keys -->
        <div class="keyboard-main">
          <!-- Left Side: Common Keys -->
          <div class="keyboard-left">
            <div class="keyboard-row">
              <button class="key-btn" @click="handleSpecialKey('at')">@</button>
              <button class="key-btn" @click="handleSpecialKey('slash')">/</button>
              <button class="key-btn" @click="handleSpecialKey('tab')">Tab</button>
            </div>
            <div class="keyboard-row">
              <button class="key-btn" @click="handleSpecialKey('home')">Home</button>
              <button class="key-btn" @click="handleSpecialKey('end')">End</button>
              <button class="key-btn" @click="handleSpecialKey('backspace')">⌫</button>
            </div>
            <div class="keyboard-row">
              <button class="key-btn danger" @click="handleSpecialKey('c-c')">Ctrl+C</button>
              <button class="key-btn" @click="handleSpecialKey('c-l')">Ctrl+L</button>
              <button class="key-btn" @click="handleSpecialKey('c-d')">Ctrl+D</button>
            </div>
          </div>

          <!-- Right Side: Arrow Keys in Cross Layout -->
          <div class="keyboard-arrows">
            <div class="arrow-spacer"></div>
            <div class="arrow-row">
              <div class="arrow-placeholder"></div>
              <button class="key-btn arrow" @click="handleSpecialKey('arrow-up')">↑</button>
              <div class="arrow-placeholder"></div>
            </div>
            <div class="arrow-row">
              <button class="key-btn arrow" @click="handleSpecialKey('arrow-left')">←</button>
              <button class="key-btn arrow" @click="handleSpecialKey('arrow-down')">↓</button>
              <button class="key-btn arrow" @click="handleSpecialKey('arrow-right')">→</button>
            </div>
          </div>
        </div>

        <!-- Bottom Row: ESC and Enter Keys -->
        <div class="keyboard-row bottom-row">
          <button class="key-btn" @click="handleSpecialKey('escape')">ESC</button>
          <button class="key-btn primary" @click="handleSpecialKey('enter')">Enter</button>
        </div>
      </div>

      <!-- Input Row -->
      <div class="terminal-input">
        <!-- Keyboard Toggle -->
        <n-button
          quaternary
          size="small"
          :type="showCustomKeyboard ? 'primary' : 'default'"
          @click="showCustomKeyboard = !showCustomKeyboard"
        >
          <template #icon>
            <n-icon><KeyboardIcon /></n-icon>
          </template>
        </n-button>

        <!-- Voice Input Button -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              quaternary
              size="small"
              :type="isListening ? 'error' : 'default'"
              :disabled="!voiceSupported || !connected"
              @click="toggleVoiceInput"
            >
              <template #icon>
                <n-icon :class="{ 'mic-active': isListening }">
                  <MicIcon />
                </n-icon>
              </template>
            </n-button>
          </template>
          {{ isListening ? t('terminal.listening') : (voiceSupported ? t('terminal.voiceInput') : t('terminal.voiceNotSupported')) }}
        </n-tooltip>

        <!-- Mode Toggle -->
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button
              quaternary
              size="small"
              :type="realtimeMode ? 'primary' : 'default'"
              @click="realtimeMode = !realtimeMode"
            >
              <template #icon>
                <n-icon><BoltIcon /></n-icon>
              </template>
            </n-button>
          </template>
          {{ realtimeMode ? t('terminal.realtimeMode') : t('terminal.commandMode') }}
        </n-tooltip>

        <n-input
          ref="commandInputRef"
          v-model:value="currentCommand"
          :placeholder="realtimeMode ? t('terminal.typeCommandRealtime') : t('terminal.typeCommand')"
          :disabled="!connected"
          @keydown="handleKeyDown"
          @update:value="handleInputChange"
          @blur="handleInputBlur"
        >
          <template #prefix>
            <n-icon><TerminalIcon /></n-icon>
          </template>
        </n-input>
        <n-button type="primary" :disabled="!connected || !currentCommand.trim()" @click="sendCurrentCommand">
          {{ t('common.send') }}
        </n-button>
      </div>
    </div>

    <!-- Connection Overlay -->
    <div v-if="!connected && !error" class="connection-overlay">
      <n-spin size="large" />
      <p>{{ t('common.connecting') }}</p>
    </div>

    <!-- Error Overlay -->
    <div v-if="error" class="error-overlay">
      <n-icon size="48" color="var(--n-error-color)">
        <WarningIcon />
      </n-icon>
      <p>{{ error }}</p>
      <n-button type="primary" @click="handleReconnect">
        {{ t('terminal.reconnect') }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { SearchAddon } from '@xterm/addon-search'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'
import {
  NTag,
  NButton,
  NIcon,
  NInput,
  NSpin,
  NDropdown,
  NPopover,
  useMessage
} from 'naive-ui'
import type { DropdownOption } from 'naive-ui'
import {
  Terminal as TerminalIcon,
  TrashOutline as TrashIcon,
  Refresh as RefreshIcon,
  Warning as WarningIcon,
  GameController as KeyboardIcon,
  FolderOutline as FolderIcon,
  DocumentOutline as DocumentIcon,
  Close as CloseIcon,
  Flash as BoltIcon,
  Mic as MicIcon,
  ArrowUp as ArrowUpIcon,
  ArrowDown as ArrowDownIcon,
  ChevronUp as ChevronUpIcon,
  ChevronDown as ChevronDownIcon,
  SwapVertical as SwapVerticalIcon
} from '@vicons/ionicons5'
import { useWebSocket } from '@/composables/useWebSocket'
import { useSessionStore } from '@/stores/session'
import * as api from '@/api/client'
import type { FileItem } from '@/api/client'

const { t } = useI18n()
const message = useMessage()
const sessionStore = useSessionStore()

// Props
const props = defineProps<{
  sessionName: string
}>()

// Emits
const emit = defineEmits<{
  connected: []
  disconnected: []
  error: [message: string]
}>()

// Computed
const workDir = computed(() => {
  const session = sessionStore.getSessionByName(props.sessionName)
  return session?.work_dir || null
})

// Refs
const terminalContainer = ref<HTMLElement>()
const commandInputRef = ref<InstanceType<typeof NInput>>()
const currentCommand = ref('')
const commandHistory = ref<string[]>([])
const historyIndex = ref(-1)

// File reference state
const showFilePopover = ref(false)
const loadingFiles = ref(false)
const fileList = ref<FileItem[]>([])
const fileSearchQuery = ref('')
const atPosition = ref(-1)

// Input mode state
const realtimeMode = ref(true) // true = realtime input, false = command mode
const lastSentLength = ref(0) // Track last sent position for realtime mode
const showCustomKeyboard = ref(false) // Show custom keyboard on mobile

// Voice input state
const isListening = ref(false)
const voiceSupported = ref(false)
let recognition: SpeechRecognition | null = null
let voiceTimeout: ReturnType<typeof setTimeout> | null = null

// Special key options
const specialKeyOptions: DropdownOption[] = [
  { label: '@ (文件引用)', key: 'at' },
  { label: '/ (路径)', key: 'slash' },
  { type: 'divider', key: 'd0' },
  { label: '↑ 上', key: 'arrow-up' },
  { label: '↓ 下', key: 'arrow-down' },
  { label: '← 左', key: 'arrow-left' },
  { label: '→ 右', key: 'arrow-right' },
  { type: 'divider', key: 'd-arrows' },
  { label: 'Home', key: 'home' },
  { label: 'End', key: 'end' },
  { label: 'Page Up', key: 'pageup' },
  { label: 'Page Down', key: 'pagedown' },
  { type: 'divider', key: 'd-nav' },
  { label: 'ESC', key: 'escape' },
  { label: 'ESC ESC (退出Claude)', key: 'escape-escape' },
  { label: 'Ctrl+C (中断)', key: 'c-c' },
  { label: 'Ctrl+D (退出)', key: 'c-d' },
  { label: 'Ctrl+Z (暂停)', key: 'c-z' },
  { label: 'Ctrl+L (清屏)', key: 'c-l' },
  { type: 'divider', key: 'd1' },
  { label: 'Enter', key: 'enter' },
  { label: 'Tab', key: 'tab' },
  { label: 'Backspace', key: 'backspace' },
  { label: 'Delete', key: 'delete' }
]

// Computed filtered files
const filteredFiles = computed(() => {
  if (!fileSearchQuery.value) {
    return fileList.value.slice(0, 20) // Limit to 20 items
  }
  const query = fileSearchQuery.value.toLowerCase()
  return fileList.value
    .filter(f => f.name.toLowerCase().includes(query))
    .slice(0, 20)
})

// Terminal instance
let terminal: Terminal | null = null
let fitAddon: FitAddon | null = null
let searchAddon: SearchAddon | null = null

// WebSocket connection
const { connected, error, kicked, connect, disconnect, sendCommand, sendKeys, onMessage, onConnect, onDisconnect } =
  useWebSocket(props.sessionName)

// Scroll mode: 'local' (scroll web view) or 'remote' (send to terminal)
const scrollMode = ref<'local' | 'remote'>('local')

/**
 * Handle special key selection
 */
function handleSpecialKey(key: string) {
  if (!connected.value) {
    message.warning(t('terminal.notConnected'))
    return
  }

  // Handle @ and / separately to trigger file reference
  if (key === 'at') {
    sendKeys('@')
    if (terminal) {
      terminal.write('@')
    }
    return
  }

  if (key === 'slash') {
    sendKeys('/')
    if (terminal) {
      terminal.write('/')
    }
    return
  }

  const specialKeys: Record<string, string> = {
    'escape': '\x1b',
    'escape-escape': '\x1b\x1b',
    'c-c': '\x03',
    'c-d': '\x04',
    'c-z': '\x1a',
    'c-l': '\x0c',
    'enter': '\r',
    'tab': '\t',
    'backspace': '\x7f',
    'delete': '\x1b[3~',
    // Arrow keys
    'arrow-up': '\x1b[A',
    'arrow-down': '\x1b[B',
    'arrow-right': '\x1b[C',
    'arrow-left': '\x1b[D',
    // Navigation keys
    'home': '\x1b[H',
    'end': '\x1b[F',
    'pageup': '\x1b[5~',
    'pagedown': '\x1b[6~'
  }

  const keyCode = specialKeys[key]
  if (keyCode) {
    // Use sendKeys for special keys to avoid auto-adding Enter
    sendKeys(keyCode)
  }
}

/**
 * Initialize voice recognition
 */
function initVoiceRecognition() {
  // Check for browser support
  const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition

  if (!SpeechRecognition) {
    voiceSupported.value = false
    return
  }

  voiceSupported.value = true
  recognition = new SpeechRecognition()

  // Configure recognition
  recognition.continuous = false
  recognition.interimResults = true
  recognition.lang = 'zh-CN' // Default to Chinese, can be changed

  recognition.onstart = () => {
    isListening.value = true
  }

  recognition.onend = () => {
    isListening.value = false
    if (voiceTimeout) {
      clearTimeout(voiceTimeout)
      voiceTimeout = null
    }
  }

  recognition.onerror = (event: SpeechRecognitionErrorEvent) => {
    console.error('Speech recognition error:', event.error)
    isListening.value = false

    if (event.error === 'not-allowed') {
      message.error(t('terminal.voicePermissionDenied'))
    } else if (event.error === 'no-speech') {
      message.warning(t('terminal.noSpeechDetected'))
    } else {
      message.error(t('terminal.voiceError') + ': ' + event.error)
    }
  }

  recognition.onresult = (event: SpeechRecognitionEvent) => {
    let interimTranscript = ''
    let finalTranscript = ''

    for (let i = event.resultIndex; i < event.results.length; i++) {
      const transcript = event.results[i][0].transcript
      if (event.results[i].isFinal) {
        finalTranscript += transcript
      } else {
        interimTranscript += transcript
      }
    }

    // Update command with recognized text
    if (finalTranscript) {
      // Handle final result
      const text = finalTranscript.trim()
      if (realtimeMode.value) {
        // In realtime mode, send each character
        sendKeys(text)
        if (terminal) {
          terminal.write(text)
        }
      } else {
        // In command mode, append to current command
        currentCommand.value += text
      }

      // Auto-send after a short delay if in command mode and text ends with certain patterns
      if (!realtimeMode.value) {
        // Don't auto-send, let user review and send manually
        message.success(t('terminal.voiceRecognized') + ': ' + text)
      }
    }
  }
}

/**
 * Toggle voice input
 */
function toggleVoiceInput() {
  if (!recognition) {
    initVoiceRecognition()
    if (!recognition) {
      message.error(t('terminal.voiceNotSupported'))
      return
    }
  }

  if (isListening.value) {
    recognition.stop()
    isListening.value = false
  } else {
    try {
      recognition.start()

      // Auto-stop after 10 seconds
      voiceTimeout = setTimeout(() => {
        if (recognition && isListening.value) {
          recognition.stop()
          message.info(t('terminal.voiceTimeout'))
        }
      }, 10000)
    } catch (e) {
      console.error('Failed to start recognition:', e)
      message.error(t('terminal.voiceStartFailed'))
    }
  }
}

/**
 * Initialize terminal
 */
function initTerminal() {
  if (!terminalContainer.value) return

  // Detect mobile for font size
  const isMobile = window.innerWidth < 768
  const fontSize = isMobile ? 12 : 14

  // Create terminal instance
  terminal = new Terminal({
    cursorBlink: true,
    fontSize: fontSize,
    fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    theme: {
      background: '#1e1e1e',
      foreground: '#d4d4d4',
      cursor: '#d4d4d4',
      cursorAccent: '#1e1e1e',
      selection: '#264f78',
      black: '#000000',
      red: '#cd3131',
      green: '#0dbc79',
      yellow: '#e5e510',
      blue: '#2472c8',
      magenta: '#bc3fbc',
      cyan: '#11a8cd',
      white: '#e5e5e5',
      brightBlack: '#666666',
      brightRed: '#f14c4c',
      brightGreen: '#23d18b',
      brightYellow: '#f5f543',
      brightBlue: '#3b8eea',
      brightMagenta: '#d670d6',
      brightCyan: '#29b8db',
      brightWhite: '#e5e5e5'
    },
    scrollback: 10000,
    allowProposedApi: true
  })

  // Load addons
  fitAddon = new FitAddon()
  searchAddon = new SearchAddon()
  const webLinksAddon = new WebLinksAddon()

  terminal.loadAddon(fitAddon)
  terminal.loadAddon(searchAddon)
  terminal.loadAddon(webLinksAddon)

  // Open terminal
  terminal.open(terminalContainer.value)

  // Handle mouse wheel - send to remote terminal when in remote mode
  let lastScrollTime = 0
  const scrollThrottle = 100 // ms between scroll events

  terminalContainer.value.addEventListener('wheel', (e: WheelEvent) => {
    if (scrollMode.value !== 'remote') {
      // Local mode: let the browser handle scrolling
      return
    }

    // Remote mode: prevent default and send to terminal
    e.preventDefault()

    if (!connected.value) return

    // Throttle scroll events
    const now = Date.now()
    if (now - lastScrollTime < scrollThrottle) return
    lastScrollTime = now

    if (e.deltaY < 0) {
      // Scroll up - send PageUp
      sendKeys('\x1b[5~') // PageUp escape sequence
    } else {
      // Scroll down - send PageDown
      sendKeys('\x1b[6~') // PageDown escape sequence
    }
  }, { passive: false })

  // Fit to container
  nextTick(() => {
    fitAddon?.fit()
  })

  // Welcome message
  terminal.writeln(t('terminal.welcome'))
  terminal.writeln(`${t('terminal.session')}: ${props.sessionName}`)
  terminal.writeln('')

  // Handle terminal resize
  const resizeObserver = new ResizeObserver(() => {
    fitAddon?.fit()
  })
  resizeObserver.observe(terminalContainer.value)

  // Handle mobile keyboard show/hide - refit terminal when keyboard closes
  const handleViewportResize = () => {
    // Delay to allow DOM to settle after keyboard animation
    setTimeout(() => {
      fitAddon?.fit()
    }, 100)
  }

  // Use visualViewport API for better mobile keyboard handling
  if (window.visualViewport) {
    window.visualViewport.addEventListener('resize', handleViewportResize)
  }

  // Also handle window resize as fallback
  window.addEventListener('resize', handleViewportResize)

  // Store cleanup function
  const cleanupViewportListeners = () => {
    if (window.visualViewport) {
      window.visualViewport.removeEventListener('resize', handleViewportResize)
    }
    window.removeEventListener('resize', handleViewportResize)
    resizeObserver.disconnect()
  }

  // Attach to terminal instance for cleanup
  if (terminal) {
    (terminal as any)._cleanupViewportListeners = cleanupViewportListeners
  }

  // Note: terminal.onData is disabled - we use command input instead
  // If you want direct terminal input, uncomment the following:
  // terminal.onData((data) => {
  //   if (connected.value) {
  //     sendCommand(data)
  //   }
  // })
}

/**
 * Handle keyboard events in command input
 */
function handleKeyDown(e: KeyboardEvent) {
  switch (e.key) {
    case 'Enter':
      sendCurrentCommand()
      break
    case 'ArrowUp':
      e.preventDefault()
      navigateHistory(-1)
      break
    case 'ArrowDown':
      e.preventDefault()
      navigateHistory(1)
      break
    case 'c':
      if (e.ctrlKey) {
        // Clear current command
        currentCommand.value = ''
      }
      break
    case 'l':
      if (e.ctrlKey) {
        // Clear terminal
        e.preventDefault()
        handleClear()
      }
      break
  }
}

/**
 * Navigate command history
 */
function navigateHistory(direction: number) {
  if (commandHistory.value.length === 0) return

  const newIndex = historyIndex.value + direction
  if (newIndex >= -1 && newIndex < commandHistory.value.length) {
    historyIndex.value = newIndex
    if (newIndex === -1) {
      currentCommand.value = ''
    } else {
      currentCommand.value = commandHistory.value[commandHistory.value.length - 1 - newIndex]
    }
  }
}

/**
 * Handle input change - detect @ for file reference and realtime mode
 */
function handleInputChange(value: string) {
  // Realtime mode: send each character to terminal
  if (realtimeMode.value && connected.value) {
    // Calculate the difference and send only new characters
    const newLength = value.length
    if (newLength > lastSentLength.value) {
      // Send new characters
      const newChars = value.slice(lastSentLength.value)
      sendKeys(newChars)
    } else if (newLength < lastSentLength.value) {
      // Characters were deleted - send backspace for each deleted char
      const deletedCount = lastSentLength.value - newLength
      for (let i = 0; i < deletedCount; i++) {
        sendKeys('\x7f') // Backspace
      }
    }
    lastSentLength.value = newLength
  }

  // Detect @ for file reference
  const lastAtIndex = value.lastIndexOf('@')
  if (lastAtIndex !== -1) {
    // Check if @ is followed by space or at the end
    const afterAt = value.slice(lastAtIndex + 1)
    if (!afterAt.includes(' ')) {
      atPosition.value = lastAtIndex
      fileSearchQuery.value = afterAt
      showFilePopover.value = true
      loadFiles()
      return
    }
  }
  showFilePopover.value = false
}

/**
 * Handle input blur - refit terminal when mobile keyboard closes
 */
function handleInputBlur() {
  // Delay to allow keyboard close animation to complete
  setTimeout(() => {
    fitAddon?.fit()
  }, 300)
}

/**
 * Load files for reference
 */
async function loadFiles() {
  if (fileList.value.length > 0) return // Already loaded

  loadingFiles.value = true
  try {
    const basePath = workDir.value || '/'
    const response = await api.fileApi.list(basePath)
    fileList.value = response.data.items || []
  } catch (e) {
    console.error('Failed to load files:', e)
  } finally {
    loadingFiles.value = false
  }
}

/**
 * Insert file reference into command
 */
function insertFileReference(file: FileItem) {
  if (atPosition.value === -1) return

  const relativePath = workDir.value
    ? file.path.replace(workDir.value, '').replace(/^\//, '') || file.name
    : file.path

  // Replace from @ to cursor with file path
  const before = currentCommand.value.slice(0, atPosition.value)
  const after = currentCommand.value.slice(atPosition.value + fileSearchQuery.value.length + 1)

  currentCommand.value = before + relativePath + after
  showFilePopover.value = false
  atPosition.value = -1
  fileSearchQuery.value = ''
}

/**
 * Send current command (Enter key)
 */
function sendCurrentCommand() {
  const cmd = currentCommand.value.trim()
  if (!cmd || !connected.value) return

  // Add to history
  commandHistory.value.push(cmd)
  if (commandHistory.value.length > 100) {
    commandHistory.value.shift()
  }
  historyIndex.value = -1

  if (realtimeMode.value) {
    // In realtime mode, just send Enter to execute
    sendKeys('\r')
  } else {
    // In command mode, send the full command with Enter
    sendCommand(cmd + '\n')
  }

  // Clear input and reset tracking
  currentCommand.value = ''
  lastSentLength.value = 0
}

/**
 * Setup WebSocket handlers
 */
function setupWebSocketHandlers() {
  onConnect(() => {
    emit('connected')
    if (terminal) {
      terminal.writeln(t('terminal.connected'))
    }
  })

  onDisconnect(() => {
    emit('disconnected')
    if (terminal) {
      terminal.writeln('')
      terminal.writeln(t('terminal.disconnected'))
    }
  })

  onMessage((msg) => {
    if (!terminal) return

    switch (msg.type) {
      case 'output':
        const data = msg.data as { text: string; timestamp: number }
        // Clear and write output
        terminal.clear()
        terminal.write(data.text.replace(/\n/g, '\r\n'))
        break
      case 'error':
        const errorMsg = String(msg.data)
        terminal.writeln('')
        terminal.writeln(`\x1b[31m${t('terminal.error')}: ${errorMsg}\x1b[0m`)
        emit('error', errorMsg)
        break
      case 'kicked':
        // Connection was kicked by another device
        terminal.writeln('')
        terminal.writeln('\x1b[33m' + String(msg.data) + '\x1b[0m')
        terminal.writeln('\x1b[33mThis session is now being used from another device.\x1b[0m')
        terminal.writeln('\x1b[33mPlease refresh the page if you want to reconnect.\x1b[0m')
        message.warning('Connection replaced by another device')
        break
      case 'status':
        // Handle status updates
        break
      case 'pong':
        // Heartbeat response
        break
    }
  })
}

/**
 * Clear terminal
 */
function handleClear() {
  if (terminal) {
    terminal.clear()
  }
}

/**
 * Reconnect to session
 */
function handleReconnect() {
  // Don't reconnect if kicked
  if (kicked.value) {
    message.warning('This session is being used from another device. Please refresh the page.')
    return
  }

  disconnect()
  setTimeout(() => {
    connect()
  }, 500)
}

/**
 * Scroll terminal to top
 */
function scrollToTop() {
  if (terminal) {
    terminal.scrollToTop()
  }
}

/**
 * Scroll terminal to bottom
 */
function scrollToBottom() {
  if (terminal) {
    terminal.scrollToBottom()
  }
}

/**
 * Scroll terminal one page up
 */
function scrollPageUp() {
  if (terminal) {
    const rows = terminal.rows
    terminal.scrollLines(-rows)
  }
}

/**
 * Scroll terminal one page down
 */
function scrollPageDown() {
  if (terminal) {
    const rows = terminal.rows
    terminal.scrollLines(rows)
  }
}

/**
 * Cleanup terminal
 */
function cleanup() {
  if (terminal) {
    // Cleanup viewport listeners
    if ((terminal as any)._cleanupViewportListeners) {
      ;(terminal as any)._cleanupViewportListeners()
    }
    terminal.dispose()
    terminal = null
  }
  fitAddon = null
  searchAddon = null
  disconnect()
}

// Watch for session name changes
watch(
  () => props.sessionName,
  (newName, oldName) => {
    if (newName !== oldName) {
      cleanup()
      initTerminal()
      connect()
    }
  }
)

// Lifecycle
onMounted(() => {
  initTerminal()
  initVoiceRecognition()
  setupWebSocketHandlers()
  connect()
})

onUnmounted(() => {
  cleanup()
})
</script>

<style scoped>
.enhanced-terminal {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #2B2B2B;
  position: relative;
}

.terminal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #3C3F41;
  border-bottom: 1px solid #4E5052;
  min-height: 40px;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-left h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #D4D4D4;
}

.work-dir {
  font-size: 11px;
  color: #808080;
  font-family: 'Courier New', monospace;
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.terminal-container {
  flex: 1;
  padding: 8px;
  overflow: hidden;
  background: #1E1E1E;
  position: relative;
  display: flex;
  flex-direction: column;
}

/* Enable xterm.js scrollbar with prominent styling */
:deep(.xterm) {
  height: 100%;
  padding-right: 16px;
}

:deep(.xterm-screen) {
  padding-right: 8px;
}

:deep(.xterm-viewport) {
  overflow-y: scroll !important;
  scrollbar-width: auto;
  scrollbar-color: #4A9CFF #1E1E1E;
  /* Ensure viewport takes full height */
  height: 100% !important;
}

/* Webkit scrollbar (Chrome, Safari, Edge) */
:deep(.xterm-viewport::-webkit-scrollbar) {
  width: 16px;
  background: #1E1E1E;
}

:deep(.xterm-viewport::-webkit-scrollbar-track) {
  background: #2B2B2B;
  border-left: 2px solid #3C3F41;
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb) {
  background: linear-gradient(180deg, #4A9CFF 0%, #667eea 100%);
  border-radius: 0;
  border: 3px solid #2B2B2B;
  min-height: 60px;
  box-shadow: 0 0 4px rgba(74, 156, 255, 0.5);
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:hover) {
  background: linear-gradient(180deg, #5CB3FF 0%, #7D8FEE 100%);
  box-shadow: 0 0 8px rgba(74, 156, 255, 0.8);
}

:deep(.xterm-viewport::-webkit-scrollbar-thumb:active) {
  background: linear-gradient(180deg, #3D8CE8 0%, #5A73D6 100%);
}

/* Hide scrollbar buttons */
:deep(.xterm-viewport::-webkit-scrollbar-button:start:decrement),
:deep(.xterm-viewport::-webkit-scrollbar-button:end:increment) {
  display: block;
  height: 20px;
  background: #3C3F41;
  border: none;
}

:deep(.xterm-viewport::-webkit-scrollbar-button:start:decrement:after),
:deep(.xterm-viewport::-webkit-scrollbar-button:end:increment:after) {
  content: '';
  display: block;
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  margin: 3px auto;
}

:deep(.xterm-viewport::-webkit-scrollbar-button:start:decrement:after) {
  border-bottom: 6px solid #888;
}

:deep(.xterm-viewport::-webkit-scrollbar-button:end:increment:after) {
  border-top: 6px solid #888;
}

:deep(.xterm-viewport::-webkit-scrollbar-button:hover:after) {
  border-bottom-color: #4A9CFF;
  border-top-color: #4A9CFF;
}

.terminal-input {
  display: flex;
  gap: 8px;
  padding: 8px 16px;
  background: #3C3F41;
  border-top: 1px solid #4E5052;
}

/* Custom Keyboard Styles */
.terminal-input-wrapper {
  display: flex;
  flex-direction: column;
  background: #3C3F41;
}

.custom-keyboard {
  padding: 8px;
  background: #2B2B2B;
  border-top: 1px solid #4E5052;
}

.keyboard-main {
  display: flex;
  gap: 12px;
  margin-bottom: 6px;
}

.keyboard-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.keyboard-arrows {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
  justify-content: flex-end;
}

.arrow-spacer {
  flex: 1;
}

.arrow-row {
  display: flex;
  gap: 4px;
}

.arrow-placeholder {
  width: 48px;
  height: 40px;
}

.keyboard-row {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}

.keyboard-row:last-child {
  margin-bottom: 0;
}

.bottom-row {
  margin-top: 6px;
}

.key-btn {
  flex: 1;
  min-width: 48px;
  height: 40px;
  background: #3C3F41;
  border: 1px solid #4E5052;
  border-radius: 8px;
  color: #A9B7C6;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.key-btn.arrow {
  width: 48px;
  flex: none;
  background: #36393B;
}

.key-btn.arrow:hover {
  background: #4E5052;
}

.key-btn:hover {
  background: #4E5052;
  border-color: #5E6062;
}

.key-btn:active {
  background: #36393B;
  transform: scale(0.95);
}

.key-btn.danger {
  color: #F14C4C;
  border-color: rgba(241, 76, 76, 0.3);
}

.key-btn.danger:hover {
  background: rgba(241, 76, 76, 0.15);
  border-color: rgba(241, 76, 76, 0.5);
}

.key-btn.primary {
  color: #4A9CFF;
  border-color: rgba(74, 156, 255, 0.3);
}

.key-btn.primary:hover {
  background: rgba(74, 156, 255, 0.15);
  border-color: rgba(74, 156, 255, 0.5);
}

.connection-overlay,
.error-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #D4D4D4;
  z-index: 10;
  background: #3C3F41;
  padding: 32px;
  border-radius: 8px;
  border: 1px solid #4E5052;
}

.error-overlay p {
  margin: 0;
  color: #D56C6C;
}

:deep(.n-input) {
  background: #3C3F41;
}

:deep(.n-input .n-input__input-el) {
  color: #A9B7C6;
  caret-color: #A9B7C6;
}

:deep(.n-input .n-input__placeholder) {
  color: #6B7B8C;
}

/* File Popover Overlay */
.file-popover-overlay {
  position: absolute;
  top: 60px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
}

.file-popover {
  background: #3C3F41;
  border: 1px solid #4E5052;
  border-radius: 12px;
  padding: 12px;
  min-width: 350px;
  max-width: 500px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.file-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: #D4D4D4;
  font-size: 14px;
  font-weight: 500;
}

/* File Reference Styles */
.file-reference-list {
  max-height: 250px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.2s;
}

.file-item:hover {
  background: rgba(74, 156, 255, 0.15);
}

.file-name {
  color: #D4D4D4;
  font-size: 13px;
  font-weight: 500;
}

.file-path {
  color: #808080;
  font-size: 11px;
  font-family: 'Courier New', monospace;
  margin-left: auto;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-files {
  padding: 16px;
  text-align: center;
  color: #808080;
  font-size: 13px;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .terminal-header {
    padding: 6px 12px;
    min-height: 36px;
  }

  .header-left {
    flex-wrap: wrap;
    gap: 6px;
  }

  .header-left h3 {
    font-size: 13px;
  }

  .header-right {
    gap: 4px;
  }

  /* Make scroll buttons more prominent on mobile */
  :deep(.n-button-group) {
    .n-button {
      padding: 0 8px;
    }
  }

  .work-dir {
    font-size: 10px;
    max-width: 200px;
  }

  .terminal-container {
    padding: 6px;
  }

  .terminal-input {
    padding: 6px 12px;
    flex-direction: row;
    gap: 8px;
  }

  .custom-keyboard {
    padding: 6px;
  }

  .keyboard-main {
    gap: 8px;
    margin-bottom: 4px;
  }

  .keyboard-left {
    gap: 4px;
  }

  .keyboard-row {
    gap: 4px;
    margin-bottom: 4px;
  }

  .key-btn {
    height: 36px;
    font-size: 13px;
    min-width: 40px;
  }

  .key-btn.arrow {
    width: 40px;
    height: 40px;
  }

  .arrow-placeholder {
    width: 40px;
    height: 36px;
  }

  :deep(.n-input) {
    font-size: 14px;
    flex: 1;
  }

  :deep(.n-button) {
    min-height: 36px;
  }
}

@media (max-width: 480px) {
  .terminal-header {
    padding: 6px 10px;
    min-height: 32px;
  }

  .header-left h3 {
    font-size: 12px;
  }

  .work-dir {
    font-size: 9px;
    max-width: 150px;
  }

  .terminal-container {
    padding: 4px;
  }

  .terminal-input {
    padding: 6px 10px;
  }

  :deep(.n-input) {
    font-size: 13px;
  }
}

/* Voice Input Styles */
.mic-active {
  animation: pulse 1s ease-in-out infinite;
  color: #f44336 !important;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
}

/* Voice input tooltip on mobile */
@media (max-width: 768px) {
  .mic-active {
    animation: pulse-red 0.8s ease-in-out infinite;
  }
}

@keyframes pulse-red {
  0%, 100% {
    color: #f44336;
  }
  50% {
    color: #ff6b6b;
  }
}
</style>
