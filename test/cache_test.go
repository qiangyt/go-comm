package test

import (
	comm "github.com/qiangyt/go-comm/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestFileCache(t *testing.T) {
	// 创建临时缓存目录
	tmpDir := t.TempDir()
	cache := comm.NewFileCache(tmpDir)

	// 测试 GetCacheDir
	if cache.GetCacheDir() != tmpDir {
		t.Errorf("GetCacheDir() = %s, want %s", cache.GetCacheDir(), tmpDir)
	}

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test_source.txt")
	testContent := []byte("test content for cache")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 测试 Has - 缓存不存在
	key := "test_key"
	if cache.Has(key) {
		t.Error("Has() should return false for non-existent cache")
	}

	// 测试 Get - 缓存不存在
	if got := cache.Get(key); got != "" {
		t.Errorf("Get() = %s, want empty string for non-existent cache", got)
	}

	// 测试 Put
	cache.Put(testFile, key)

	// 测试 Has - 缓存存在
	if !cache.Has(key) {
		t.Error("Has() should return true after Put()")
	}

	// 测试 Get - 缓存存在
	cachedPath := cache.Get(key)
	if cachedPath == "" {
		t.Error("Get() should return path after Put()")
	}

	// 验证缓存文件内容
	cachedContent, err := os.ReadFile(cachedPath)
	if err != nil {
		t.Fatalf("Failed to read cached file: %v", err)
	}
	if string(cachedContent) != string(testContent) {
		t.Errorf("Cached content = %s, want %s", cachedContent, testContent)
	}

	// 测试 CopyTo
	destFile := filepath.Join(tmpDir, "dest.txt")
	cache.CopyTo(key, destFile)

	destContent, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(destContent) != string(testContent) {
		t.Errorf("Destination content = %s, want %s", destContent, testContent)
	}

	// 测试 Size
	size := cache.Size()
	if size == 0 {
		t.Error("Size() should return non-zero for non-empty cache")
	}

	// 测试 Delete
	cache.Delete(key)
	if cache.Has(key) {
		t.Error("Has() should return false after Delete()")
	}

	// 重新添加文件用于测试 Clear
	cache.Put(testFile, key)
	cache.Put(testFile, "another_key")

	// 测试 Clear
	cache.Clear()
	if cache.Has(key) {
		t.Error("Has() should return false after Clear()")
	}
	if cache.Has("another_key") {
		t.Error("Has() should return false for all keys after Clear()")
	}
}

func TestFileCacheEmptyKey(t *testing.T) {
	tmpDir := t.TempDir()
	cache := comm.NewFileCache(tmpDir)

	// 测试空key
	if cache.Get("") != "" {
		t.Error("Get() should return empty string for empty key")
	}

	if cache.Has("") {
		t.Error("Has() should return false for empty key")
	}

	// 测试 Delete 空key不应panic
	cache.Delete("")
}

func TestFileCachePanic(t *testing.T) {
	tmpDir := t.TempDir()
	cache := comm.NewFileCache(tmpDir)

	// 测试 Put 空key应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Put() with empty key should panic")
		}
	}()
	cache.Put("nonexistent", "")
}

func TestFileCacheCopyToPanic(t *testing.T) {
	tmpDir := t.TempDir()
	cache := comm.NewFileCache(tmpDir)

	// 测试 CopyTo 不存在的key应该panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("CopyTo() with non-existent key should panic")
		}
	}()
	cache.CopyTo("nonexistent", filepath.Join(tmpDir, "dest.txt"))
}

func TestNewFileCacheDefaultDir(t *testing.T) {
	// 测试使用默认缓存目录
	cache := comm.NewFileCache("")
	cacheDir := cache.GetCacheDir()

	homeDir, _ := os.UserHomeDir()
	expectedDir := filepath.Join(homeDir, ".cache", "amcopy")

	if cacheDir != expectedDir {
		t.Errorf("GetCacheDir() = %s, want %s", cacheDir, expectedDir)
	}

	// 清理
	os.RemoveAll(cacheDir)
}
