package qio

import (
	"path/filepath"
	"strings"
)

// IsPathAllowed 检查路径是否在允许的目录列表内
// 用于安全检查，防止路径遍历攻击
//
// 参数:
//   - path: 要检查的路径
//   - allowedDirs: 允许的目录列表
//
// 返回:
//   - true: 路径在允许的目录内
//   - false: 路径不在允许的目录内或路径无效
//
// 示例:
//
//	IsPathAllowed("/home/user/data/file.txt", []string{"/home/user/data"}) // true
//	IsPathAllowed("/home/user/data/../etc/passwd", []string{"/home/user/data"}) // false
//	IsPathAllowed("/home/user/data/subdir/file.txt", []string{"/home/user/data", "/tmp"}) // true
func IsPathAllowed(path string, allowedDirs []string) bool {
	if len(allowedDirs) == 0 {
		return false
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// 解析符号链接和相对路径
	absPath = filepath.Clean(absPath)

	for _, dir := range allowedDirs {
		if isInsideDir(absPath, dir) {
			return true
		}
	}
	return false
}

// isInsideDir 检查 absPath 是否在 dir 目录内
func isInsideDir(absPath string, dir string) bool {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	absDir = filepath.Clean(absDir)

	// 检查路径是否在允许的目录内
	rel, err := filepath.Rel(absDir, absPath)
	if err != nil {
		return false
	}

	// 如果相对路径不以 .. 开头，说明在允许的目录内
	// rel == "." 表示路径就是允许的目录本身
	// rel[0] != '.' 表示相对路径不以 .. 开头（即不是上级目录）
	if rel == "." {
		return true
	}
	if rel != "" && !strings.HasPrefix(rel, "..") {
		return true
	}
	return false
}

// IsPathAllowedP 是 IsPathAllowed 的 panic 版本
// 当路径无效时 panic 而不是返回 false
func IsPathAllowedP(path string, allowedDirs []string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	absPath = filepath.Clean(absPath)

	for _, dir := range allowedDirs {
		if isInsideDir(absPath, dir) {
			return true
		}
	}
	return false
}
