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
 * "AS IS" BASIS, WITHOUT ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fileApi, type FileItem, type FileContent } from '@/api/client'

export const useFileStore = defineStore('file', () => {
  // State
  const currentPath = ref('/')
  const files = ref<FileItem[]>([])
  const selectedFile = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const fileTree = ref<TreeTreeNode[]>([])

  // Tree node interface for NTree component
  interface TreeTreeNode {
    key: string
    label: string
    isLeaf: boolean
    children?: TreeTreeNode[]
    prefix?: () => void
    data: FileItem
  }

  /**
   * Fetch files from a specific path
   */
  async function fetchFiles(path: string) {
    loading.value = true
    error.value = null
    try {
      const response = await fileApi.list(path)
      // 修复: 使用 items 而不是 files，匹配后端响应格式
      files.value = response.data.items || []
      currentPath.value = response.data.path || path
      return files.value
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch files'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Get file content
   */
  async function getFileContent(path: string): Promise<FileContent> {
    loading.value = true
    error.value = null
    try {
      const response = await fileApi.getContent(path)
      return response.data
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to get file content'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new file or directory
   */
  async function createFile(path: string, type: 'file' | 'directory', content?: string) {
    loading.value = true
    error.value = null
    try {
      const response = await fileApi.create(path, type, content)
      // Refresh the current directory
      await fetchFiles(currentPath.value)
      return response.data
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to create file'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Rename a file or directory
   */
  async function renameFile(oldPath: string, newPath: string) {
    loading.value = true
    error.value = null
    try {
      await fileApi.rename(oldPath, newPath)
      // Refresh the current directory
      await fetchFiles(currentPath.value)
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to rename file'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Delete a file or directory
   */
  async function deleteFile(path: string) {
    loading.value = true
    error.value = null
    try {
      await fileApi.delete(path)
      // Refresh the current directory
      await fetchFiles(currentPath.value)
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to delete file'
      throw e
    } finally {
      loading.value = false
    }
  }

  /**
   * Set selected file
   */
  function setSelectedFile(path: string | null) {
    selectedFile.value = path
  }

  /**
   * Convert files to tree structure for NTree component
   */
  function buildTreeNodes(items: FileItem[], parentPath: string = ''): TreeTreeNode[] {
    return items.map((item) => {
      const fullPath = parentPath ? `${parentPath}/${item.name}` : item.name
      return {
        key: fullPath,
        label: item.name,
        isLeaf: item.type === 'file',
        children: item.type === 'directory' ? [] : undefined,
        data: item
      }
    })
  }

  /**
   * Update file tree
   */
  function updateFileTree(nodes: TreeTreeNode[]) {
    fileTree.value = nodes
  }

  return {
    // State
    currentPath,
    files,
    selectedFile,
    loading,
    error,
    fileTree,
    // Actions
    fetchFiles,
    getFileContent,
    createFile,
    renameFile,
    deleteFile,
    setSelectedFile,
    buildTreeNodes,
    updateFileTree
  }
})
