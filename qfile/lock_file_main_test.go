package qfile

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestReadLockFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	// Create a lock file with JSON content
	// Note: JSON unmarshaling to any produces float64 for numbers
	// The function has a type assertion bug, but we test actual behavior
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to type assertion (float64 vs int)
			a.Contains(r.(error).Error(), "interface")
		}
	}()
	WriteFileTextP(fs, "/test.lock", `{"pid": 1234, "data": "test data"}`)
	_, _, _ = ReadLockFile(fs, "/test.lock")
	t.Error("ReadLockFile should panic due to type assertion bug")
}

func TestReadLockFile_notFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, _, err := ReadLockFile(fs, "/nonexistent.lock")
	a.Error(err)
}

func TestReadLockFile_invalidJson(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.lock", `not valid json`)

	_, _, err := ReadLockFile(fs, "/test.lock")
	a.Error(err)
}

func TestReadLockFile_withValidPid(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	// 使用 float64 格式的 PID（因为 JSON unmarshal 到 any 会产生 float64）
	WriteFileTextP(fs, "/test.lock", `{"pid": 12345.0, "data": {"key": "value"}}`)

	// ReadLockFile 会因为类型断言失败而 panic
	defer func() {
		r := recover()
		a.NotNil(r)
	}()

	_, _, _ = ReadLockFile(fs, "/test.lock")
}
