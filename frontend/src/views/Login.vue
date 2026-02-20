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
  <div class="login-container">
    <!-- Background decoration -->
    <div class="bg-circle bg-circle-1"></div>
    <div class="bg-circle bg-circle-2"></div>

    <div class="login-card">
      <!-- Logo -->
      <div class="logo-wrapper">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
      </div>

      <h1 class="login-title">Remote Code</h1>
      <p class="login-subtitle">Sign in to continue</p>

      <div class="form-wrapper">
        <div class="input-group">
          <span class="input-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </span>
          <input
            v-model="form.username"
            type="text"
            class="custom-input"
            placeholder="Username"
            @keydown.enter="handleLogin"
          />
        </div>

        <div class="input-group">
          <span class="input-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </span>
          <input
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            class="custom-input"
            placeholder="Password"
            @keydown.enter="handleLogin"
          />
          <button type="button" class="toggle-password" @click="showPassword = !showPassword">
            <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
              <line x1="1" y1="1" x2="23" y2="23"/>
            </svg>
          </button>
        </div>

        <button
          type="button"
          class="login-button"
          :disabled="loading"
          @click="handleLogin"
        >
          <span v-if="loading" class="spinner"></span>
          <span v-else>Sign In</span>
        </button>

        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import * as api from '@/api/client'

const router = useRouter()

const loading = ref(false)
const error = ref('')
const showPassword = ref(false)

const form = reactive({
  username: '',
  password: ''
})

async function handleLogin() {
  if (!form.username.trim()) {
    error.value = 'Please enter your username'
    return
  }
  if (!form.password.trim()) {
    error.value = 'Please enter your password'
    return
  }

  try {
    loading.value = true
    error.value = ''

    const response = await api.authApi.login(form)
    localStorage.setItem('token', response.data.token)

    router.push('/dashboard')
  } catch (e: unknown) {
    error.value = 'Invalid username or password'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

/* Background decoration circles */
.bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
  pointer-events: none;
}

.bg-circle-1 {
  width: 500px;
  height: 500px;
  background: #4A9CFF;
  top: -150px;
  right: -150px;
}

.bg-circle-2 {
  width: 400px;
  height: 400px;
  background: #667eea;
  bottom: -100px;
  left: -100px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 1;
}

.logo-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.logo-icon {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #4A9CFF 0%, #667eea 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(74, 156, 255, 0.4);
}

.logo-icon svg {
  width: 32px;
  height: 32px;
  color: #fff;
}

.login-title {
  margin: 0 0 8px 0;
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  text-align: center;
}

.login-subtitle {
  margin: 0 0 32px 0;
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  text-align: center;
}

.form-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-group {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 14px;
  width: 20px;
  height: 20px;
  color: #4A9CFF;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  z-index: 1;
}

.input-icon svg {
  width: 20px;
  height: 20px;
}

.custom-input {
  width: 100%;
  height: 48px;
  padding: 0 44px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: #fff;
  font-size: 15px;
  outline: none;
  transition: all 0.3s ease;
}

.custom-input::placeholder {
  color: rgba(255, 255, 255, 0.35);
}

.custom-input:hover {
  border-color: rgba(255, 255, 255, 0.2);
}

.custom-input:focus {
  border-color: #4A9CFF;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: 0 0 0 3px rgba(74, 156, 255, 0.15);
}

.toggle-password {
  position: absolute;
  right: 14px;
  width: 20px;
  height: 20px;
  background: none;
  border: none;
  color: #4A9CFF;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: all 0.2s;
}

.toggle-password:hover {
  color: #667eea;
  transform: scale(1.1);
}

.toggle-password svg {
  width: 20px;
  height: 20px;
}

.login-button {
  height: 48px;
  margin-top: 8px;
  background: linear-gradient(135deg, #4A9CFF 0%, #667eea 100%);
  border: none;
  border-radius: 12px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(74, 156, 255, 0.4);
}

.login-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-message {
  padding: 12px 16px;
  background: rgba(244, 67, 54, 0.15);
  border: 1px solid rgba(244, 67, 54, 0.3);
  border-radius: 10px;
  color: #f44336;
  font-size: 14px;
  text-align: center;
}

/* Mobile Responsive */
@media (max-width: 480px) {
  .login-card {
    padding: 32px 24px;
  }

  .logo-icon {
    width: 50px;
    height: 50px;
  }

  .logo-icon svg {
    width: 28px;
    height: 28px;
  }

  .login-title {
    font-size: 22px;
  }

  .custom-input {
    height: 44px;
    font-size: 14px;
  }

  .login-button {
    height: 44px;
    font-size: 15px;
  }
}
</style>
