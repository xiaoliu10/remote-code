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
    language: '语言',
    chinese: '中文',
    english: 'English',
    logout: '退出',
    refresh: '刷新',
    create: '创建',
    delete: '删除',
    open: '打开',
    send: '发送',
    save: '保存',
    cancel: '取消',
    confirm: '确认',
    success: '成功',
    error: '错误',
    loading: '加载中...',
    connecting: '连接中...',
    connected: '已连接',
    disconnected: '已断开'
  },
  dashboard: {
    title: 'Remote Code',
    createSession: '创建新会话',
    sessionName: '会话名称',
    sessionNamePlaceholder: 'my-project',
    workDirectory: '工作目录',
    workDirectoryPlaceholder: '/home/user/projects',
    sessions: '会话列表',
    noSessions: '暂无会话',
    active: '运行中',
    inactive: '未运行',
    default: '默认',
    sessionCreated: '会话创建成功',
    sessionDeleted: '会话删除成功',
    createFailed: '创建会话失败',
    deleteFailed: '删除会话失败',
    deleteConfirm: '确定要删除这个会话吗？',
    refreshSuccess: '会话列表已刷新',
    enterSessionName: '请输入会话名称'
  },
  session: {
    terminal: '远程终端',
    typeCommand: '输入命令...',
    leaveSession: '离开会话',
    leaveConfirm: 'WebSocket 连接将被关闭，确定离开吗？',
    stay: '留下',
    leave: '离开',
    welcomeMessage: 'Remote Code 终端',
    sessionLabel: '会话',
    connecting: '正在连接...'
  },
  login: {
    title: '登录',
    username: '用户名',
    password: '密码',
    loginButton: '登录',
    usernamePlaceholder: '请输入用户名',
    passwordPlaceholder: '请输入密码',
    loginSuccess: '登录成功',
    loginFailed: '登录失败，请检查用户名和密码'
  },
  sidebar: {
    sessions: '会话',
    newSession: '新建会话',
    createSession: '创建会话',
    sessionName: '会话名称',
    sessionNamePlaceholder: 'my-project',
    sessionNameReadOnly: '会话名称创建后不可修改',
    workDirectory: '工作目录',
    workDirectoryPlaceholder: '/home/user/projects',
    workDirectoryNote: '工作目录仅用于显示，修改需重启会话',
    sessionConfig: '会话配置',
    sessionInfoUpdated: '会话信息已更新',
    noSessions: '暂无会话'
  },
  fileExplorer: {
    title: '文件浏览器',
    searchPlaceholder: '搜索文件...',
    emptyFolder: '空文件夹',
    newFile: '新建文件',
    newFolder: '新建文件夹',
    rename: '重命名',
    name: '名称',
    namePlaceholder: '输入名称',
    nameRequired: '请输入名称',
    newName: '新名称',
    newNamePlaceholder: '输入新名称',
    deleteConfirm: '确认删除',
    deleteConfirmMessage: '确定要删除 "{name}" 吗？',
    createSuccess: '创建成功',
    createFailed: '创建失败',
    renameSuccess: '重命名成功',
    renameFailed: '重命名失败',
    deleteSuccess: '删除成功',
    deleteFailed: '删除失败',
    refreshFailed: '刷新失败',
    navigateFailed: '导航失败'
  },
  terminal: {
    title: '远程终端',
    typeCommand: '输入命令... (输入 @ 引用文件)',
    typeCommandRealtime: '实时输入模式... (字符实时发送)',
    welcome: 'Remote Code 终端',
    session: '会话',
    connected: '已连接到会话',
    disconnected: '与服务器断开连接',
    reconnect: '重新连接',
    error: '错误',
    clear: '清屏',
    noFilesFound: '未找到文件',
    notConnected: '未连接到会话',
    selectFile: '选择文件',
    searchFiles: '搜索文件...',
    realtimeMode: '实时输入模式 - 字符实时发送到远端',
    commandMode: '命令模式 - 按 Enter 或点击发送执行命令',
    voiceInput: '语音输入',
    listening: '正在听...',
    voiceNotSupported: '浏览器不支持语音',
    voicePermissionDenied: '请授予麦克风权限',
    noSpeechDetected: '未检测到语音',
    voiceError: '语音识别错误',
    voiceRecognized: '语音识别成功',
    voiceTimeout: '语音输入超时',
    voiceStartFailed: '无法启动语音识别'
  },
  sessions: {
    selectSession: '请从左侧选择一个会话',
    sessionNotFound: '会话未找到'
  }
}
