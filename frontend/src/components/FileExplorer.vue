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
  <div class="file-explorer">
    <!-- Header -->
    <div class="explorer-header">
      <h3>{{ t('fileExplorer.title') }}</h3>
      <n-space>
        <n-button quaternary size="small" @click="handleRefresh">
          <template #icon>
            <n-icon><RefreshIcon /></n-icon>
          </template>
        </n-button>
        <n-button quaternary size="small" @click="showCreateFileMenu">
          <template #icon>
            <n-icon><AddIcon /></n-icon>
          </template>
        </n-button>
      </n-space>
    </div>

    <!-- Search (Optional) -->
    <div class="explorer-search">
      <n-input
        v-model:value="searchText"
        :placeholder="t('fileExplorer.searchPlaceholder')"
        clearable
        size="small"
      >
        <template #prefix>
          <n-icon><SearchIcon /></n-icon>
        </template>
      </n-input>
    </div>

    <!-- Current Path -->
    <div class="explorer-path">
      <n-breadcrumb>
        <n-breadcrumb-item @click="navigateToPath('/')">
          <n-icon><HomeIcon /></n-icon>
        </n-breadcrumb-item>
        <n-breadcrumb-item
          v-for="(segment, index) in pathSegments"
          :key="index"
          @click="navigateToPath(getPathUpTo(index))"
        >
          {{ segment }}
        </n-breadcrumb-item>
      </n-breadcrumb>
    </div>

    <!-- File Tree -->
    <div class="explorer-content">
      <n-spin :show="fileStore.loading">
        <n-empty
          v-if="treeData.length === 0 && !fileStore.loading"
          :description="t('fileExplorer.emptyFolder')"
          size="small"
        />
        <n-tree
          v-else
          :data="treeData"
          :pattern="searchText"
          :selected-keys="selectedKeys"
          :expanded-keys="expandedKeys"
          :render-prefix="renderPrefix"
          :render-label="renderLabel"
          block-line
          virtual-scroll
          style="max-height: 100%"
          @update:selected-keys="handleSelect"
          @update:expanded-keys="handleExpand"
          @load="handleLoad"
          @contextmenu="handleContextMenu"
        />
      </n-spin>
    </div>

    <!-- Context Menu -->
    <n-dropdown
      :show="showContextMenu"
      :x="contextMenuX"
      :y="contextMenuY"
      :options="contextMenuOptions"
      placement="bottom-start"
      @select="handleContextMenuSelect"
      @clickoutside="showContextMenu = false"
    />

    <!-- Create Modal -->
    <n-modal
      v-model:show="showCreateModal"
      preset="dialog"
      :title="createModalTitle"
      :positive-text="t('common.create')"
      :negative-text="t('common.cancel')"
      :loading="creating"
      @positive-click="handleCreateItem"
    >
      <n-form :model="createForm" label-placement="left" label-width="auto">
        <n-form-item :label="t('fileExplorer.name')" required>
          <n-input
            v-model:value="createForm.name"
            :placeholder="t('fileExplorer.namePlaceholder')"
          />
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- Rename Modal -->
    <n-modal
      v-model:show="showRenameModal"
      preset="dialog"
      :title="t('fileExplorer.rename')"
      :positive-text="t('common.confirm')"
      :negative-text="t('common.cancel')"
      :loading="renaming"
      @positive-click="handleRenameItem"
    >
      <n-form :model="renameForm" label-placement="left" label-width="auto">
        <n-form-item :label="t('fileExplorer.newName')" required>
          <n-input
            v-model:value="renameForm.newName"
            :placeholder="t('fileExplorer.newNamePlaceholder')"
          />
        </n-form-item>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  NButton,
  NIcon,
  NSpace,
  NInput,
  NBreadcrumb,
  NBreadcrumbItem,
  NSpin,
  NEmpty,
  NTree,
  NDropdown,
  NModal,
  NForm,
  NFormItem,
  useMessage,
  useDialog,
  type TreeOption,
  type DropdownOption
} from 'naive-ui'
import {
  Refresh as RefreshIcon,
  Add as AddIcon,
  Search as SearchIcon,
  Home as HomeIcon,
  Folder as FolderIcon,
  Document as FileIcon,
  CreateOutline as RenameIcon,
  TrashOutline as DeleteIcon
} from '@vicons/ionicons5'
import { useFileStore } from '@/stores/file'
import { useSessionStore } from '@/stores/session'
import type { FileItem } from '@/api/client'

const { t } = useI18n()
const message = useMessage()
const dialog = useDialog()
const fileStore = useFileStore()
const sessionStore = useSessionStore()

// State
const searchText = ref('')
const selectedKeys = ref<string[]>([])
const expandedKeys = ref<string[]>([])
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const contextMenuTarget = ref<TreeOption | null>(null)

// Create modal
const showCreateModal = ref(false)
const createType = ref<'file' | 'directory'>('file')
const creating = ref(false)
const createForm = ref({ name: '' })

// Rename modal
const showRenameModal = ref(false)
const renaming = ref(false)
const renameForm = ref({ newName: '', oldPath: '' })

