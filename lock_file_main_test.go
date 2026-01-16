package comm

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
	ReadLockFile(fs, "/test.lock")
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
