package test

import (
	"path/filepath"
	"testing"

	"github.com/fastgh/go-comm/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_WorkDir_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("defaultDir", comm.WorkDir("http://test", "defaultDir"))
	a.Equal("defaultDir", comm.WorkDir("file://hello.txt", "defaultDir"))
	a.Equal("defaultDir", comm.WorkDir("hello.txt", "defaultDir"))
	a.Equal(filepath.Join("defaultDir", "home"), comm.WorkDir("home/hello.txt", "defaultDir"))

	if comm.IsWindows() {
		a.Equal("c:\\home", comm.WorkDir("file://c:/home/hello.txt", "defaultDir"))
		a.Equal("c:\\home", comm.WorkDir("c:/home/hello.txt", "defaultDir"))

		a.Equal("C:\\", comm.WorkDir("C:\\hello.txt", "defaultDir"))
	} else {
		a.Equal("/home", comm.WorkDir("file:///home/hello.txt", "defaultDir"))
		a.Equal("/home", comm.WorkDir("/home/hello.txt", "defaultDir"))
	}
}

func Test_IsRemote_happy(t *testing.T) {
	a := require.New(t)

	a.True(comm.IsRemote("http://test.local"))
	a.True(comm.IsRemote("HTTP://test.local"))
	a.True(comm.IsRemote("https://test.local"))
	a.True(comm.IsRemote("HTTPS://test.local"))
	a.True(comm.IsRemote("HTTPs://test.local"))

	a.True(comm.IsRemote("ftp://test.local"))
	a.True(comm.IsRemote("FTP://test.local"))
	a.True(comm.IsRemote("ftps://test.local"))
	a.True(comm.IsRemote("FTPS://test.local"))
	a.True(comm.IsRemote("FTPs://test.local"))

	a.True(comm.IsRemote("sftp://test.local"))
	a.True(comm.IsRemote("SFTP://test.local"))

	a.True(comm.IsRemote("s3://test.local"))
	a.True(comm.IsRemote("S3://test.local"))
}

func Test_NewFile(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	fRemote := comm.NewFileP(nil, "https://google.com", nil, 0)
	_, isRemoteFile := fRemote.(comm.RemoteFile)
	a.True(isRemoteFile)

	fLocal := comm.NewFileP(afs, "file://test.txt", nil, 0)
	_, isLocalFile := fLocal.(comm.AferoFile)
	a.True(isLocalFile)

	fLocal = comm.NewFileP(afs, "test.txt", nil, 0)
	_, isLocalFile = fLocal.(comm.AferoFile)
	a.True(isLocalFile)
}

func Test_ShortDescription_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("", comm.ShortDescription(""))

	a.Equal("AB/12345678.hosts", comm.ShortDescription("AB/12345678.hosts"))
	a.Equal("ABC...12345678.hosts", comm.ShortDescription("ABC/12345678.hosts"))

	a.Equal("https://AB/12345678.hosts", comm.ShortDescription("https://AB/12345678.hosts"))
	a.Equal("SFTP://ABC...12345678.hosts", comm.ShortDescription("SFTP://ABC/12345678.hosts"))
}

func Test_DownloadText_happy(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	comm.WriteFileTextP(afs, "test.txt", "Test_DownloadText_happy")

	actual := comm.DownloadTextP(nil, "", afs, "test.txt", nil, 0)
	a.Equal("Test_DownloadText_happy", actual)
}

func Test_DownloadText_Remote(t *testing.T) {
	a := require.New(t)
	afs := afero.NewMemMapFs()

	fallbackDir := "/fallback"
	comm.MkdirP(afs, fallbackDir)

	url := "https://mirror.sjtu.edu.cn/debian/README.mirrors.txt"
	a.False(comm.HasFallbackFile(fallbackDir, afs, url))

	actual := comm.DownloadTextP(nil, fallbackDir, afs, url, nil, 0)
	a.Equal("The list of Debian mirror sites is available here: https://www.debian.org/mirror/list\n", actual)

	a.True(comm.HasFallbackFile(fallbackDir, afs, url))
}

func Test_MapFromYamlFileP_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "test.yaml", `k: v`)

	configMap := comm.MapFromYamlFileP(fs, "test.yaml", false)

	a.Len(configMap, 1)
	a.Equal("v", configMap["k"])
}

func Test_MapFromJsonFileP_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	comm.WriteFileTextP(fs, "test.json", `{"k": "v"}`)

	configMap := comm.MapFromJsonFileP(fs, "test.json", false)

	a.Len(configMap, 1)
	a.Equal("v", configMap["k"])
}
