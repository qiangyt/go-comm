package comm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRemoteFile_happy(t *testing.T) {
	a := require.New(t)

	remoteFile, err := NewRemoteFile("http://example.com/file.txt", nil, 30*time.Second)
	a.NoError(err)
	a.NotNil(remoteFile)
	a.Equal("http://example.com/file.txt", remoteFile.Url())
	a.Equal("http", remoteFile.Protocol())
	a.Equal(30*time.Second, remoteFile.Timeout())
}

func TestNewRemoteFileP_happy(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, 30*time.Second)
	a.NotNil(remoteFile)
	a.Equal("http://example.com/file.txt", remoteFile.Url())
}

func TestNewRemoteFileP_invalidUrl(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewRemoteFileP should panic on invalid URL")
		}
	}()

	NewRemoteFileP(":", nil, 30*time.Second)
}

func TestRemoteFile_Name(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, 30*time.Second)
	a.NotEmpty(remoteFile.Name())
}

func TestRemoteFile_Dir(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/path/to/file.txt", nil, 30*time.Second)
	a.NotEmpty(remoteFile.Dir())
}

func TestRemoteFile_Url(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, 30*time.Second)
	a.Equal("http://example.com/file.txt", remoteFile.Url())
}

func TestRemoteFile_Protocol(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, 30*time.Second)
	a.Equal("http", remoteFile.Protocol())

	remoteFile = NewRemoteFileP("https://example.com/file.txt", nil, 30*time.Second)
	a.Equal("https", remoteFile.Protocol())
}

func TestRemoteFile_URL(t *testing.T) {
	a := require.New(t)

	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, 30*time.Second)
	parsedURL := remoteFile.URL()
	a.NotNil(parsedURL)
	a.Equal("http", parsedURL.Scheme)
	a.Equal("example.com", parsedURL.Host)
}

func TestRemoteFile_Credentials(t *testing.T) {
	a := require.New(t)

	creds := &CredentialsT{
		User: "testuser",
	}
	remoteFile := NewRemoteFileP("http://example.com/file.txt", creds, 30*time.Second)

	fileCreds := remoteFile.Credentials()
	a.NotNil(fileCreds)
	a.Equal("testuser", fileCreds.User)
}

func TestRemoteFile_Timeout(t *testing.T) {
	a := require.New(t)

	timeout := 45 * time.Second
	remoteFile := NewRemoteFileP("http://example.com/file.txt", nil, timeout)
	a.Equal(timeout, remoteFile.Timeout())
}

func TestRemoteFile_Download_error(t *testing.T) {
	a := require.New(t)

	// This will fail since the URL doesn't exist
	remoteFile := NewRemoteFileP("http://this-url-does-not-exist-12345.com/file.txt", nil, 5*time.Second)
	_, err := remoteFile.Download()
	a.Error(err)
}

func TestRemoteFile_DownloadP_panicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("DownloadP should panic on error")
		}
	}()

	remoteFile := NewRemoteFileP("http://this-url-does-not-exist-12345.com/file.txt", nil, 5*time.Second)
	remoteFile.DownloadP()
}
