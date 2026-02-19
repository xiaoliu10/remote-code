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

package handlers

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// MaxFileSize 最大文件大小限制 (10MB)
	MaxFileSize = 10 * 1024 * 1024
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 100
	// MaxPageSize 最大分页大小
	MaxPageSize = 1000
)

// PathValidator 路径验证器，确保所有文件操作都在允许的目录内
type PathValidator struct {
	allowedDir string
}

// NewPathValidator 创建路径验证器
func NewPathValidator(allowedDir string) (*PathValidator, error) {
	// 确保允许的目录是绝对路径
	absDir, err := filepath.Abs(allowedDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve allowed directory: %w", err)
	}

	// 确保允许的目录存在
	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("allowed directory does not exist: %w", err)
	}
	if !info.IsDir() {
		return nil, errors.New("allowed directory path is not a directory")
	}

	return &PathValidator{allowedDir: absDir}, nil
}

// Validate 验证路径是否在允许的目录内
// 返回清理后的绝对路径，如果验证失败则返回错误
func (v *PathValidator) Validate(path string) (string, error) {
	// 1. 清理路径（去除 ..、. 等）
	cleanPath := filepath.Clean(path)

	// 2. 拼接完整路径
	fullPath := filepath.Join(v.allowedDir, cleanPath)

	// 3. 解析为绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	// 4. 验证路径在允许的目录内
	if !strings.HasPrefix(absPath, v.allowedDir) {
		return "", errors.New("path is outside allowed directory")
	}

	// 5. 检查符号链接，防止逃逸
	if err := v.checkSymlinkSafety(absPath); err != nil {
		return "", err
	}

	return absPath, nil
}

// checkSymlinkSafety 检查路径中的符号链接是否安全
func (v *PathValidator) checkSymlinkSafety(path string) error {
	// 逐级检查路径中的每个组件
	current := v.allowedDir
	components := strings.Split(strings.TrimPrefix(path, v.allowedDir), string(filepath.Separator))

	for _, comp := range components {
		if comp == "" {
			continue
		}

		current = filepath.Join(current, comp)

		// 检查当前路径组件是否为符号链接
		info, err := os.Lstat(current)
		if err != nil {
			// 文件不存在是允许的（创建文件时）
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to stat path: %w", err)
		}

		// 如果是符号链接，解析并验证目标
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(current)
			if err != nil {
				return fmt.Errorf("failed to read symlink: %w", err)
			}

			// 如果是相对路径，转换为绝对路径
			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Dir(current), target)
			}

			// 清理目标路径
			target = filepath.Clean(target)

			// 验证目标路径是否在允许的目录内
			if !strings.HasPrefix(target, v.allowedDir) {
				return errors.New("symlink target is outside allowed directory")
			}
		}
	}

	return nil
}

// GetAllowedDir 返回允许的根目录
func (v *PathValidator) GetAllowedDir() string {
	return v.allowedDir
}

// FileHandler 文件操作处理器
type FileHandler struct {
	validator *PathValidator
}

// NewFileHandler 创建文件处理器
func NewFileHandler(validator *PathValidator) *FileHandler {
	return &FileHandler{
		validator: validator,
	}
}

// FileInfo 文件信息结构
type FileInfo struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Type       string    `json:"type"` // "file" or "directory"
	Size       int64     `json:"size"`
	ModTime    time.Time `json:"modTime"`
	Permission string    `json:"permission"`
}

