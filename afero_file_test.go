package comm

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestNewAferoFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	credentials := &CredentialsT{}
	timeout := 30 * time.Second

	// Test with simple path
	af, err := NewAferoFile(afs, "/path/to/file.txt", credentials, timeout)
	if err != nil {
		t.Fatalf("NewAferoFile error: %v", err)
	}

	if af == nil {
		t.Fatal("NewAferoFile returned nil")
	}

	if af.Name() != "file.txt" {
		t.Errorf("Expected name 'file.txt', got '%s'", af.Name())
	}

	// Dir() returns platform-specific path, just check it's not empty
	if af.Dir() == "" {
		t.Error("Dir() should not be empty")
	}

	if af.Protocol() != "file" {
		t.Errorf("Expected protocol 'file', got '%s'", af.Protocol())
	}

	if af.Timeout() != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, af.Timeout())
	}
}

func TestNewAferoFile_WithFileProtocol(t *testing.T) {
	afs := afero.NewMemMapFs()
	credentials := &CredentialsT{}

	af, err := NewAferoFile(afs, "file:///path/to/file.txt", credentials, 0)
	if err != nil {
		t.Fatalf("NewAferoFile error: %v", err)
	}

	if af.Protocol() != "file" {
		t.Errorf("Expected protocol 'file', got '%s'", af.Protocol())
	}

	if af.Name() != "file.txt" {
		t.Errorf("Expected name 'file.txt', got '%s'", af.Name())
	}
}

func TestNewAferoFile_InvalidUrl(t *testing.T) {
	afs := afero.NewMemMapFs()

	// Test with invalid URL
	_, err := NewAferoFile(afs, ":invalid", nil, 0)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestNewAferoFileP(t *testing.T) {
	afs := afero.NewMemMapFs()

	// Test panic on error
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewAferoFileP should panic on error")
		}
	}()

	NewAferoFileP(afs, ":invalid", nil, 0)
}

func TestAferoFile_Fs(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	if af.Fs() != afs {
		t.Error("Fs() should return the same filesystem")
	}
}

func TestAferoFile_Name(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	if af.Name() != "file.txt" {
		t.Errorf("Expected 'file.txt', got '%s'", af.Name())
	}
}

func TestAferoFile_Dir(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	// Dir() returns platform-specific path
	dir := af.Dir()
	if dir == "" {
		t.Error("Dir() should not be empty")
	}
}

func TestAferoFile_Url(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	url := af.Url()
	if !strings.Contains(url, "file://") {
		t.Errorf("Expected URL to contain 'file://', got '%s'", url)
	}
	if !strings.Contains(url, "file.txt") {
		t.Errorf("Expected URL to contain 'file.txt', got '%s'", url)
	}
}

func TestAferoFile_Protocol(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	if af.Protocol() != "file" {
		t.Errorf("Expected protocol 'file', got '%s'", af.Protocol())
	}
}

func TestAferoFile_URL(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	parsedURL := af.URL()
	if parsedURL == nil {
		t.Fatal("URL() returned nil")
	}

	if parsedURL.Scheme != "file" {
		t.Errorf("Expected URL scheme 'file', got '%s'", parsedURL.Scheme)
	}
}

func TestAferoFile_Credentials(t *testing.T) {
	afs := afero.NewMemMapFs()
	credentials := &CredentialsT{}
	af, _ := NewAferoFile(afs, "/path/to/file.txt", credentials, 0)

	if af.Credentials() != credentials {
		t.Error("Credentials() should return the same credentials")
	}
}

func TestAferoFile_Timeout(t *testing.T) {
	afs := afero.NewMemMapFs()
	timeout := 60 * time.Second
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, timeout)

	if af.Timeout() != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, af.Timeout())
	}
}

func TestAferoFile_DownloadP(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	content := af.DownloadP()
	if content == nil {
		t.Fatal("DownloadP() returned nil")
	}

	if content.Name != "file.txt" {
		t.Errorf("Expected content name 'file.txt', got '%s'", content.Name)
	}

	if !strings.Contains(content.Path, "file.txt") {
		t.Errorf("Expected content path to contain 'file.txt', got '%s'", content.Path)
	}

	if content.Blob == nil {
		t.Error("Content Blob should not be nil")
	}
}

