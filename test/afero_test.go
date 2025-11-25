package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_CopyFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.MkdirP(fs, "/Test_CopyFile_happy/c1")
	comm.WriteFileTextP(fs, "/Test_CopyFile_happy/c1/src.txt", "hello")

	comm.MkdirP(fs, "/Test_CopyFile_happy/c2")
	_, err := comm.CopyFile(fs, "/Test_CopyFile_happy/c1/src.txt", "/Test_CopyFile_happy/c2/dest.txt")
	a.NoError(err)

	actual := comm.ReadFileTextP(fs, "/Test_CopyFile_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_CopyFile_SourceFileNotFound(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.MkdirP(fs, "/Test_CopyFile_SourceFileNotFound/c1")
	comm.MkdirP(fs, "/Test_CopyFile_SourceFileNotFound/c2")

	a.Panicsf(func() {
		comm.CopyFileP(fs, "/Test_CopyFile_SourceFileNotFound/c1/src.txt", "/Test_CopyFile_SourceFileNotFound/c2/dest.txt")
	}, "file not exists: /Test_CopyFile_SourceFileNotFound/c1/src.txt")
}

func Test_Rename_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.MkdirP(fs, "/Test_Rename_happy/c1")
	comm.WriteFileTextP(fs, "/Test_Rename_happy/c1/src.txt", "hello")

	comm.MkdirP(fs, "/Test_Rename_happy/c2")
	comm.RenameP(fs, "/Test_Rename_happy/c1/src.txt", "/Test_Rename_happy/c2/dest.txt")

	actual := comm.ReadFileTextP(fs, "/Test_Rename_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func Test_ReadFileLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.MkdirP(fs, "/Test_ReadAsLines_happy")
	comm.WriteFileLinesP(fs, "/Test_ReadAsLines_happy/f.txt",
		"line 1",
		"line 2")

	actual := comm.ReadFileLinesP(fs, "/Test_ReadAsLines_happy/f.txt")
	a.Equal([]string{
		"line 1",
		"line 2",
	}, actual)
}

func Test_ListSuffixedFiles_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.MkdirP(fs, "Test_ListFilesWithExt_happy")
	comm.MkdirP(fs, "Test_ListFilesWithExt_happy/.")
	comm.MkdirP(fs, "Test_ListFilesWithExt_happy/..")
	comm.MkdirP(fs, "Test_ListFilesWithExt_happy/d.hosts.txt")

	comm.WriteFileTextP(fs, "Test_ListFilesWithExt_happy/1.hosts.txt", "1")
	comm.WriteFileTextP(fs, "Test_ListFilesWithExt_happy/2.hosts.txt.not", "2")
	comm.WriteFileTextP(fs, "Test_ListFilesWithExt_happy/3.hosts.not.text", "3")
	comm.WriteFileTextP(fs, "Test_ListFilesWithExt_happy/4_empty.hosts.txt", "")
	comm.WriteFileTextP(fs, "Test_ListFilesWithExt_happy/.hosts.txt", "5")

	a.Equal(map[string]string{
		"1": filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
	}, comm.ListSuffixedFilesP(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", true))

	a.Equal(map[string]string{
		"1":       filepath.Join("Test_ListFilesWithExt_happy", "1.hosts.txt"),
		"4_empty": filepath.Join("Test_ListFilesWithExt_happy", "4_empty.hosts.txt"),
	}, comm.ListSuffixedFilesP(fs, "Test_ListFilesWithExt_happy", ".hosts.txt", false))
}

func Test_ExtractTitle_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("abc", comm.ExtractTitle("/Test_ExtractTitle_happy/abc.xyz"))
	a.Equal("abc", comm.ExtractTitle("/Test_ExtractTitle_happy/abc"))
	a.Equal("", comm.ExtractTitle("/Test_ExtractTitle_happy/.xyz"))
}

func Test_FileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.False(comm.FileExists(fs, "/f.txt"))
	comm.WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	a.True(comm.FileExists(fs, "/f.txt"))

	comm.MkdirP(fs, "/d")
	a.Panics(func() { comm.FileExistsP(fs, "/d") }, "expect /d be file, but it is directory")
}

func Test_EnsureFileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { comm.EnsureFileExistsP(fs, "/F.txt") }, "file not found: %s")
	comm.WriteFileTextIfNotFoundP(fs, "/F.txt", "blah")
	comm.EnsureFileExistsP(fs, "/F.txt")

	comm.MkdirP(fs, "/D")
	a.Panics(func() { comm.EnsureFileExistsP(fs, "/D") }, "expect /D be file, but it is directory")
}

