package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestFallbackFilePath_happy(t *testing.T) {
	a := require.New(t)

	result := FallbackFilePath("/fallback/dir", "http://example.com/file.txt")
	a.NotEmpty(result)
	// Just check that it's not empty and contains a hash-like string
	a.True(len(result) > 20)
}

func TestHasFallbackFile_notExists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	exists, err := HasFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.NoError(err)
	a.False(exists)
}

func TestHasFallbackFile_exists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/fallback")
	WriteFileTextP(fs, FallbackFilePath("/fallback", "http://example.com/file.txt"), "content")

	exists, err := HasFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.NoError(err)
	a.True(exists)
}

func TestReadFallbackFile_notExists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	path, bytes, err := ReadFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.NoError(err)
	a.Empty(path)
	a.Nil(bytes)
}

func TestReadFallbackFile_exists(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/fallback")
	fallbackPath := FallbackFilePath("/fallback", "http://example.com/file.txt")
	WriteFileTextP(fs, fallbackPath, "test content")

	path, bytes, err := ReadFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.NoError(err)
	a.Equal(fallbackPath, path)
	a.Equal([]byte("test content"), bytes)
}

func TestWriteFallbackFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/fallback")

	path, err := WriteFallbackFile("/fallback", fs, "http://example.com/file.txt", []byte("test content"))
	a.NoError(err)
	a.NotEmpty(path)

	// Verify file exists
	exists, _ := HasFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.True(exists)
}

func TestWriteFallbackFile_overwrite(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/fallback")

	// Write first time
	path1, err := WriteFallbackFile("/fallback", fs, "http://example.com/file.txt", []byte("content1"))
	a.NoError(err)

	// Overwrite
	path2, err := WriteFallbackFile("/fallback", fs, "http://example.com/file.txt", []byte("content2"))
	a.NoError(err)
	a.Equal(path1, path2)

	// Verify content was updated
	_, content, _ := ReadFallbackFile("/fallback", fs, "http://example.com/file.txt")
	a.Equal([]byte("content2"), content)
}
