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
