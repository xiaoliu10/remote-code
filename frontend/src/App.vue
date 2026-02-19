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
  <n-config-provider :theme="darkTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <router-view />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { darkTheme } from 'naive-ui'
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui'
import type { GlobalThemeOverrides } from 'naive-ui'

// IntelliJ IDEA Darcula theme colors
const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#4A9CFF',
    primaryColorHover: '#66B1FF',
    primaryColorPressed: '#3385E6',
    primaryColorSuppl: '#4A9CFF',
    bodyColor: '#2B2B2B',
    cardColor: '#3C3F41',
    modalColor: '#3C3F41',
    popoverColor: '#3C3F41',
    tableColor: '#2B2B2B',
    inputColor: '#3C3F41',
    actionColor: '#3C3F41',
    hoverColor: 'rgba(255, 255, 255, 0.09)',
    borderColor: '#4E5052',
    dividerColor: '#4E5052',
    scrollbarColor: 'rgba(255, 255, 255, 0.2)',
    scrollbarColorHover: 'rgba(255, 255, 255, 0.3)',
  },
  Menu: {
    color: '#2B2B2B',
    itemColorHover: 'rgba(255, 255, 255, 0.08)',
    itemColorActive: 'rgba(74, 156, 255, 0.15)',
    itemColorActiveHover: 'rgba(74, 156, 255, 0.2)',
    itemTextColor: '#A9B7C6',
    itemTextColorHover: '#D4D4D4',
    itemTextColorActive: '#4A9CFF',
    itemTextColorActiveHover: '#66B1FF',
    itemTextColorChildActive: '#4A9CFF',
    arrowColor: '#808080',
  },
  Button: {
    textColor: '#A9B7C6',
    textColorHover: '#D4D4D4',
    textColorPressed: '#D4D4D4',
    color: '#3C3F41',
    colorHover: '#4E5052',
    colorPressed: '#36393B',
    borderColor: '#4E5052',
    borderColorHover: '#5E6062',
    borderColorPressed: '#4E5052',
  },
  Input: {
    textColor: '#A9B7C6',
    textColorHover: '#D4D4D4',
    placeholderColor: '#6B7B8C',
    placeholderColorHover: '#808994',
    color: '#3C3F41',
    colorHover: '#45484A',
    colorFocus: '#3C3F41',
    borderColor: '#4E5052',
    borderColorHover: '#5E6062',
    borderColorFocus: '#4A9CFF',
  },
  Tree: {
    textColor: '#A9B7C6',
    textColorHover: '#D4D4D4',
    nodeColorHover: 'rgba(255, 255, 255, 0.08)',
    nodeColorActive: 'rgba(74, 156, 255, 0.15)',
    nodeTextColorActive: '#4A9CFF',
    arrowColor: '#808080',
  },
  Tag: {
    textColor: '#A9B7C6',
  },
  Card: {
    color: '#3C3F41',
    textColor: '#A9B7C6',
    titleTextColor: '#D4D4D4',
    borderColor: '#4E5052',
  },
  Breadcrumb: {
    textColor: '#808080',
    textColorHover: '#4A9CFF',
  },
  Empty: {
    textColor: '#808080',
  },
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  overflow: hidden;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell,
    'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #2B2B2B;
  color: #A9B7C6;
  /* Prevent iOS bounce scrolling */
  position: fixed;
  width: 100%;
  overscroll-behavior: none;
}

#app {
  width: 100%;
  height: 100%;
  background: #2B2B2B;
  overflow: auto;
  /* Use dvh for better mobile keyboard handling */
  height: 100dvh;
}

/* Mobile keyboard handling */
@media (max-width: 768px) {
  body {
    /* Prevent zoom on input focus on iOS */
    -webkit-text-size-adjust: 100%;
  }
}
</style>
