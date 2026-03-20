package qfile

import (
	"path/filepath"
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_WorkDir_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("defaultDir", WorkDir("http://test", "defaultDir"))
	a.Equal("defaultDir", WorkDir("file://hello.txt", "defaultDir"))
	a.Equal("defaultDir", WorkDir("hello.txt", "defaultDir"))
	a.Equal(filepath.Join("defaultDir", "home"), WorkDir("home/hello.txt", "defaultDir"))

	if comm.IsWindows() {
		a.Equal("c:\\home", WorkDir("file://c:/home/hello.txt", "defaultDir"))
		a.Equal("c:\\home", WorkDir("c:/home/hello.txt", "defaultDir"))

		a.Equal("C:\\", WorkDir("C:\\hello.txt", "defaultDir"))
	} else {
		a.Equal("/home", WorkDir("file:///home/hello.txt", "defaultDir"))
		a.Equal("/home", WorkDir("/home/hello.txt", "defaultDir"))
	}
}

func Test_IsRemote_happy(t *testing.T) {
	a := require.New(t)

	a.True(IsRemote("http://test.local"))
	a.True(IsRemote("HTTP://test.local"))
	a.True(IsRemote("https://test.local"))
	a.True(IsRemote("HTTPS://test.local"))
	a.True(IsRemote("HTTPs://test.local"))

	a.True(IsRemote("ftp://test.local"))
	a.True(IsRemote("FTP://test.local"))
	a.True(IsRemote("ftps://test.local"))
	a.True(IsRemote("FTPS://test.local"))
	a.True(IsRemote("FTPs://test.local"))

	a.True(IsRemote("sftp://test.local"))
	a.True(IsRemote("SFTP://test.local"))

	a.True(IsRemote("s3://test.local"))
	a.True(IsRemote("S3://test.local"))
}

func Test_NewFile(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	fRemote := NewFileP(nil, "https://google.com", nil, 0)
	_, isRemoteFile := fRemote.(RemoteFile)
	a.True(isRemoteFile)

	fLocal := NewFileP(afs, "file://test.txt", nil, 0)
	_, isLocalFile := fLocal.(AferoFile)
	a.True(isLocalFile)

	fLocal = NewFileP(afs, "test.txt", nil, 0)
	_, isLocalFile = fLocal.(AferoFile)
	a.True(isLocalFile)
}

func Test_ShortDescription_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("", ShortDescription(""))

	a.Equal("AB/12345678.hosts", ShortDescription("AB/12345678.hosts"))
	a.Equal("ABC...12345678.hosts", ShortDescription("ABC/12345678.hosts"))

	a.Equal("https://AB/12345678.hosts", ShortDescription("https://AB/12345678.hosts"))
	a.Equal("SFTP://ABC...12345678.hosts", ShortDescription("SFTP://ABC/12345678.hosts"))
}

func Test_DownloadText_happy(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	WriteFileTextP(afs, "test.txt", "Test_DownloadText_happy")

	actual := DownloadTextP(nil, "", afs, "test.txt", nil, 0)
	a.Equal("Test_DownloadText_happy", actual)
}

func Test_DownloadText_Remote(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	fallbackDir := "/fallback"
	MkdirP(afs, fallbackDir)

	url := "https://mirror.sjtu.edu.cn/debian/README.mirrors.txt"
	a.False(HasFallbackFile(fallbackDir, afs, url))

	actual := DownloadTextP(nil, fallbackDir, afs, url, nil, 0)
	a.Equal("The list of Debian mirror sites is available here: https://www.debian.org/mirror/list\n", actual)

	a.True(HasFallbackFile(fallbackDir, afs, url))
}
