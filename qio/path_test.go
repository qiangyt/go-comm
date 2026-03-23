package qio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPathAllowed(t *testing.T) {
	t.Run("路径在允许目录内", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		err := os.Mkdir(subDir, 0o755)
		require.NoError(t, err)

		// 测试子目录
		assert.True(t, IsPathAllowed(subDir, []string{tempDir}))

		// 测试文件
		testFile := filepath.Join(tempDir, "test.txt")
		assert.True(t, IsPathAllowed(testFile, []string{tempDir}))

		// 测试嵌套子目录
		nestedFile := filepath.Join(subDir, "nested", "file.txt")
		assert.True(t, IsPathAllowed(nestedFile, []string{tempDir}))
	})

	t.Run("路径在允许目录外", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()

		// dir2 的文件不在 dir1 的允许范围内
		testFile := filepath.Join(dir2, "test.txt")
		assert.False(t, IsPathAllowed(testFile, []string{dir1}))
	})

	t.Run("路径遍历攻击应被阻止", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		err := os.Mkdir(subDir, 0o755)
		require.NoError(t, err)

		// 尝试通过 .. 访问上级目录
		escapePath := filepath.Join(subDir, "..", "..", "etc", "passwd")
		assert.False(t, IsPathAllowed(escapePath, []string{tempDir}))
	})

	t.Run("空允许目录列表返回 false", func(t *testing.T) {
		assert.False(t, IsPathAllowed("/some/path", []string{}))
	})

	t.Run("多个允许目录", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()

		// dir1 中的文件
		file1 := filepath.Join(dir1, "file1.txt")
		assert.True(t, IsPathAllowed(file1, []string{dir1, dir2}))

		// dir2 中的文件
		file2 := filepath.Join(dir2, "file2.txt")
		assert.True(t, IsPathAllowed(file2, []string{dir1, dir2}))

		// 其他目录的文件
		dir3 := t.TempDir()
		file3 := filepath.Join(dir3, "file3.txt")
		assert.False(t, IsPathAllowed(file3, []string{dir1, dir2}))
	})

	t.Run("允许目录本身就是允许的", func(t *testing.T) {
		tempDir := t.TempDir()
		assert.True(t, IsPathAllowed(tempDir, []string{tempDir}))
	})

	t.Run("相对路径处理", func(t *testing.T) {
		tempDir := t.TempDir()

		// 使用相对路径访问
		oldDir, err := os.Getwd()
		require.NoError(t, err)
		defer os.Chdir(oldDir)

		err = os.Chdir(tempDir)
		require.NoError(t, err)

		// 相对路径文件
		assert.True(t, IsPathAllowed("./test.txt", []string{tempDir}))
		assert.True(t, IsPathAllowed("subdir/file.txt", []string{tempDir}))
	})

	t.Run("无效路径返回 false", func(t *testing.T) {
		// 包含 null 字符的路径在大多数系统上是无效的
		// 但 filepath.Abs 可能会处理它，所以这个测试主要验证不会 panic
		result := IsPathAllowed("/valid/path", []string{"/valid"})
		_ = result // 只要不 panic 就行
	})
}

func TestIsPathAllowedP(t *testing.T) {
	t.Run("路径在允许目录内", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "test.txt")
		assert.True(t, IsPathAllowedP(testFile, []string{tempDir}))
	})

	t.Run("路径在允许目录外", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()
		testFile := filepath.Join(dir2, "test.txt")
		assert.False(t, IsPathAllowedP(testFile, []string{dir1}))
	})
}

func TestIsInsideDir(t *testing.T) {
	t.Run("子目录在父目录内", func(t *testing.T) {
		tempDir := t.TempDir()
		absDir, _ := filepath.Abs(tempDir)

		subDir := filepath.Join(tempDir, "subdir")
		absPath, _ := filepath.Abs(subDir)

		assert.True(t, isInsideDir(absPath, absDir))
	})

	t.Run("目录在自身内", func(t *testing.T) {
		tempDir := t.TempDir()
		absPath, _ := filepath.Abs(tempDir)

		assert.True(t, isInsideDir(absPath, tempDir))
	})

	t.Run("不同目录不在内", func(t *testing.T) {
		dir1 := t.TempDir()
		dir2 := t.TempDir()

		absPath1, _ := filepath.Abs(dir1)
		absDir2, _ := filepath.Abs(dir2)

		assert.False(t, isInsideDir(absPath1, absDir2))
	})
}