// ListDirectoryResponse 列出目录响应
type ListDirectoryResponse struct {
	Path     string     `json:"path"`
	Items    []FileInfo `json:"items"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
}

// ListDirectoryRequest 列出目录请求参数
type ListDirectoryRequest struct {
	Path     string `form:"path" binding:"-"`
	Page     int    `form:"page" binding:"-"`
	PageSize int    `form:"pageSize" binding:"-"`
}

// ListDirectory 列出目录内容
// GET /api/files?path=/project&page=1&pageSize=100
func (h *FileHandler) ListDirectory(c *gin.Context) {
	var req ListDirectoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request parameters",
			"code":  "INVALID_PARAMS",
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}
	if req.PageSize > MaxPageSize {
		req.PageSize = MaxPageSize
	}

	// 验证路径
	validPath, err := h.validator.Validate(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 读取目录
	entries, err := os.ReadDir(validPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "directory not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		if os.IsPermission(err) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
				"code":  "PERMISSION_DENIED",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read directory",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 转换为 FileInfo 列表
	items := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue // 跳过无法读取的文件
		}

		// 计算相对路径（相对于 allowedDir）
		relPath, _ := filepath.Rel(h.validator.GetAllowedDir(), filepath.Join(validPath, entry.Name()))

		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}

		items = append(items, FileInfo{
			Name:       entry.Name(),
			Path:       normalizePath(relPath),
			Type:       fileType,
			Size:       info.Size(),
			ModTime:    info.ModTime(),
			Permission: info.Mode().String(),
		})
	}

	// 排序：目录在前，然后按名称排序
	sort.Slice(items, func(i, j int) bool {
		if items[i].Type != items[j].Type {
			return items[i].Type == "directory"
		}
		return items[i].Name < items[j].Name
	})

	// 分页
	total := len(items)
	start := (req.Page - 1) * req.PageSize
	if start > total {
		start = total
	}
	end := start + req.PageSize
	if end > total {
		end = total
	}

	// 计算相对路径用于响应
	relPath, _ := filepath.Rel(h.validator.GetAllowedDir(), validPath)

	c.JSON(http.StatusOK, ListDirectoryResponse{
		Path:     normalizePath(relPath),
		Items:    items[start:end],
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// GetFileContentRequest 获取文件内容请求参数
type GetFileContentRequest struct {
	Path string `form:"path" binding:"required"`
}

// GetFileContent 获取文件内容
// GET /api/files/content?path=/project/file.txt
func (h *FileHandler) GetFileContent(c *gin.Context) {
	var req GetFileContentRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is required",
			"code":  "INVALID_PARAMS",
		})
		return
	}

	// 验证路径
	validPath, err := h.validator.Validate(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 检查文件是否存在且是文件（不是目录）
	info, err := os.Stat(validPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		if os.IsPermission(err) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
				"code":  "PERMISSION_DENIED",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to stat file",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 检查是否为目录
	if info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "path is a directory, not a file",
			"code":  "IS_DIRECTORY",
		})
		return
	}

	// 检查文件大小
	if info.Size() > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("file too large (max %d bytes)", MaxFileSize),
			"code":  "FILE_TOO_LARGE",
		})
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(validPath)
	if err != nil {
		if os.IsPermission(err) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "permission denied",
				"code":  "PERMISSION_DENIED",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to read file",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 返回文本内容
	c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
}

// CreateFileFolderRequest 创建文件/文件夹请求
type CreateFileFolderRequest struct {
	Path    string `json:"path" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=file directory"`
	Content string `json:"content"`
}

// CreateFileFolder 创建文件或文件夹
// POST /api/files
func (h *FileHandler) CreateFileFolder(c *gin.Context) {
	var req CreateFileFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"code":    "INVALID_PARAMS",
			"details": err.Error(),
		})
		return
	}

	// 验证路径
	validPath, err := h.validator.Validate(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 检查是否已存在
	if _, err := os.Stat(validPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "file or directory already exists",
			"code":  "ALREADY_EXISTS",
		})
		return
	}

	switch req.Type {
	case "directory":
		// 创建目录（包括所有父目录）
		if err := os.MkdirAll(validPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create directory",
				"code":  "INTERNAL_ERROR",
			})
			return
		}

	case "file":
		// 确保父目录存在
		parentDir := filepath.Dir(validPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create parent directory",
				"code":  "INTERNAL_ERROR",
			})
			return
		}

		// 写入文件内容
		content := []byte(req.Content)
		if len(content) > MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("content too large (max %d bytes)", MaxFileSize),
				"code":  "CONTENT_TOO_LARGE",
			})
			return
		}

		if err := os.WriteFile(validPath, content, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create file",
				"code":  "INTERNAL_ERROR",
			})
			return
		}
	}

	// 返回创建的资源信息
	info, _ := os.Stat(validPath)
	relPath, _ := filepath.Rel(h.validator.GetAllowedDir(), validPath)

	c.JSON(http.StatusCreated, FileInfo{
		Name:       filepath.Base(validPath),
		Path:       normalizePath(relPath),
		Type:       req.Type,
		Size:       info.Size(),
		ModTime:    info.ModTime(),
		Permission: info.Mode().String(),
	})
}

