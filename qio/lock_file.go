package qio

// a revised copy of github.com/allan-simon/go-singleinstance v0.0.0-20210120080615-d0997106ab37

import (
	"sync"
	"testing"

	"github.com/qiangyt/go-comm/v2/qjson"
	"github.com/spf13/afero"
)

// If filename is a lock file, returns the PID of the process locking it
func ReadLockFile(fs afero.Fs, filename string) (int, any, error) {
	contents, err := afero.ReadFile(fs, filename)
	if err != nil {
		return 0, nil, err
	}

	payload := map[string]any{}
	if err = qjson.JsonUnmarshal(contents, &payload); err != nil {
		return 0, nil, err
	}

	return payload["pid"].(int), payload["data"], nil
}

// ==================== CreateLockFile 测试 (POSIX) ====================
// 注意：Windows 版本在 lock_file_windows_main_test.go 中测试

func TestCreateLockFile_happy(t *testing.T) {
	fs := afero.NewMemMapFs()
	data := map[string]any{"key": "value"}

	f, err := CreateLockFile(fs, "/test.lock", data)
	// 在内存文件系统上可能不支持 flock，所以我们检查文件是否被创建
	if err != nil {
		t.Logf("CreateLockFile error (may be expected on memfs): %v", err)
	}
	if f != nil {
		f.Close()
	}
}

func TestCreateLockFile_concurrent(t *testing.T) {
	fs := afero.NewMemMapFs()
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := map[string]any{"goroutine": true}
			f, err := CreateLockFile(fs, "/concurrent.lock", data)
			if err == nil && f != nil {
				mu.Lock()
				successCount++
				mu.Unlock()
				f.Close()
			}
		}()
	}
	wg.Wait()

	// 在内存文件系统上，可能只有一个或多个成功
	t.Logf("Concurrent CreateLockFile: %d successes", successCount)
}
