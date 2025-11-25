package comm

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCopyFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "/TestCopyFile_happy/c1")
	WriteFileTextP(fs, "/TestCopyFile_happy/c1/src.txt", "hello")

	MkdirP(fs, "/TestCopyFile_happy/c2")
	_, err := CopyFile(fs, "/TestCopyFile_happy/c1/src.txt", "/TestCopyFile_happy/c2/dest.txt")
	a.NoError(err)

	actual := ReadFileTextP(fs, "/TestCopyFile_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func TestCopyFileP_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "/TestCopyFileP_happy/c1")
	WriteFileTextP(fs, "/TestCopyFileP_happy/c1/src.txt", "hello")

	MkdirP(fs, "/TestCopyFileP_happy/c2")
	n := CopyFileP(fs, "/TestCopyFileP_happy/c1/src.txt", "/TestCopyFileP_happy/c2/dest.txt")
	a.Equal(int64(5), n)

	actual := ReadFileTextP(fs, "/TestCopyFileP_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func TestCopyFile_SourceFileNotFound(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "/TestCopyFile_SourceFileNotFound/c1")
	MkdirP(fs, "/TestCopyFile_SourceFileNotFound/c2")

	a.Panics(func() {
		CopyFileP(fs, "/TestCopyFile_SourceFileNotFound/c1/src.txt", "/TestCopyFile_SourceFileNotFound/c2/dest.txt")
	})
}

func TestRename_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "/TestRename_happy/c1")
	WriteFileTextP(fs, "/TestRename_happy/c1/src.txt", "hello")

	MkdirP(fs, "/TestRename_happy/c2")
	RenameP(fs, "/TestRename_happy/c1/src.txt", "/TestRename_happy/c2/dest.txt")

	actual := ReadFileTextP(fs, "/TestRename_happy/c2/dest.txt")
	a.Equal("hello", actual)
}

func TestRename_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := Rename(fs, "/nonexistent.txt", "/dest.txt")
	a.Error(err)
}

func TestStat_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	WriteFileTextP(fs, "/f.txt", "content")

	fi, err := Stat(fs, "/f.txt", false)
	a.NoError(err)
	a.NotNil(fi)
	a.False(fi.IsDir())

	fi, err = Stat(fs, "/nonexistent.txt", false)
	a.NoError(err)
	a.Nil(fi)

	fi, err = Stat(fs, "/nonexistent.txt", true)
	a.Error(err)
	a.Nil(fi)

	fi = StatP(fs, "/f.txt", false)
	a.NotNil(fi)
}

func TestReadFileLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "/TestReadAsLines_happy")
	WriteFileLinesP(fs, "/TestReadAsLines_happy/f.txt", "line 1", "line 2")

	actual := ReadFileLinesP(fs, "/TestReadAsLines_happy/f.txt")
	a.Equal([]string{"line 1", "line 2"}, actual)
}

func TestListSuffixedFiles_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	MkdirP(fs, "TestListFilesWithExt_happy")
	MkdirP(fs, "TestListFilesWithExt_happy/.")
	MkdirP(fs, "TestListFilesWithExt_happy/..")
	MkdirP(fs, "TestListFilesWithExt_happy/d.hosts.txt")

	WriteFileTextP(fs, "TestListFilesWithExt_happy/1.hosts.txt", "1")
	WriteFileTextP(fs, "TestListFilesWithExt_happy/2.hosts.txt.not", "2")
	WriteFileTextP(fs, "TestListFilesWithExt_happy/3.hosts.not.text", "3")
	WriteFileTextP(fs, "TestListFilesWithExt_happy/4_empty.hosts.txt", "")
	WriteFileTextP(fs, "TestListFilesWithExt_happy/.hosts.txt", "5")

	result := ListSuffixedFilesP(fs, "TestListFilesWithExt_happy", ".hosts.txt", true)
	a.Equal(1, len(result))
	a.Contains(result, "1")

	result = ListSuffixedFilesP(fs, "TestListFilesWithExt_happy", ".hosts.txt", false)
	a.Equal(2, len(result))
	a.Contains(result, "1")
	a.Contains(result, "4_empty")
}

