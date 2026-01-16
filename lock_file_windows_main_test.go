//go:build windows
// +build windows

package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCreateLockFile_Happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	filename := "/test.lock"
	data := "test data"

	f, err := CreateLockFile(fs, filename, data)
	a.NoError(err)
	a.NotNil(f)
	f.Close()
}

func TestCreateLockFile_ReplaceExisting(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	filename := "/test.lock"

	// Create first lock file
	f1, err := CreateLockFile(fs, filename, "data1")
	a.NoError(err)
	a.NotNil(f1)
	f1.Close()

	// Create second lock file - should replace the first
	f2, err := CreateLockFile(fs, filename, "data2")
	a.NoError(err)
	a.NotNil(f2)
	f2.Close()
}
