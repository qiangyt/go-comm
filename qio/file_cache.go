package qio

import (
	"io"
	"os"
	"path/filepath"

	"github.com/qiangyt/go-comm/v3/qerr"
)

// FileCache 文件缓存管理器
type FileCacheT struct {
	cacheDir string
}

type FileCache = *FileCacheT

// NewFileCache 创建文件缓存管理器
// cacheDir: 缓存目录路径，如果为空则使用 ~/.cache/<appName>
func NewFileCache(cacheDir string) FileCache {
	if cacheDir == "" {
		// 获取用户缓存目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(qerr.NewSystemError("get user home directory", err))
		}
		cacheDir = filepath.Join(homeDir, ".cache", "amcopy")
	}

	// 创建缓存目录
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		panic(qerr.NewSystemError("create cache directory", err))
	}

	return &FileCacheT{
		cacheDir: cacheDir,
	}
}

// Get 获取缓存文件路径
// 如果文件存在返回路径，否则返回空字符串
func (me FileCache) Get(key string) string {
	if key == "" {
		return ""
	}

	cachedPath := me.getCachedPath(key)
	if _, err := os.Stat(cachedPath); err == nil {
		return cachedPath
	}

	return ""
}

func (me FileCache) getCachedPath(key string) string {
	if key == "" {
		return ""
	}

	return filepath.Join(me.cacheDir, key)
}

// Has 检查缓存是否存在
func (me FileCache) Has(key string) bool {
	return me.Get(key) != ""
}

// Put 将文件添加到缓存
func (me FileCache) Put(srcPath string, key string) {
	if key == "" {
		panic(qerr.NewBusinessError("cache key is required", nil))
	}

	cachedPath := me.getCachedPath(key)

	// 如果缓存已存在，不需要再复制
	if _, err := os.Stat(cachedPath); err == nil {
		return
	}

	// 复制文件到缓存目录
	src, err := os.Open(srcPath)
	if err != nil {
		panic(qerr.NewSystemError("open source file", err))
	}
	defer src.Close()

	dst, err := os.Create(cachedPath)
	if err != nil {
		panic(qerr.NewSystemError("create cache file", err))
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		// 删除不完整的缓存文件
		os.Remove(cachedPath)
		panic(qerr.NewSystemError("copy file to cache", err))
	}
}

// CopyTo 从缓存复制文件到目标位置
func (me FileCache) CopyTo(key string, destPath string) {
	cachedPath := me.Get(key)
	if cachedPath == "" {
		panic(qerr.NewBusinessError("file not found in cache", nil))
	}

	// 确保目标目录存在
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		panic(qerr.NewSystemError("create destination directory", err))
	}

	// 复制文件
	src, err := os.Open(cachedPath)
	if err != nil {
		panic(qerr.NewSystemError("open cached file", err))
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		panic(qerr.NewSystemError("create destination file", err))
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(destPath)
		panic(qerr.NewSystemError("copy from cache", err))
	}
}

// Delete 删除缓存文件
func (me FileCache) Delete(key string) {
	if key == "" {
		return
	}

	cachedPath := me.getCachedPath(key)
	if err := os.Remove(cachedPath); err != nil && !os.IsNotExist(err) {
		panic(qerr.NewSystemError("delete cache file", err))
	}
}

// Clear 清空所有缓存
func (me FileCache) Clear() {
	entries, err := os.ReadDir(me.cacheDir)
	if err != nil {
		panic(qerr.NewSystemError("read cache directory", err))
	}

	for _, entry := range entries {
		path := me.getCachedPath(entry.Name())
		if err := os.RemoveAll(path); err != nil {
			panic(qerr.NewSystemError("remove cache entry "+entry.Name(), err))
		}
	}
}

// GetCacheDir 获取缓存目录路径
func (me FileCache) GetCacheDir() string {
	return me.cacheDir
}

// Size 获取缓存总大小（字节）
func (me FileCache) Size() int64 {
	var size int64

	err := filepath.Walk(me.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		panic(qerr.NewSystemError("calculate cache size", err))
	}

	return size
}