// Tree data
interface FileTreeOption extends TreeOption {
  isLeaf: boolean
  data: FileItem
}

const treeData = ref<FileTreeOption[]>([])

// Computed
const createModalTitle = computed(() =>
  createType.value === 'file'
    ? t('fileExplorer.newFile')
    : t('fileExplorer.newFolder')
)

const pathSegments = computed(() => {
  return fileStore.currentPath
    .split('/')
    .filter((s) => s.length > 0)
})

const contextMenuOptions = computed<DropdownOption[]>(() => [
  {
    label: t('fileExplorer.newFile'),
    key: 'newFile',
    icon: () => h(NIcon, null, { default: () => h(FileIcon) })
  },
  {
    label: t('fileExplorer.newFolder'),
    key: 'newFolder',
    icon: () => h(NIcon, null, { default: () => h(FolderIcon) })
  },
  { type: 'divider', key: 'd1' },
  {
    label: t('fileExplorer.rename'),
    key: 'rename',
    icon: () => h(NIcon, null, { default: () => h(RenameIcon) })
  },
  {
    label: t('common.delete'),
    key: 'delete',
    icon: () => h(NIcon, null, { default: () => h(DeleteIcon) })
  }
])

/**
 * Get path up to a certain segment index
 */
function getPathUpTo(index: number) {
  const segments = pathSegments.value.slice(0, index + 1)
  return '/' + segments.join('/')
}

/**
 * Navigate to a specific path
 */
async function navigateToPath(path: string) {
  try {
    await loadDirectory(path)
    fileStore.currentPath = path
  } catch (error) {
    message.error(t('fileExplorer.navigateFailed'))
  }
}

/**
 * Load directory contents
 */
async function loadDirectory(path: string): Promise<FileTreeOption[]> {
  try {
    const files = await fileStore.fetchFiles(path)
    return files.map((file: FileItem) => ({
      key: file.path,
      label: file.name,
      isLeaf: file.type === 'file',
      data: file,
      children: file.type === 'directory' ? [] : undefined
    }))
  } catch (error) {
    console.error('Failed to load directory:', error)
    return []
  }
}

/**
 * Render tree node prefix (icon)
 */
function renderPrefix({ option }: { option: TreeOption }) {
  const fileOption = option as FileTreeOption
  return h(NIcon, null, {
    default: () =>
      h(fileOption.isLeaf ? FileIcon : FolderIcon, {
        style: { color: fileOption.isLeaf ? 'var(--n-text-color)' : '#f0a020' }
      })
  })
}

/**
 * Render tree node label
 */
function renderLabel({ option }: { option: TreeOption }) {
  return h('span', { class: 'file-label' }, option.label)
}

/**
 * Handle tree node selection
 */
async function handleSelect(keys: string[], option: TreeOption[]) {
  selectedKeys.value = keys
  if (keys.length > 0) {
    const selectedOption = option[0] as FileTreeOption
    fileStore.setSelectedFile(selectedOption.key as string)

    // If it's a file, emit event or load content
    if (selectedOption.isLeaf) {
      // Could emit event to parent to show file content in terminal
      console.log('Selected file:', selectedOption.key)
    }
  }
}

/**
 * Handle tree node expansion
 */
function handleExpand(keys: string[]) {
  expandedKeys.value = keys
}

/**
 * Show create file menu
 */
function showCreateFileMenu() {
  createType.value = 'file'
  contextMenuTarget.value = null
  showCreateModal.value = true
}

/**
 * Handle lazy loading of tree nodes
 */
async function handleLoad(option: TreeOption) {
  const fileOption = option as FileTreeOption
  if (fileOption.isLeaf) return

  const children = await loadDirectory(fileOption.key as string)
  fileOption.children = children
}

/**
 * Handle right-click context menu
 */
function handleContextMenu(e: MouseEvent, option: TreeOption[]) {
  e.preventDefault()
  if (option.length > 0) {
    contextMenuTarget.value = option[0]
  } else {
    contextMenuTarget.value = null
  }
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  showContextMenu.value = true
}

/**
 * Handle context menu selection
 */
function handleContextMenuSelect(key: string) {
  showContextMenu.value = false

  switch (key) {
    case 'newFile':
      createType.value = 'file'
      showCreateModal.value = true
      break
    case 'newFolder':
      createType.value = 'directory'
      showCreateModal.value = true
      break
    case 'rename':
      if (contextMenuTarget.value) {
        const fileOption = contextMenuTarget.value as FileTreeOption
        renameForm.value.oldPath = fileOption.key as string
        renameForm.value.newName = fileOption.data.name
        showRenameModal.value = true
      }
      break
    case 'delete':
      if (contextMenuTarget.value) {
        const fileOption = contextMenuTarget.value as FileTreeOption
        confirmDelete(fileOption.key as string, fileOption.data.name)
      }
      break
  }
}

/**
 * Create new file or folder
 */