func Test_DirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.False(comm.DirExists(fs, "/d"))
	comm.MkdirP(fs, "/d")
	a.True(comm.DirExists(fs, "/d"))

	comm.WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { comm.DirExistsP(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_EnsureDirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { comm.EnsureDirExistsP(fs, "/D") }, "directory not found: %s")
	comm.MkdirP(fs, "/D")
	comm.EnsureDirExistsP(fs, "/D")

	comm.WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { comm.EnsureDirExistsP(fs, "/f.txt") }, "expect /f.txt be directory, but it is file")
}

func Test_RemoveFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.RemoveFileP(fs, "/f.html")
	comm.WriteFileTextIfNotFoundP(fs, "/f.html", "<html></html>")
	a.True(comm.FileExists(fs, "/f.html"))
	comm.RemoveFileP(fs, "/f.html")
	a.False(comm.FileExists(fs, "/f.html"))

	comm.MkdirP(fs, "/D")
	a.Panics(func() { comm.RemoveFileP(fs, "/D") }, "expect /D be file, but it is directory")
	a.True(comm.DirExists(fs, "/D"))
}

func Test_RemoveDir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.RemoveDirP(fs, "/d")
	comm.MkdirP(fs, "/d")
	comm.RemoveDirP(fs, "/d")

	comm.WriteFileTextIfNotFoundP(fs, "/f.html", "<html></html>")
	a.Panics(func() { comm.RemoveDirP(fs, "/f.html") }, "expect /f.html be directory, but it is file")
}

func Test_WriteFileIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(comm.WriteFileIfNotFoundP(fs, "/f.txt", []byte("hello")))
	a.Equal("hello", comm.ReadFileTextP(fs, "/f.txt"))

	a.False(comm.WriteFileIfNotFound(fs, "/f.txt", []byte("updated")))
	a.Equal("hello", comm.ReadFileTextP(fs, "/f.txt"))
}

func Test_WriteFileTextIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(comm.WriteFileTextIfNotFoundP(fs, "/f.txt", "hello"))
	a.Equal("hello", comm.ReadFileTextP(fs, "/f.txt"))

	a.False(comm.WriteFileTextIfNotFoundP(fs, "/f.txt", "updated"))
	a.Equal("hello", comm.ReadFileTextP(fs, "/f.txt"))
}

func Test_TempFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	actual := comm.TempFileP(fs, "xyz")
	a.True(strings.Contains(actual, "xyz"))
	a.NotEqual("xyz", actual)
}

func Test_ExpandHomePath_happy(t *testing.T) {
	a := require.New(t)
	a.Equal("none", comm.ExpandHomePathP("none"))

	a.Equal(comm.UserHomeDirP(), comm.ExpandHomePathP("~"))
}

func Test_Stat_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "/f.txt", "content")

	// Stat existing file
	fi, err := comm.Stat(fs, "/f.txt", false)
	a.NoError(err)
	a.NotNil(fi)
	a.False(fi.IsDir())

	// Stat non-existing file without ensureExists
	fi, err = comm.Stat(fs, "/nonexistent.txt", false)
	a.NoError(err)
	a.Nil(fi)

	// Stat non-existing file with ensureExists
	fi, err = comm.Stat(fs, "/nonexistent.txt", true)
	a.Error(err)
	a.Nil(fi)

	// StatP
	fi = comm.StatP(fs, "/f.txt", false)
	a.NotNil(fi)
}

func Test_EnsureFileNotExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// File doesn't exist - should pass
	err := comm.EnsureFileNotExists(fs, "/nonexistent.txt")
	a.NoError(err)

	// Create file
	comm.WriteFileTextP(fs, "/f.txt", "content")

	// File exists - should error
	err = comm.EnsureFileNotExists(fs, "/f.txt")
	a.Error(err)

	// EnsureFileNotExistsP with non-existing file
	comm.EnsureFileNotExistsP(fs, "/another_nonexistent.txt")

	// EnsureFileNotExistsP with existing file should panic
	a.Panics(func() {
		comm.EnsureFileNotExistsP(fs, "/f.txt")
	})
}

func Test_EnsureDirNotExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// Dir doesn't exist - should pass
	err := comm.EnsureDirNotExists(fs, "/nonexistent_dir")
	a.NoError(err)

	// Create dir
	comm.MkdirP(fs, "/mydir")

	// Dir exists - should error
	err = comm.EnsureDirNotExists(fs, "/mydir")
	a.Error(err)

	// EnsureDirNotExistsP with non-existing dir
	comm.EnsureDirNotExistsP(fs, "/another_nonexistent_dir")

	// EnsureDirNotExistsP with existing dir should panic
	a.Panics(func() {
		comm.EnsureDirNotExistsP(fs, "/mydir")
	})
}

