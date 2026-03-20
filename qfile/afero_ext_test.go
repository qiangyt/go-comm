package qfile

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// EnsureDirWithSubdirs 测试

func TestEnsureDirWithSubdirs_CreatesMainDirAndSubdirs(t *testing.T) {
	fs := afero.NewMemMapFs()
	mainDir := "/test/config"
	subdirs := []string{"logs", "cache"}

	err := EnsureDirWithSubdirs(fs, mainDir, subdirs...)
	require.NoError(t, err)

	// 验证主目录存在
	exists, err := DirExists(fs, mainDir)
	require.NoError(t, err)
	assert.True(t, exists, "Main directory should exist")

	// 验证子目录存在
	for _, subdir := range subdirs {
		fullPath := mainDir + "/" + subdir
		exists, err := DirExists(fs, fullPath)
		require.NoError(t, err)
		assert.True(t, exists, "Subdirectory "+subdir+" should exist")
	}
}

func TestEnsureDirWithSubdirs_MainDirAlreadyExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	mainDir := "/test/config"

	// 先创建主目录
	Mkdir(fs, mainDir)

	// 再次调用应该成功
	err := EnsureDirWithSubdirs(fs, mainDir, "logs")
	require.NoError(t, err)

	// 验证目录存在
	exists, err := DirExists(fs, mainDir+"/logs")
	require.NoError(t, err)
	assert.True(t, exists, "Logs directory should exist")
}

func TestEnsureDirWithSubdirs_EmptySubdirs(t *testing.T) {
	fs := afero.NewMemMapFs()
	mainDir := "/test/config"

	// 空子目录列表
	err := EnsureDirWithSubdirs(fs, mainDir)
	require.NoError(t, err)

	// 验证主目录存在
	exists, err := DirExists(fs, mainDir)
	require.NoError(t, err)
	assert.True(t, exists, "Main directory should exist")
}

func TestEnsureDirWithSubdirsP_Success(t *testing.T) {
	fs := afero.NewMemMapFs()
	mainDir := "/test/config"

	// Panic 版本应该不会 panic
	EnsureDirWithSubdirsP(fs, mainDir, "logs", "cache")

	// 验证目录存在
	exists, err := DirExists(fs, mainDir)
	require.NoError(t, err)
	assert.True(t, exists, "Main directory should exist")

	exists, err = DirExists(fs, mainDir+"/logs")
	require.NoError(t, err)
	assert.True(t, exists, "Logs directory should exist")

	exists, err = DirExists(fs, mainDir+"/cache")
	require.NoError(t, err)
	assert.True(t, exists, "Cache directory should exist")
}

// EnsureFileWithContent 测试

func TestEnsureFileWithContent_FileNotExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/config.yaml"
	content := []byte("test: content")

	created, err := EnsureFileWithContent(fs, filePath, content)
	require.NoError(t, err)
	assert.True(t, created, "File should be created")

	// 验证文件存在
	exists, err := FileExists(fs, filePath)
	require.NoError(t, err)
	assert.True(t, exists, "File should exist")

	// 验证文件内容
	readContent, err := ReadFileBytes(fs, filePath)
	require.NoError(t, err)
	assert.Equal(t, content, readContent, "File content should match")
}

func TestEnsureFileWithContent_FileAlreadyExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/config.yaml"
	originalContent := []byte("original: content")

	// 先创建文件
	Mkdir(fs, "/test/config")
	WriteFile(fs, filePath, originalContent)

	// 尝试再次创建
	newContent := []byte("new: content")
	created, err := EnsureFileWithContent(fs, filePath, newContent)
	require.NoError(t, err)
	assert.False(t, created, "File should not be created")

	// 验证原始内容被保留
	readContent, err := ReadFileBytes(fs, filePath)
	require.NoError(t, err)
	assert.Equal(t, originalContent, readContent, "Original content should be preserved")
}

func TestEnsureFileWithContent_CreatesParentDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/nested/config.yaml"
	content := []byte("test: content")

	created, err := EnsureFileWithContent(fs, filePath, content)
	require.NoError(t, err)
	assert.True(t, created, "File should be created")

	// 验证父目录被创建
	exists, err := DirExists(fs, "/test/config/nested")
	require.NoError(t, err)
	assert.True(t, exists, "Parent directory should be created")
}

func TestEnsureFileWithContent_WithSubdirs(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/config.yaml"
	content := []byte("test: content")

	created, err := EnsureFileWithContent(fs, filePath, content, "logs", "cache")
	require.NoError(t, err)
	assert.True(t, created, "File should be created")

	// 验证子目录被创建
	exists, err := DirExists(fs, "/test/config/logs")
	require.NoError(t, err)
	assert.True(t, exists, "Logs subdirectory should be created")

	exists, err = DirExists(fs, "/test/config/cache")
	require.NoError(t, err)
	assert.True(t, exists, "Cache subdirectory should be created")
}

func TestEnsureFileWithContentP_Success(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/config.yaml"
	content := []byte("test: content")

	// Panic 版本应该不会 panic
	created := EnsureFileWithContentP(fs, filePath, content, "logs")
	assert.True(t, created, "File should be created")

	// 验证文件存在
	exists, err := FileExists(fs, filePath)
	require.NoError(t, err)
	assert.True(t, exists, "File should exist")

	// 验证子目录被创建
	exists, err = DirExists(fs, "/test/config/logs")
	require.NoError(t, err)
	assert.True(t, exists, "Logs subdirectory should be created")
}

func TestEnsureFileWithContentP_FileAlreadyExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := "/test/config/config.yaml"
	originalContent := []byte("original: content")

	// 先创建文件
	Mkdir(fs, "/test/config")
	WriteFile(fs, filePath, originalContent)

	// Panic 版本应该不会 panic
	newContent := []byte("new: content")
	created := EnsureFileWithContentP(fs, filePath, newContent)
	assert.False(t, created, "File should not be created")

	// 验证原始内容被保留
	readContent, err := ReadFileBytes(fs, filePath)
	require.NoError(t, err)
	assert.Equal(t, originalContent, readContent, "Original content should be preserved")
}
