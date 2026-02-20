<template>
  <div class="app-container">
    <!-- Setup Screen -->
    <div v-if="!isSetupComplete" class="setup-container">
      <div class="setup-card">
        <!-- Logo -->
        <div class="logo">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>

        <h1>Remote Code</h1>
        <p class="subtitle">Initial Setup</p>

        <!-- Platform Notice -->
        <div v-if="!config.isLocalSupported" class="platform-notice">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <span>Windows detected: Local sessions not supported. You can connect to remote servers.</span>
        </div>

        <!-- Setup Steps -->
        <div class="setup-steps">
          <!-- Step 1: Basic Settings -->
          <div class="step" :class="{ active: currentStep === 1 }">
            <h2>1. Basic Settings</h2>
            <div class="form-group">
              <label>Admin Password</label>
              <div class="input-wrapper">
                <input
                  :type="showPassword ? 'text' : 'password'"
                  v-model="config.adminPassword"
                  placeholder="Enter admin password"
                />
                <button class="toggle-btn" @click="showPassword = !showPassword">
                  {{ showPassword ? 'Hide' : 'Show' }}
                </button>
              </div>
            </div>

            <div v-if="config.isLocalSupported" class="form-group">
              <label>Working Directory</label>
              <div class="input-wrapper">
                <input
                  type="text"
                  v-model="config.allowedDir"
                  placeholder="/Users/username/projects"
                />
                <button class="browse-btn" @click="selectDirectory">Browse</button>
              </div>
            </div>
          </div>

          <!-- Step 2: Network Settings -->
          <div class="step" :class="{ active: currentStep === 2 }">
            <h2>2. Network Settings (Optional)</h2>

            <div class="form-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="config.frpEnabled" />
                <span>Enable FRP Tunnel</span>
              </label>
            </div>

            <div v-if="config.frpEnabled" class="frp-settings">
              <div class="form-group">
                <label>FRP Server Address</label>
                <input
                  type="text"
                  v-model="config.frpServer"
                  placeholder="your-server.com"
                />
              </div>
              <div class="form-group">
                <label>FRP Server Port</label>
                <input
                  type="number"
                  v-model.number="config.frpPort"
                  placeholder="7000"
                />
              </div>
              <div class="form-group">
                <label>FRP Token</label>
                <input
                  type="password"
                  v-model="config.frpToken"
                  placeholder="your-frp-token"
                />
              </div>
            </div>
          </div>

          <!-- Step 3: Ready -->
          <div class="step" :class="{ active: currentStep === 3 }">
            <h2>3. Ready to Start</h2>
            <div class="summary">
              <div class="summary-item">
                <span class="label">Admin Password:</span>
                <span class="value">{{ config.adminPassword ? '••••••••' : 'Not set' }}</span>
              </div>
              <div v-if="config.isLocalSupported" class="summary-item">
                <span class="label">Working Directory:</span>
                <span class="value">{{ config.allowedDir || 'Not set' }}</span>
              </div>
              <div class="summary-item">
                <span class="label">FRP Tunnel:</span>
                <span class="value">{{ config.frpEnabled ? 'Enabled' : 'Disabled' }}</span>
              </div>
              <div class="summary-item">
                <span class="label">Backend Port:</span>
                <span class="value">{{ config.backendPort }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation -->
        <div class="navigation">
          <button
            v-if="currentStep > 1"
            class="btn btn-secondary"
            @click="currentStep--"
          >
            Back
          </button>
          <button
            v-if="currentStep < 3"
            class="btn btn-primary"
            @click="currentStep++"
          >
            Next
          </button>
          <button
            v-else
            class="btn btn-primary"
            @click="saveAndStart"
            :disabled="isStarting"
          >
            {{ isStarting ? 'Starting...' : 'Start Application' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Main App -->
    <div v-else class="main-app">
      <!-- Loading Screen -->
      <div v-if="isLoading" class="loading-screen">
        <div class="spinner"></div>
        <p>Starting services...</p>
      </div>

      <!-- Main Content -->
      <div v-else class="main-content">
        <div class="app-header">
          <h1>Remote Code</h1>
          <div class="header-actions">
            <button class="btn-icon" @click="openInBrowser" title="Open in Browser">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                <polyline points="15 3 21 3 21 9"/>
                <line x1="10" y1="14" x2="21" y2="3"/>
              </svg>
            </button>
            <button class="btn-icon" @click="showSettings = true" title="Settings">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Status Bar -->
        <div class="status-bar">
          <div class="status-item">
            <span class="status-dot" :class="status.backendRunning ? 'online' : 'offline'"></span>
            <span>Backend: {{ status.backendRunning ? 'Running' : 'Stopped' }}</span>
          </div>
          <div v-if="status.frpEnabled" class="status-item">
            <span class="status-dot online"></span>
            <span>FRP: {{ status.frpServer }}</span>
          </div>
        </div>

        <!-- Web Content -->
        <div class="web-content">
          <div class="placeholder">
            <p>Application is running!</p>
            <p>Access at: <a href="#" @click.prevent="openInBrowser">http://localhost:{{ config.backendPort }}</a></p>
            <button class="btn btn-primary" @click="openInBrowser">Open in Browser</button>
          </div>
        </div>
      </div>

      <!-- Settings Modal -->
      <div v-if="showSettings" class="modal-overlay" @click.self="showSettings = false">
        <div class="modal">
          <h2>Settings</h2>
          <div class="form-group">
            <label>Admin Password</label>
            <input type="password" v-model="config.adminPassword" />
          </div>
          <div v-if="config.isLocalSupported" class="form-group">
            <label>Working Directory</label>
            <input type="text" v-model="config.allowedDir" />
          </div>
          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="config.frpEnabled" />
              <span>Enable FRP Tunnel</span>
            </label>
          </div>
          <div class="modal-actions">
            <button class="btn btn-secondary" @click="showSettings = false">Cancel</button>
            <button class="btn btn-primary" @click="saveSettings">Save</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'

export default {
  name: 'App',
  setup() {
    const isSetupComplete = ref(false)
    const currentStep = ref(1)
    const showPassword = ref(false)
    const isStarting = ref(false)
    const isLoading = ref(true)
    const showSettings = ref(false)

    const config = reactive({
      backendPort: 9090,
      frontendPort: 5173,
      jwtSecret: '',
      adminPassword: '',
      allowedDir: '',
      frpEnabled: false,
      frpServer: '',
      frpPort: 7000,
      frpToken: '',
      platform: '',
      isLocalSupported: true
    })

    const status = reactive({
      isRunning: false,
      backendRunning: false,
      backendPort: 9090,
      frpEnabled: false,
      frpServer: ''
    })

    // Load config on mount
    onMounted(async () => {
      try {
        // Check if setup is complete
        const setupComplete = await window.go.main.App.IsSetupComplete()
        isSetupComplete.value = setupComplete

        // Load existing config
        const existingConfig = await window.go.main.App.GetConfig()
        Object.assign(config, existingConfig)

        if (setupComplete) {
          // Start services and load main app
          await startServices()
        }
      } catch (err) {
        console.error('Failed to load config:', err)
      }
    })

    const selectDirectory = async () => {
      try {
        const dir = await window.go.main.App.SelectDirectory()
        if (dir) {
          config.allowedDir = dir
        }
      } catch (err) {
        console.error('Failed to select directory:', err)
      }
    }

    const saveAndStart = async () => {
      isStarting.value = true
      try {
        await window.go.main.App.SaveConfiguration(config)
        isSetupComplete.value = true
        await startServices()
      } catch (err) {
        console.error('Failed to save config:', err)
        alert('Failed to save configuration: ' + err)
      } finally {
        isStarting.value = false
      }
    }

    const startServices = async () => {
      isLoading.value = true
      try {
        await window.go.main.App.StartServices()

        // Poll for status
        const updateStatus = async () => {
          const serviceStatus = await window.go.main.App.GetServiceStatus()
          Object.assign(status, serviceStatus)
        }

        await updateStatus()
        setInterval(updateStatus, 2000)
      } catch (err) {
        console.error('Failed to start services:', err)
      } finally {
        isLoading.value = false
      }
    }

    const openInBrowser = async () => {
      const url = `http://localhost:${config.backendPort}`
      await window.go.main.App.OpenInBrowser(url)
    }

    const saveSettings = async () => {
      try {
        await window.go.main.App.SaveConfiguration(config)
        showSettings.value = false
      } catch (err) {
        console.error('Failed to save settings:', err)
      }
    }

    return {
      isSetupComplete,
      currentStep,
      showPassword,
      isStarting,
      isLoading,
      showSettings,
      config,
      status,
      selectDirectory,
      saveAndStart,
      openInBrowser,
      saveSettings
    }
  }
}
</script>

<style>
/* Global Styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  min-height: 100vh;
  color: #fff;
}

/* App Container */
.app-container {
  width: 100%;
  min-height: 100vh;
}

/* Setup Container */
.setup-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.setup-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 500px;
}

.logo {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #4A9CFF 0%, #667eea 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  box-shadow: 0 4px 20px rgba(74, 156, 255, 0.4);
}

.logo svg {
  width: 32px;
  height: 32px;
  color: #fff;
}

h1 {
  text-align: center;
  font-size: 24px;
  margin-bottom: 8px;
}

.subtitle {
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  margin-bottom: 30px;
}

.platform-notice {
  background: rgba(255, 193, 7, 0.1);
  border: 1px solid rgba(255, 193, 7, 0.3);
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
}

.platform-notice svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  color: #ffc107;
}

/* Steps */
.step {
  display: none;
}

.step.active {
  display: block;
}

.step h2 {
  font-size: 16px;
  margin-bottom: 20px;
  color: #4A9CFF;
}

/* Form */
.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 6px;
}

