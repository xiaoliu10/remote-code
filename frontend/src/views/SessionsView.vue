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
  <div class="sessions-view">
    <!-- No Session Selected -->
    <div v-if="!currentSession" class="no-session">
      <n-empty :description="t('sessions.selectSession')" />
    </div>

    <!-- Terminal View -->
    <EnhancedTerminal
      v-else
      :session-name="currentSession.name"
      @connected="handleConnected"
      @disconnected="handleDisconnected"
      @error="handleError"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { NEmpty, useMessage } from 'naive-ui'
import EnhancedTerminal from '@/components/EnhancedTerminal.vue'
import { useSessionStore } from '@/stores/session'

const { t } = useI18n()
const route = useRoute()
const message = useMessage()
const sessionStore = useSessionStore()

// Current session
const currentSession = computed(() => sessionStore.currentSession)

/**
 * Handle connection established
 */
function handleConnected() {
  console.log('Terminal connected')
}

/**
 * Handle disconnection
 */
function handleDisconnected() {
  console.log('Terminal disconnected')
}

/**
 * Handle terminal error
 */
function handleError(errorMsg: string) {
  message.error(errorMsg)
}

/**
 * Load session from route param
 */
async function loadSessionFromRoute() {
  const sessionName = route.params.name as string
  if (sessionName) {
    // Try to find the session
    let session = sessionStore.getSessionByName(sessionName)

    if (!session) {
      // Fetch sessions if not loaded
      await sessionStore.fetchSessions()
      session = sessionStore.getSessionByName(sessionName)
    }

    if (session) {
      sessionStore.setCurrentSession(session)
    } else {
      message.warning(t('sessions.sessionNotFound'))
    }
  }
}

// Watch route changes
watch(
  () => route.params.name,
  () => {
    loadSessionFromRoute()
  }
)

// Initialize
onMounted(async () => {
  await sessionStore.fetchSessions()
  loadSessionFromRoute()
})
</script>

<style scoped>
.sessions-view {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.no-session {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--n-color);
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .no-session {
    padding: 16px;
  }

  :deep(.n-empty__description) {
    font-size: 14px;
  }

  :deep(.n-button) {
    min-height: 36px;
  }
}

@media (max-width: 480px) {
  .no-session {
    padding: 12px;
  }

  :deep(.n-empty__description) {
    font-size: 13px;
  }

  :deep(.n-button) {
    min-height: 32px;
    font-size: 13px;
  }
}
</style>