func TestExtractTitle_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("abc", ExtractTitle("/TestExtractTitle_happy/abc.xyz"))
	a.Equal("abc", ExtractTitle("/TestExtractTitle_happy/abc"))
	a.Equal("", ExtractTitle("/TestExtractTitle_happy/.xyz"))
}

func TestFileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	exists, _ := FileExists(fs, "/f.txt")
	a.False(exists)

	WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	exists, _ = FileExists(fs, "/f.txt")
	a.True(exists)

	MkdirP(fs, "/d")
	a.Panics(func() { FileExistsP(fs, "/d") })
}

func TestEnsureFileExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { EnsureFileExistsP(fs, "/F.txt") })
	WriteFileTextIfNotFoundP(fs, "/F.txt", "blah")
	EnsureFileExistsP(fs, "/F.txt")

	MkdirP(fs, "/D")
	a.Panics(func() { EnsureFileExistsP(fs, "/D") })
}

func TestEnsureFileNotExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := EnsureFileNotExists(fs, "/nonexistent.txt")
	a.NoError(err)

	WriteFileTextP(fs, "/f.txt", "content")

	err = EnsureFileNotExists(fs, "/f.txt")
	a.Error(err)

	EnsureFileNotExistsP(fs, "/another_nonexistent.txt")

	a.Panics(func() {
		EnsureFileNotExistsP(fs, "/f.txt")
	})
}

func TestDirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	exists, _ := DirExists(fs, "/d")
	a.False(exists)

	MkdirP(fs, "/d")
	exists, _ = DirExists(fs, "/d")
	a.True(exists)

	WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { DirExistsP(fs, "/f.txt") })
}

func TestEnsureDirExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.Panics(func() { EnsureDirExistsP(fs, "/D") })
	MkdirP(fs, "/D")
	EnsureDirExistsP(fs, "/D")

	WriteFileTextIfNotFoundP(fs, "/f.txt", "blah")
	a.Panics(func() { EnsureDirExistsP(fs, "/f.txt") })
}

func TestEnsureDirNotExists_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := EnsureDirNotExists(fs, "/nonexistent_dir")
	a.NoError(err)

	MkdirP(fs, "/mydir")

	err = EnsureDirNotExists(fs, "/mydir")
	a.Error(err)

	EnsureDirNotExistsP(fs, "/another_nonexistent_dir")

	a.Panics(func() {
		EnsureDirNotExistsP(fs, "/mydir")
	})
}

func TestRemoveFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	RemoveFileP(fs, "/f.html")
	WriteFileTextIfNotFoundP(fs, "/f.html", "<html></html>")
	exists, _ := FileExists(fs, "/f.html")
	a.True(exists)
	RemoveFileP(fs, "/f.html")
	exists, _ = FileExists(fs, "/f.html")
	a.False(exists)

	MkdirP(fs, "/D")
	a.Panics(func() { RemoveFileP(fs, "/D") })
	exists, _ = DirExists(fs, "/D")
	a.True(exists)
}

func TestRemoveDir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	RemoveDirP(fs, "/d")
	MkdirP(fs, "/d")
	RemoveDirP(fs, "/d")

	WriteFileTextIfNotFoundP(fs, "/f.html", "<html></html>")
	a.Panics(func() { RemoveDirP(fs, "/f.html") })
}

func TestWriteFileIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(WriteFileIfNotFoundP(fs, "/f.txt", []byte("hello")))
	a.Equal("hello", ReadFileTextP(fs, "/f.txt"))

	a.False(WriteFileIfNotFound(fs, "/f.txt", []byte("updated")))
	a.Equal("hello", ReadFileTextP(fs, "/f.txt"))
}