.input-wrapper {
  display: flex;
  gap: 8px;
}

.input-wrapper input {
  flex: 1;
}

input[type="text"],
input[type="password"],
input[type="number"] {
  width: 100%;
  height: 40px;
  padding: 0 12px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  outline: none;
  transition: all 0.3s;
}

input:focus {
  border-color: #4A9CFF;
  background: rgba(255, 255, 255, 0.1);
}

input::placeholder {
  color: rgba(255, 255, 255, 0.35);
}

.checkbox-label {
  display: flex !important;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #4A9CFF;
}

.frp-settings {
  margin-top: 16px;
  padding-left: 26px;
}

/* Buttons */
.btn {
  height: 40px;
  padding: 0 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: linear-gradient(135deg, #4A9CFF 0%, #667eea 100%);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(74, 156, 255, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.browse-btn,
.toggle-btn {
  padding: 0 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
}

.browse-btn:hover,
.toggle-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

/* Navigation */
.navigation {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 30px;
}

/* Summary */
.summary {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 16px;
}

.summary-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.summary-item:last-child {
  border-bottom: none;
}

.summary-item .label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
}

.summary-item .value {
  color: #fff;
  font-size: 13px;
}

/* Main App */
.main-app {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.loading-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(255, 255, 255, 0.1);
  border-top-color: #4A9CFF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-screen p {
  margin-top: 16px;
  color: rgba(255, 255, 255, 0.7);
}

.main-content {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.app-header h1 {
  font-size: 18px;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background: rgba(255, 255, 255, 0.15);
}

.btn-icon svg {
  width: 18px;
  height: 18px;
}

/* Status Bar */
.status-bar {
  display: flex;
  gap: 20px;
  padding: 8px 24px;
  background: rgba(0, 0, 0, 0.1);
  font-size: 12px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: #4caf50;
}

.status-dot.offline {
  background: #f44336;
}

/* Web Content */
.web-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder {
  text-align: center;
}

.placeholder p {
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 12px;
}

.placeholder a {
  color: #4A9CFF;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #1a1a2e;
  border-radius: 16px;
  padding: 24px;
  width: 100%;
  max-width: 400px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal h2 {
  font-size: 18px;
  margin-bottom: 20px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
</style>