func TestAferoFile_Download(t *testing.T) {
	afs := afero.NewMemMapFs()
	af, _ := NewAferoFile(afs, "/path/to/file.txt", nil, 0)

	content, err := af.Download()
	if err != nil {
		t.Fatalf("Download() error: %v", err)
	}

	if content == nil {
		t.Fatal("Download() returned nil content")
	}

	if content.Name != "file.txt" {
		t.Errorf("Expected content name 'file.txt', got '%s'", content.Name)
	}
}

func TestNewAferoBlob(t *testing.T) {
	afs := afero.NewMemMapFs()
	ab := NewAferoBlob(afs, "/path/to/file.txt")

	if ab == nil {
		t.Fatal("NewAferoBlob returned nil")
	}

	if ab.Path() != "/path/to/file.txt" {
		t.Errorf("Expected path '/path/to/file.txt', got '%s'", ab.Path())
	}

	if ab.Fs() != afs {
		t.Error("Fs() should return the same filesystem")
	}
}

func TestAferoBlob_Path(t *testing.T) {
	afs := afero.NewMemMapFs()
	ab := NewAferoBlob(afs, "/path/to/file.txt")

	if ab.Path() != "/path/to/file.txt" {
		t.Errorf("Expected '/path/to/file.txt', got '%s'", ab.Path())
	}
}

func TestAferoBlob_Fs(t *testing.T) {
	afs := afero.NewMemMapFs()
	ab := NewAferoBlob(afs, "/path/to/file.txt")

	if ab.Fs() != afs {
		t.Error("Fs() should return the same filesystem")
	}
}

func TestAferoBlob_Read(t *testing.T) {
	afs := afero.NewMemMapFs()
	afs.Create("/path/to/file.txt")

	ab := NewAferoBlob(afs, "/path/to/file.txt")
	buf := make([]byte, 10)

	n, err := ab.Read(buf)
	// Empty file should return 0 bytes and EOF
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}

	if n != 0 {
		t.Errorf("Expected to read 0 bytes, got %d", n)
	}
}

func TestAferoBlob_Read_NotFound(t *testing.T) {
	afs := afero.NewMemMapFs()
	ab := NewAferoBlob(afs, "/nonexistent/file.txt")
	buf := make([]byte, 10)

	_, err := ab.Read(buf)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestAferoBlob_Close(t *testing.T) {
	afs := afero.NewMemMapFs()
	afs.Create("/path/to/file.txt")

	ab := NewAferoBlob(afs, "/path/to/file.txt")

	// Close without reading should not error
	err := ab.Close()
	if err != nil {
		t.Fatalf("Close error: %v", err)
	}
}

func TestAferoBlob_Close_AfterRead(t *testing.T) {
	afs := afero.NewMemMapFs()
	afero.WriteFile(afs, "/path/to/file.txt", []byte("hello"), 0644)

	ab := NewAferoBlob(afs, "/path/to/file.txt")
	buf := make([]byte, 10)
	ab.Read(buf)

	err := ab.Close()
	if err != nil {
		t.Fatalf("Close error: %v", err)
	}

	// File should be removed
	exists, _ := afero.Exists(afs, "/path/to/file.txt")
	if exists {
		t.Error("File should be removed after Close()")
	}
}

func TestAferoFile_Url_WithHttpProtocol(t *testing.T) {
	afs := afero.NewMemMapFs()
	// Note: NewAferoFile always treats input as file paths
	af, _ := NewAferoFile(afs, "http://example.com/path/file.txt", nil, 0)

	// Name is extracted from the path
	if af.Name() != "file.txt" {
		t.Errorf("Expected name \"file.txt\", got \"%s\"", af.Name())
	}

	// Protocol is "file" since we are using AferoFile
	if af.Protocol() != "file" {
		t.Errorf("Expected protocol \"file\", got \"%s\"", af.Protocol())
	}
}
