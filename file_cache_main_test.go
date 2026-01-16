package comm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFileCache_withDir(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache")
	cache := NewFileCache(cacheDir)

	a.NotNil(cache)
	a.Equal(cacheDir, cache.GetCacheDir())

	// Clean up
	os.RemoveAll(cacheDir)
}

func TestNewFileCache_defaultDir(t *testing.T) {
	a := require.New(t)

	cache := NewFileCache("")
	a.NotNil(cache)
	a.NotEmpty(cache.GetCacheDir())
}

func TestFileCache_Has(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-has")
	cache := NewFileCache(cacheDir)
	defer os.RemoveAll(cacheDir)

	a.False(cache.Has("nonexistent"))
}

func TestFileCache_Get_emptyKey(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-get")
	cache := NewFileCache(cacheDir)
	defer os.RemoveAll(cacheDir)

	result := cache.Get("")
	a.Equal("", result)
}

func TestFileCache_Get_notFound(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-get2")
	cache := NewFileCache(cacheDir)
	defer os.RemoveAll(cacheDir)

	result := cache.Get("nonexistent-key")
	a.Equal("", result)
}

func TestFileCache_Put(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-put")
	defer os.RemoveAll(cacheDir)

	// Create a source file
	srcFile := filepath.Join(os.TempDir(), "src-file.txt")
	os.WriteFile(srcFile, []byte("test content"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile, "test-key")

	result := cache.Get("test-key")
	a.NotEmpty(result)
	a.Equal(true, cache.Has("test-key"))
}

func TestFileCache_Put_emptyKey(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-put2")
	defer os.RemoveAll(cacheDir)

	srcFile := filepath.Join(os.TempDir(), "src-file2.txt")
	os.WriteFile(srcFile, []byte("test"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Put should panic on empty key")
		}
	}()

	cache.Put(srcFile, "")
}

func TestFileCache_Put_duplicate(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-put3")
	defer os.RemoveAll(cacheDir)

	srcFile := filepath.Join(os.TempDir(), "src-file3.txt")
	os.WriteFile(srcFile, []byte("test content"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile, "test-key")
	cache.Put(srcFile, "test-key") // Should not panic, just skip

	result := cache.Get("test-key")
	a.NotEmpty(result)
}

func TestFileCache_CopyTo(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-copy")
	defer os.RemoveAll(cacheDir)

	destFile := filepath.Join(os.TempDir(), "dest-file.txt")
	defer os.Remove(destFile)

	srcFile := filepath.Join(os.TempDir(), "src-file-copy.txt")
	os.WriteFile(srcFile, []byte("test content"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile, "test-key")

	cache.CopyTo("test-key", destFile)

	content, err := os.ReadFile(destFile)
	a.NoError(err)
	a.Equal("test content", string(content))
}

func TestFileCache_CopyTo_notFound(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-copy2")
	defer os.RemoveAll(cacheDir)

	destFile := filepath.Join(os.TempDir(), "dest-file2.txt")
	defer os.Remove(destFile)

	cache := NewFileCache(cacheDir)

	defer func() {
		if r := recover(); r == nil {
			t.Error("CopyTo should panic when key not found")
		}
	}()

	cache.CopyTo("nonexistent", destFile)
}

func TestFileCache_Delete(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-delete")
	defer os.RemoveAll(cacheDir)

	srcFile := filepath.Join(os.TempDir(), "src-file-del.txt")
	os.WriteFile(srcFile, []byte("test"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile, "test-key")

	a.True(cache.Has("test-key"))
	cache.Delete("test-key")
	a.False(cache.Has("test-key"))
}

func TestFileCache_Delete_emptyKey(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "test-cache-delete2")
	defer os.RemoveAll(cacheDir)

	cache := NewFileCache(cacheDir)

	// Should not panic
	cache.Delete("")
}

func TestFileCache_Clear(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-clear")
	defer os.RemoveAll(cacheDir)

	srcFile1 := filepath.Join(os.TempDir(), "src-file1.txt")
	srcFile2 := filepath.Join(os.TempDir(), "src-file2.txt")
	os.WriteFile(srcFile1, []byte("test1"), 0644)
	os.WriteFile(srcFile2, []byte("test2"), 0644)
	defer os.Remove(srcFile1)
	defer os.Remove(srcFile2)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile1, "key1")
	cache.Put(srcFile2, "key2")

	a.True(cache.Has("key1"))
	a.True(cache.Has("key2"))

	cache.Clear()

	a.False(cache.Has("key1"))
	a.False(cache.Has("key2"))
}

func TestFileCache_Size(t *testing.T) {
	a := require.New(t)

	cacheDir := filepath.Join(os.TempDir(), "test-cache-size")
	defer os.RemoveAll(cacheDir)

	srcFile := filepath.Join(os.TempDir(), "src-file-size.txt")
	os.WriteFile(srcFile, []byte("test content"), 0644)
	defer os.Remove(srcFile)

	cache := NewFileCache(cacheDir)
	cache.Put(srcFile, "test-key")

	size := cache.Size()
	a.Greater(size, int64(0))
}
