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
  <div class="layout-container">
    <!-- Top Header -->
    <header class="layout-header">
      <div class="header-left">
        <n-button
          v-if="sidebarCollapsed"
          quaternary
          size="small"
          @click="toggleSidebar"
        >
          <template #icon>
            <n-icon><MenuIcon /></n-icon>
          </template>
        </n-button>
        <h1>{{ t('dashboard.title') }}</h1>
      </div>
      <div class="header-right">
        <LanguageSwitcher />
        <n-button quaternary @click="handleLogout">
          <template #icon>
            <n-icon><LogoutIcon /></n-icon>
          </template>
          {{ t('common.logout') }}
        </n-button>
      </div>
    </header>

    <!-- Main Two-Column Layout -->
    <div class="main-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
      <!-- Session Sidebar -->
      <SessionSidebar
        class="sidebar-panel"
        :style="{ width: sidebarWidth + 'px' }"
        @toggle="toggleSidebar"
      />
      <div
        v-if="!sidebarCollapsed"
        class="resize-handle"
        @mousedown="startResizeSidebar"
      />

      <!-- Main Content / Terminal -->
      <div class="main-content-wrapper">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NButton, NIcon } from 'naive-ui'
import { LogOut as LogoutIcon, Menu as MenuIcon } from '@vicons/ionicons5'
import SessionSidebar from '@/components/SessionSidebar.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const { t } = useI18n()
const router = useRouter()

// Panel widths
const sidebarWidth = ref(250)
const sidebarCollapsed = ref(false)

// Mobile detection
const isMobile = ref(false)
const MOBILE_BREAKPOINT = 768

// Check if device is mobile
function checkMobile() {
  isMobile.value = window.innerWidth < MOBILE_BREAKPOINT
  if (isMobile.value) {
    // Auto-collapse sidebar on mobile
    sidebarCollapsed.value = true
    // Set smaller panel width
    sidebarWidth.value = 200
  }
}

// Constraints (responsive)
const SIDEBAR_MIN = computed(() => isMobile.value ? 180 : 200)
const SIDEBAR_MAX = computed(() => isMobile.value ? 280 : 400)

// Resize state
let isResizing = false

/**
 * Start resizing sidebar
 */
function startResizeSidebar(e: MouseEvent) {
  e.preventDefault()
  isResizing = true
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

/**
 * Handle resize movement
 */
function handleResize(e: MouseEvent) {
  if (!isResizing) return

  const newWidth = Math.max(SIDEBAR_MIN.value, Math.min(SIDEBAR_MAX.value, e.clientX))
  sidebarWidth.value = newWidth
}

/**
 * Stop resizing
 */
function stopResize() {
  isResizing = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''

  // Save to localStorage
  localStorage.setItem('sidebarWidth', String(sidebarWidth.value))
}

/**
 * Toggle sidebar collapsed state
 */
function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('sidebarCollapsed', String(sidebarCollapsed.value))
}

/**
 * Load saved panel widths from localStorage
 */
function loadSavedWidths() {
  // Only load saved widths on desktop
  if (isMobile.value) return

  const savedSidebarWidth = localStorage.getItem('sidebarWidth')
  const savedSidebarCollapsed = localStorage.getItem('sidebarCollapsed')

  if (savedSidebarWidth) {
    sidebarWidth.value = Math.max(SIDEBAR_MIN.value, Math.min(SIDEBAR_MAX.value, parseInt(savedSidebarWidth)))
  }
  if (savedSidebarCollapsed) {
    sidebarCollapsed.value = savedSidebarCollapsed === 'true'
  }
}

/**
 * Handle logout
 */
function handleLogout() {
  localStorage.removeItem('token')
  router.push('/login')
}

/**
 * Handle window resize (with debounce)
 */
let resizeTimeout: number | null = null
function handleWindowResize() {
  if (resizeTimeout) {
    clearTimeout(resizeTimeout)
  }
  resizeTimeout = window.setTimeout(() => {
    checkMobile()
    if (!isMobile.value) {
      loadSavedWidths()
    }
    resizeTimeout = null
  }, 100)
}

// Lifecycle
onMounted(() => {
  checkMobile()
  loadSavedWidths()
  window.addEventListener('resize', handleWindowResize)
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  window.removeEventListener('resize', handleWindowResize)
})
</script>

<style scoped>
.layout-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
  background: #2B2B2B;
}

.layout-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  height: 48px;
  background: #3C3F41;
  border-bottom: 1px solid #4E5052;
  flex-shrink: 0;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h1 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #D4D4D4;
}

.main-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
  background: #2B2B2B;
}

.sidebar-panel {
  flex-shrink: 0;
  overflow: hidden;
  transition: width 0.3s ease;
  background: #2B2B2B;
}

.resize-handle {
  width: 4px;
  flex-shrink: 0;
  background: #4E5052;
  cursor: col-resize;
  transition: background 0.2s;
  position: relative;
  z-index: 10;
}

.resize-handle:hover {
  background: #4A9CFF;
}

.resize-handle::after {
  content: '';
  position: absolute;
  top: 0;
  left: -4px;
  right: -4px;
  bottom: 0;
}

.main-content-wrapper {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: #2B2B2B;
}

/* Collapsed state */
.main-layout.sidebar-collapsed .sidebar-panel {
  width: 0 !important;
  overflow: hidden;
}

/* Mobile Responsive - Vertical Layout */
@media (max-width: 768px) {
  .layout-header {
    height: 52px;
    padding: 0 12px;
  }

  .header-left h1 {
    font-size: 14px;
  }

  .header-right {
    gap: 8px;
  }

  /* Change to vertical layout */
  .main-layout {
    flex-direction: column;
  }

  /* Hide resize handles */
  .resize-handle {
    display: none;
  }

  /* Sessions panel - full width, collapsible */
  .sidebar-panel {
    width: 100% !important;
    height: auto;
    max-height: 200px;
    border-right: none;
    border-bottom: 1px solid #4E5052;
    transition: max-height 0.3s ease;
  }

  .main-layout.sidebar-collapsed .sidebar-panel {
    max-height: 0 !important;
    border-bottom: none;
  }

  /* Terminal - takes remaining space */
  .main-content-wrapper {
    flex: 1;
    min-height: 300px;
  }
}

@media (max-width: 480px) {
  .layout-header {
    height: 48px;
    padding: 0 8px;
  }

  .header-left h1 {
    font-size: 13px;
  }

  :deep(.n-button) {
    padding: 0 8px;
  }

  /* Smaller heights on very small screens */
  .sidebar-panel {
    max-height: 180px;
  }

  .main-content-wrapper {
    min-height: 250px;
  }
}
</style>
