// Package qtest 提供测试相关的工具函数和类型
package qtest

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
)

// MockFS 用于测试的 mock 文件系统，实现 afero.Fs 接口
type MockFS struct {
	afero.Fs       // 嵌入 MemMapFs 作为基础实现
	MkdirAllErr  error
	OpenFileErr  error
	CallCount    int
}

// NewMockFS 创建 mock 文件系统
func NewMockFS() *MockFS {
	return &MockFS{
		Fs: afero.NewMemMapFs(),
	}
}

// MkdirAll 实现 afero.Fs 接口
func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllErr != nil {
		return m.MkdirAllErr
	}
	return m.Fs.MkdirAll(path, perm)
}

// OpenFile 实现 afero.Fs 接口
func (m *MockFS) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.CallCount++
	if m.OpenFileErr != nil {
		return nil, m.OpenFileErr
	}
	return m.Fs.OpenFile(name, flag, perm)
}

// MockFSWithReadableError 第一次 OpenFile 成功，第二次失败
type MockFSWithReadableError struct {
	*MockFS
}

// OpenFile 实现 afero.Fs 接口
func (m *MockFSWithReadableError) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.MockFS.CallCount++
	if m.MockFS.CallCount == 1 {
		// 第一次调用（protocol.jsonl）成功
		return m.MockFS.Fs.OpenFile(name, flag, perm)
	}
	// 第二次调用（readable.log）失败
	return nil, fmt.Errorf("mock error: readable.log 创建失败")
}

// MockFSWithExistsError 模拟 Exists 返回错误
type MockFSWithExistsError struct {
	*MockFS
}

// Exists 实现 afero.Fs 接口
func (m *MockFSWithExistsError) Exists(path string) (bool, error) {
	return false, fmt.Errorf("mock error: exists check failed")
}

// MockFSWithWriteError 模拟 WriteFile 返回错误
type MockFSWithWriteError struct {
	*MockFS
}

// WriteFile 实现 afero.Fs 接口
func (m *MockFSWithWriteError) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return fmt.Errorf("mock error: write failed")
}

// MockFSWithReadError 模拟 ReadFile 返回错误
type MockFSWithReadError struct {
	*MockFS
}

// ReadFile 实现 afero.Fs 接口
func (m *MockFSWithReadError) ReadFile(filename string) ([]byte, error) {
	return nil, fmt.Errorf("mock error: read failed")
}

// MockFSWithAppendWriteError 模拟追加写入失败
// 当 .gitignore 文件已存在且需要追加时，WriteFile 失败
type MockFSWithAppendWriteError struct {
	*MockFS
}

// WriteFile 实现 afero.Fs 接口
func (m *MockFSWithAppendWriteError) WriteFile(filename string, data []byte, perm os.FileMode) error {
	// 如果是追加 .gitignore（包含 "# crayfish" 且文件已存在）
	if strings.Contains(filename, ".gitignore") && strings.Contains(string(data), "# crayfish") {
		exists, _ := afero.Exists(m.MockFS, filename)
		if exists {
			return fmt.Errorf("mock error: append write failed")
		}
	}
	return m.MockFS.Fs.(interface {
		WriteFile(string, []byte, os.FileMode) error
	}).WriteFile(filename, data, perm)
}