func TestWriteFileTextIfNotFound_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	a.True(WriteFileTextIfNotFoundP(fs, "/f.txt", "hello"))
	a.Equal("hello", ReadFileTextP(fs, "/f.txt"))

	a.False(WriteFileTextIfNotFoundP(fs, "/f.txt", "updated"))
	a.Equal("hello", ReadFileTextP(fs, "/f.txt"))
}

func TestReadFileBytes_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	content := []byte("hello bytes")
	WriteFileP(fs, "/bytes.bin", content)

	result, err := ReadFileBytes(fs, "/bytes.bin")
	a.NoError(err)
	a.Equal(content, result)

	result = ReadFileBytesP(fs, "/bytes.bin")
	a.Equal(content, result)

	_, err = ReadFileBytes(fs, "/nonexistent.bin")
	a.Error(err)
}

func TestWriteFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	content := []byte("test content")

	err := WriteFile(fs, "/test.txt", content)
	a.NoError(err)

	result := ReadFileBytesP(fs, "/test.txt")
	a.Equal(content, result)

	WriteFileP(fs, "/test2.txt", []byte("content2"))
	result = ReadFileBytesP(fs, "/test2.txt")
	a.Equal([]byte("content2"), result)
}

func TestWriteFileText_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := WriteFileText(fs, "/text.txt", "hello text")
	a.NoError(err)

	result := ReadFileTextP(fs, "/text.txt")
	a.Equal("hello text", result)

	WriteFileTextP(fs, "/text2.txt", "text2")
	result = ReadFileTextP(fs, "/text2.txt")
	a.Equal("text2", result)
}

func TestMkdir_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := Mkdir(fs, "/new_dir")
	a.NoError(err)
	exists, _ := DirExists(fs, "/new_dir")
	a.True(exists)

	MkdirP(fs, "/another_dir")
	exists, _ = DirExists(fs, "/another_dir")
	a.True(exists)

	MkdirP(fs, "/parent/child/grandchild")
	exists, _ = DirExists(fs, "/parent/child/grandchild")
	a.True(exists)
}

func TestMkdir4File_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := Mkdir4File(fs, "/dir_for_file/nested/file.txt")
	a.NoError(err)
	exists, _ := DirExists(fs, "/dir_for_file/nested")
	a.True(exists)

	Mkdir4FileP(fs, "/another_dir/deep/path/file.txt")
	exists, _ = DirExists(fs, "/another_dir/deep/path")
	a.True(exists)
}

func TestTempFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	actual := TempFileP(fs, "xyz")
	a.True(strings.Contains(actual, "xyz"))
	a.NotEqual("xyz", actual)
}

func TestTempTextFile_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	path, err := TempTextFile(fs, "prefix", "content")
	a.NoError(err)
	a.True(strings.Contains(path, "prefix"))

	result := ReadFileTextP(fs, path)
	a.Equal("content", result)

	path = TempTextFileP(fs, "prefix2", "content2")
	a.True(strings.Contains(path, "prefix2"))
}

func TestExpandHomePath_happy(t *testing.T) {
	a := require.New(t)
	a.Equal("none", ExpandHomePathP("none"))

	a.Equal(UserHomeDirP(), ExpandHomePathP("~"))
}

func TestUserHomeDir_happy(t *testing.T) {
	a := require.New(t)

	dir, err := UserHomeDir()
	a.NoError(err)
	a.NotEmpty(dir)

	dirP := UserHomeDirP()
	a.Equal(dir, dirP)
}

func TestWriteFileLines_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	err := WriteFileLines(fs, "/lines.txt", "a", "b", "c")
	a.NoError(err)

	lines := ReadFileLinesP(fs, "/lines.txt")
	a.Equal([]string{"a", "b", "c"}, lines)
}

func TestReadFileText_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	_, err := ReadFileText(fs, "/nonexistent.txt")
	a.Error(err)
}

func TestReadFileLines_error(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	_, err := ReadFileLines(fs, "/nonexistent.txt")
	a.Error(err)
}