async function handleCreateItem() {
  if (!createForm.value.name) {
    message.error(t('fileExplorer.nameRequired'))
    return false
  }

  creating.value = true
  try {
    const parentPath = contextMenuTarget.value
      ? (contextMenuTarget.value as FileTreeOption).key
      : fileStore.currentPath
    const newPath = `${parentPath}/${createForm.value.name}`.replace(/\/+/g, '/')

    await fileStore.createFile(newPath, createType.value, '')
    message.success(t('fileExplorer.createSuccess'))
    showCreateModal.value = false
    createForm.value.name = ''

    // Refresh the tree
    await loadDirectory(fileStore.currentPath)
  } catch (error) {
    message.error(t('fileExplorer.createFailed'))
    return false
  } finally {
    creating.value = false
  }
  return true
}

/**
 * Rename file or folder
 */
async function handleRenameItem() {
  if (!renameForm.value.newName) {
    message.error(t('fileExplorer.nameRequired'))
    return false
  }

  renaming.value = true
  try {
    const oldPath = renameForm.value.oldPath
    const parentPath = oldPath.substring(0, oldPath.lastIndexOf('/'))
    const newPath = `${parentPath}/${renameForm.value.newName}`

    await fileStore.renameFile(oldPath, newPath)
    message.success(t('fileExplorer.renameSuccess'))
    showRenameModal.value = false

    // Refresh the tree
    await loadDirectory(fileStore.currentPath)
  } catch (error) {
    message.error(t('fileExplorer.renameFailed'))
    return false
  } finally {
    renaming.value = false
  }
  return true
}

/**
 * Confirm and delete file/folder
 */
function confirmDelete(path: string, name: string) {
  dialog.warning({
    title: t('fileExplorer.deleteConfirm'),
    content: t('fileExplorer.deleteConfirmMessage', { name }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await fileStore.deleteFile(path)
        message.success(t('fileExplorer.deleteSuccess'))
        // Refresh the tree
        await loadDirectory(fileStore.currentPath)
      } catch (error) {
        message.error(t('fileExplorer.deleteFailed'))
      }
    }
  })
}

/**
 * Refresh file list
 */
async function handleRefresh() {
  try {
    await loadDirectory(fileStore.currentPath)
    message.success(t('common.success'))
  } catch (error) {
    message.error(t('fileExplorer.refreshFailed'))
  }
}

// Watch for session changes to update file root
watch(
  () => sessionStore.currentSession,
  async (session) => {
    if (session) {
      const workDir = session.work_dir || '/'
      fileStore.currentPath = workDir
      treeData.value = await loadDirectory(workDir)
    }
  },
  { immediate: true }
)

// Initialize
onMounted(async () => {
  if (sessionStore.currentSession) {
    const workDir = sessionStore.currentSession.work_dir || '/'
    fileStore.currentPath = workDir
    treeData.value = await loadDirectory(workDir)
  } else {
    // Load root as fallback
    treeData.value = await loadDirectory('/')
  }
})
</script>

<style scoped>
.file-explorer {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #2B2B2B;
  border-right: 1px solid #4E5052;
}

.explorer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #4E5052;
  background: #3C3F41;
}

.explorer-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #D4D4D4;
}

.explorer-search {
  padding: 8px 16px;
  border-bottom: 1px solid #4E5052;
  background: #3C3F41;
}

.explorer-path {
  padding: 8px 16px;
  border-bottom: 1px solid #4E5052;
  overflow-x: auto;
  background: #3C3F41;
  color: #808080;
}

.explorer-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  background: #2B2B2B;
}

.file-label {
  font-size: 13px;
  color: #A9B7C6;
}

:deep(.n-tree-node-content) {
  padding: 4px 0;
}

:deep(.n-tree-node-wrapper) {
  padding: 2px 4px;
}

:deep(.n-tree-node-wrapper:hover) {
  background: rgba(255, 255, 255, 0.08);
}

:deep(.n-breadcrumb-item) {
  cursor: pointer;
}

:deep(.n-breadcrumb-item:hover) {
  color: #4A9CFF;
}

:deep(.n-breadcrumb-item:last-child) {
  color: #D4D4D4;
  font-weight: 500;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .file-explorer {
    border-right: none;
    border-bottom: 1px solid #4E5052;
  }

  .explorer-header {
    padding: 10px 12px;
  }

  .explorer-header h3 {
    font-size: 13px;
  }

  .explorer-search {
    padding: 6px 12px;
  }

  .explorer-path {
    padding: 6px 12px;
    font-size: 12px;
  }

  .explorer-content {
    padding: 6px;
    max-height: 150px;
  }

  .file-label {
    font-size: 12px;
  }

  :deep(.n-tree-node-content) {
    padding: 6px 0;
  }

  :deep(.n-button) {
    min-height: 36px;
  }
}

@media (max-width: 480px) {
  .explorer-header {
    padding: 8px 10px;
  }

  .explorer-header h3 {
    font-size: 12px;
  }

  .explorer-search {
    display: none;
  }

  .explorer-content {
    max-height: 120px;
  }
}
</style>
