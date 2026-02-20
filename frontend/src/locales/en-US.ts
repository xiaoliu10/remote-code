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

export default {
  common: {
    language: 'Language',
    chinese: '中文',
    english: 'English',
    logout: 'Logout',
    refresh: 'Refresh',
    create: 'Create',
    delete: 'Delete',
    open: 'Open',
    send: 'Send',
    save: 'Save',
    cancel: 'Cancel',
    confirm: 'Confirm',
    success: 'Success',
    error: 'Error',
    loading: 'Loading...',
    connecting: 'Connecting...',
    connected: 'Connected',
    disconnected: 'Disconnected'
  },
  dashboard: {
    title: 'Remote Code',
    createSession: 'Create New Session',
    sessionName: 'Session Name',
    sessionNamePlaceholder: 'my-project',
    workDirectory: 'Work Directory',
    workDirectoryPlaceholder: '/home/user/projects',
    sessions: 'Sessions',
    noSessions: 'No sessions yet',
    active: 'Active',
    inactive: 'Inactive',
    default: 'Default',
    sessionCreated: 'Session created successfully',
    sessionDeleted: 'Session deleted successfully',
    createFailed: 'Failed to create session',
    deleteFailed: 'Failed to delete session',
    deleteConfirm: 'Are you sure you want to delete this session?',
    refreshSuccess: 'Sessions refreshed',
    enterSessionName: 'Please enter a session name'
  },
  session: {
    terminal: 'Remote Terminal',
    typeCommand: 'Type command...',
    leaveSession: 'Leave Session?',
    leaveConfirm: 'Your WebSocket connection will be closed.',
    stay: 'Stay',
    leave: 'Leave',
    welcomeMessage: 'Remote Code Terminal',
    sessionLabel: 'Session',
    connecting: 'Connecting...'
  },
  login: {
    title: 'Login',
    username: 'Username',
    password: 'Password',
    loginButton: 'Login',
    usernamePlaceholder: 'Enter username',
    passwordPlaceholder: 'Enter password',
    loginSuccess: 'Login successful',
    loginFailed: 'Login failed. Please check your credentials'
  },
  sidebar: {
    sessions: 'Sessions',
    newSession: 'New Session',
    createSession: 'Create Session',
    sessionName: 'Session Name',
    sessionNamePlaceholder: 'my-project',
    sessionNameReadOnly: 'Session name cannot be changed after creation',
    workDirectory: 'Work Directory',
    workDirectoryPlaceholder: '/home/user/projects',
    workDirectoryReadOnly: 'Work directory cannot be changed after session creation. Please delete and recreate the session to change it',
    sessionConfig: 'Session Config',
    noSessions: 'No sessions'
  },
  fileExplorer: {
    title: 'File Explorer',
    searchPlaceholder: 'Search files...',
    emptyFolder: 'Empty folder',
    newFile: 'New File',
    newFolder: 'New Folder',
    rename: 'Rename',
    name: 'Name',
    namePlaceholder: 'Enter name',
    nameRequired: 'Name is required',
    newName: 'New Name',
    newNamePlaceholder: 'Enter new name',
    deleteConfirm: 'Confirm Delete',
    deleteConfirmMessage: 'Are you sure you want to delete "{name}"?',
    createSuccess: 'Created successfully',
    createFailed: 'Failed to create',
    renameSuccess: 'Renamed successfully',
    renameFailed: 'Failed to rename',
    deleteSuccess: 'Deleted successfully',
    deleteFailed: 'Failed to delete',
    refreshFailed: 'Failed to refresh',
    navigateFailed: 'Failed to navigate'
  },
  terminal: {
    title: 'Remote Terminal',
    typeCommand: 'Type command... (use @ to reference files)',
    typeCommandRealtime: 'Realtime input mode... (chars sent live)',
    welcome: 'Remote Code Terminal',
    session: 'Session',
    connected: 'Connected to session',
    disconnected: 'Disconnected from server',
    reconnect: 'Reconnect',
    error: 'Error',
    clear: 'Clear',
    noFilesFound: 'No files found',
    notConnected: 'Not connected to session',
    selectFile: 'Select File',
    searchFiles: 'Search files...',
    realtimeMode: 'Realtime Mode - chars sent to remote immediately',
    commandMode: 'Command Mode - press Enter or Send to execute',
    voiceInput: 'Voice Input',
    listening: 'Listening...',
    voiceNotSupported: 'Voice not supported',
    voicePermissionDenied: 'Please grant microphone permission',
    noSpeechDetected: 'No speech detected',
    voiceError: 'Voice recognition error',
    voiceRecognized: 'Voice recognized',
    voiceTimeout: 'Voice input timeout',
    voiceStartFailed: 'Failed to start voice recognition',
    scrollToTop: 'Scroll to top',
    scrollToBottom: 'Scroll to bottom',
    pageUp: 'Page up',
    pageDown: 'Page down'
  },
  sessions: {
    selectSession: 'Please select a session from the sidebar',
    sessionNotFound: 'Session not found'
  }
}