// RenameFileFolderRequest 重命名请求
type RenameFileFolderRequest struct {
	OldPath string `json:"oldPath" binding:"required"`
	NewPath string `json:"newPath" binding:"required"`
}

// RenameFileFolder 重命名文件或文件夹
// PUT /api/files/rename
func (h *FileHandler) RenameFileFolder(c *gin.Context) {
	var req RenameFileFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"code":    "INVALID_PARAMS",
			"details": err.Error(),
		})
		return
	}

	// 验证旧路径
	oldPath, err := h.validator.Validate(req.OldPath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "old path: " + err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 验证新路径
	newPath, err := h.validator.Validate(req.NewPath)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "new path: " + err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 检查源是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "source file or directory not found",
			"code":  "NOT_FOUND",
		})
		return
	}

	// 检查目标是否已存在
	if _, err := os.Stat(newPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "destination already exists",
			"code":  "ALREADY_EXISTS",
		})
		return
	}

	// 确保目标父目录存在
	parentDir := filepath.Dir(newPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create destination parent directory",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 执行重命名
	if err := os.Rename(oldPath, newPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to rename",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 返回重命名后的资源信息
	info, _ := os.Stat(newPath)
	relPath, _ := filepath.Rel(h.validator.GetAllowedDir(), newPath)

	c.JSON(http.StatusOK, FileInfo{
		Name:       filepath.Base(newPath),
		Path:       normalizePath(relPath),
		Type:       getFileType(info),
		Size:       info.Size(),
		ModTime:    info.ModTime(),
		Permission: info.Mode().String(),
	})
}

// DeleteFileFolderRequest 删除请求参数
type DeleteFileFolderRequest struct {
	Path string `form:"path" binding:"required"`
}

// DeleteFileFolder 删除文件或文件夹
// DELETE /api/files?path=/project/file.txt
func (h *FileHandler) DeleteFileFolder(c *gin.Context) {
	var req DeleteFileFolderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "path parameter is required",
			"code":  "INVALID_PARAMS",
		})
		return
	}

	// 验证路径
	validPath, err := h.validator.Validate(req.Path)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  "INVALID_PATH",
		})
		return
	}

	// 检查是否存在
	info, err := os.Stat(validPath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "file or directory not found",
				"code":  "NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to stat path",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// 执行删除
	if info.IsDir() {
		// 删除目录及其内容
		if err := os.RemoveAll(validPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete directory",
				"code":  "INTERNAL_ERROR",
			})
			return
		}
	} else {
		// 删除文件
		if err := os.Remove(validPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete file",
				"code":  "INTERNAL_ERROR",
			})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

// getFileType 获取文件类型
func getFileType(info fs.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

// CopyFile 复制文件的辅助方法（可选功能）
func (h *FileHandler) CopyFile(src, dst string) error {
	// 验证源路径
	srcPath, err := h.validator.Validate(src)
	if err != nil {
		return err
	}

	// 验证目标路径
	dstPath, err := h.validator.Validate(dst)
	if err != nil {
		return err
	}

	// 打开源文件
	sourceFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 获取源文件信息
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// 检查文件大小
	if sourceInfo.Size() > MaxFileSize {
		return fmt.Errorf("file too large (max %d bytes)", MaxFileSize)
	}

	// 创建目标文件
	destFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 复制内容
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 保持权限
	return os.Chmod(dstPath, sourceInfo.Mode())
}

// normalizePath 规范化路径显示
// 将 "." 转换为 "/"，确保路径格式正确
func normalizePath(relPath string) string {
    if relPath == "." || relPath == "" {
        return "/"
    }
    return "/" + relPath
}