func Test_ReadFileBytes_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	content := []byte("hello bytes")
	comm.WriteFileP(fs, "/bytes.bin", content)

	// ReadFileBytes
	result, err := comm.ReadFileBytes(fs, "/bytes.bin")
	a.NoError(err)
	a.Equal(content, result)

	// ReadFileBytesP
	result = comm.ReadFileBytesP(fs, "/bytes.bin")
	a.Equal(content, result)

	// Non-existing file
	_, err = comm.ReadFileBytes(fs, "/nonexistent.bin")
	a.Error(err)
}

func Test_WriteFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	content := []byte("test content")

	// WriteFile
	err := comm.WriteFile(fs, "/test.txt", content)
	a.NoError(err)

	result := comm.ReadFileBytesP(fs, "/test.txt")
	a.Equal(content, result)

	// WriteFileP
	comm.WriteFileP(fs, "/test2.txt", []byte("content2"))
	result = comm.ReadFileBytesP(fs, "/test2.txt")
	a.Equal([]byte("content2"), result)
}

func Test_WriteFileText_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// WriteFileText
	err := comm.WriteFileText(fs, "/text.txt", "hello text")
	a.NoError(err)

	result := comm.ReadFileTextP(fs, "/text.txt")
	a.Equal("hello text", result)

	// WriteFileTextP
	comm.WriteFileTextP(fs, "/text2.txt", "text2")
	result = comm.ReadFileTextP(fs, "/text2.txt")
	a.Equal("text2", result)
}

func Test_Mkdir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// Mkdir
	err := comm.Mkdir(fs, "/new_dir")
	a.NoError(err)
	a.True(comm.DirExists(fs, "/new_dir"))

	// MkdirP
	comm.MkdirP(fs, "/another_dir")
	a.True(comm.DirExists(fs, "/another_dir"))

	// Nested mkdir
	comm.MkdirP(fs, "/parent/child/grandchild")
	a.True(comm.DirExists(fs, "/parent/child/grandchild"))
}

func Test_Mkdir4File_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// Mkdir4File
	err := comm.Mkdir4File(fs, "/dir_for_file/nested/file.txt")
	a.NoError(err)
	a.True(comm.DirExists(fs, "/dir_for_file/nested"))

	// Mkdir4FileP
	comm.Mkdir4FileP(fs, "/another_dir/deep/path/file.txt")
	a.True(comm.DirExists(fs, "/another_dir/deep/path"))
}

func Test_TempTextFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// TempTextFile
	path, err := comm.TempTextFile(fs, "prefix", "content")
	a.NoError(err)
	a.True(strings.Contains(path, "prefix"))

	result := comm.ReadFileTextP(fs, path)
	a.Equal("content", result)

	// TempTextFileP
	path = comm.TempTextFileP(fs, "prefix2", "content2")
	a.True(strings.Contains(path, "prefix2"))
}

func Test_UserHomeDir_happy(t *testing.T) {
	a := require.New(t)

	// UserHomeDir
	dir, err := comm.UserHomeDir()
	a.NoError(err)
	a.NotEmpty(dir)

	// UserHomeDirP
	dirP := comm.UserHomeDirP()
	a.Equal(dir, dirP)
}

func Test_Rename_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// Rename non-existing file should error
	err := comm.Rename(fs, "/nonexistent.txt", "/dest.txt")
	a.Error(err)
}

func Test_FileExists_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// FileExists returns error-wrapped value
	exists, err := comm.FileExists(fs, "/nonexistent.txt")
	a.NoError(err)
	a.False(exists)
}

func Test_DirExists_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// DirExists returns false for non-existing dir
	exists, _ := comm.DirExists(fs, "/nonexistent_dir")
	a.False(exists)
}

func Test_EnsureDirExists_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// EnsureDirExists returns error for non-existing dir
	err := comm.EnsureDirExists(fs, "/nonexistent_dir")
	a.Error(err)
}

func Test_WriteFileLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// WriteFileLines
	err := comm.WriteFileLines(fs, "/lines.txt", "a", "b", "c")
	a.NoError(err)

	lines := comm.ReadFileLinesP(fs, "/lines.txt")
	a.Equal([]string{"a", "b", "c"}, lines)
}

func Test_ReadFileText_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// ReadFileText on non-existing file
	_, err := comm.ReadFileText(fs, "/nonexistent.txt")
	a.Error(err)
}

func Test_RemoveFile_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// RemoveFile on non-existing file - should not error
	err := comm.RemoveFile(fs, "/nonexistent.txt")
	a.NoError(err)
}

func Test_RemoveDir_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	// RemoveDir on non-existing dir - should not error
	err := comm.RemoveDir(fs, "/nonexistent_dir")
	a.NoError(err)
}
