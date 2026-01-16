package comm

import (
	"runtime"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFileOps_GetFs(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)
	result := ops.GetFs()
	a.Equal(fs, result)
}

func TestFileOps_FileExists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	// File doesn't exist
	a.False(ops.FileExists("/test.txt"))

	// Create file
	WriteFileTextP(fs, "/test.txt", "content")
	a.True(ops.FileExists("/test.txt"))
}

func TestFileOps_DirExists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	// Dir doesn't exist
	a.False(ops.DirExists("/testdir"))

	// Create dir
	MkdirP(fs, "/testdir")
	a.True(ops.DirExists("/testdir"))

	// File exists but not dir
	WriteFileTextP(fs, "/testfile.txt", "content")
	a.False(ops.DirExists("/testfile.txt"))
}

func TestFileOps_GetFileSize(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/test.txt", "hello world")
	size := ops.GetFileSize("/test.txt")
	a.Equal(int64(11), size)
}

func TestFileOps_GetFileSize_NotFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	defer func() {
		if r := recover(); r != nil {
			_ = a
			a.Contains(r.(string), "FailedToGetFileSize")
		}
	}()
	ops.GetFileSize("/nonexistent.txt")
	t.Error("GetFileSize should panic on non-existent file")
}

func TestFileOps_EnsureDir(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	ops.EnsureDir("/test/dir")
	a.True(ops.DirExists("/test/dir"))
}

func TestFileOps_CopyFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/source.txt", "test content")
	ops.CopyFile("/source.txt", "/dest.txt")

	content := ReadFileTextP(fs, "/dest.txt")
	a.Equal("test content", content)
}

func TestFileOps_MoveFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/source.txt", "test content")
	ops.MoveFile("/source.txt", "/dest.txt")

	// Source should be gone
	a.False(ops.FileExists("/source.txt"))
	// Dest should exist
	a.True(ops.FileExists("/dest.txt"))
	content := ReadFileTextP(fs, "/dest.txt")
	a.Equal("test content", content)
}

func TestFileOps_RenameFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/oldname.txt", "content")
	newPath := ops.RenameFile("/oldname.txt", "newname.txt")

	// Should return new path
	a.Contains(newPath, "newname.txt")
	a.True(ops.FileExists(newPath))
	a.False(ops.FileExists("/oldname.txt"))
}

func TestFileOps_RenameFile_EmptyName(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	result := ops.RenameFile("/test.txt", "")
	a.Equal("/test.txt", result)
}

func TestFileOps_CreateDirectory(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	ops.CreateDirectory("/test/dir", "", "", false, "")
	a.True(ops.DirExists("/test/dir"))
}

func TestFileOps_RemoveFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/test.txt", "content")
	a.True(ops.FileExists("/test.txt"))

	ops.RemoveFile("/test.txt", false, "")
	a.False(ops.FileExists("/test.txt"))
}

func TestFileOps_RemoveFile_EmptyPath(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	// Should not panic with empty path
	ops.RemoveFile("", false, "")
	_ = a
}

func TestFileOps_SetPermissions_Empty(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/test.txt", "content")
	// Should not panic with empty chmod
	ops.SetPermissions("/test.txt", "")
	_ = a
}

func TestFileOps_SetPermissions_Invalid(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/test.txt", "content")

	if runtime.GOOS != "windows" {
		defer func() {
			if r := recover(); r != nil {
				_ = a
				a.Contains(r.(string), "InvalidChmodValue")
			}
		}()
		ops.SetPermissions("/test.txt", "invalid")
		t.Error("SetPermissions should panic on invalid chmod")
	}
}

func TestFileOps_SetOwner_Empty(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	WriteFileTextP(fs, "/test.txt", "content")
	// Should not panic with empty chown
	ops.SetOwner("/test.txt", "", false, "")
	_ = a
}

func TestFileOps_CreateSymlink_Empty(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	// Should not panic with empty linkPath
	ops.CreateSymlink("/target", "", false, "")
	_ = a
}

func TestFileOps_SetLocalizeFunc(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewFileOps(fs)

	customLocalize := func(key string, args ...map[string]any) string {
		return "custom: " + key
	}
	ops.SetLocalizeFunc(customLocalize)
	// Just verify it doesn't panic - the function will be used in other operations
	_ = a
}
