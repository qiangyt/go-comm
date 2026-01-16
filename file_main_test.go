package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	// Test local file
	file, err := NewFile(fs, "/path/to/file.txt", nil, 0)
	a.NoError(err)
	a.NotNil(file)
	a.Equal("file.txt", file.Name())
}

func TestNewFileP_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	file := NewFileP(fs, "/path/to/file.txt", nil, 0)
	a.NotNil(file)
	a.Equal("file.txt", file.Name())
}

func TestNewFileP_panicsOnError(t *testing.T) {
	fs := afero.NewMemMapFs()

	defer func() {
		if r := recover(); r == nil {
			t.Error("NewFileP should panic on error")
		}
	}()

	// Use an invalid remote URL that contains spaces (causes URL parse error)
	NewFileP(fs, "http://example com/file.txt", nil, 0)
}

func TestIsRemote_true(t *testing.T) {
	a := require.New(t)

	a.True(IsRemote("http://example.com/file.txt"))
	a.True(IsRemote("HTTP://example.com/file.txt"))
	a.True(IsRemote("https://example.com/file.txt"))
	a.True(IsRemote("HTTPS://example.com/file.txt"))
	a.True(IsRemote("ftp://example.com/file.txt"))
	a.True(IsRemote("ftps://example.com/file.txt"))
	a.True(IsRemote("sftp://example.com/file.txt"))
	a.True(IsRemote("s3://bucket/file.txt"))
}

func TestIsRemote_false(t *testing.T) {
	a := require.New(t)

	a.False(IsRemote("/path/to/file.txt"))
	a.False(IsRemote("./file.txt"))
	a.False(IsRemote("file.txt"))
	a.False(IsRemote("file://path/to/file.txt"))
}

func TestWorkDir_remote(t *testing.T) {
	a := require.New(t)

	result := WorkDir("http://example.com/file.txt", "/default/dir")
	a.Equal("/default/dir", result)
}

func TestWorkDir_fileProtocol(t *testing.T) {
	a := require.New(t)

	result := WorkDir("file:///path/to/file.txt", "/default/dir")
	// On Windows, the result will be a Windows path, on Unix it will be a Unix path
	a.NotEqual("/default/dir", result)
	a.NotEmpty(result)
}

func TestWorkDir_relativePath(t *testing.T) {
	a := require.New(t)

	result := WorkDir("relative/path/file.txt", "/default/dir")
	a.NotEqual("/default/dir", result)
	a.Contains(result, "relative")
}

func TestWorkDir_absolutePath(t *testing.T) {
	a := require.New(t)

	result := WorkDir("/absolute/path/file.txt", "/default/dir")
	// On Windows, absolute paths will be treated differently
	// Just check that result is not empty
	a.NotEmpty(result)
}

func TestWorkDir_dot(t *testing.T) {
	a := require.New(t)

	result := WorkDir("file.txt", "/default/dir")
	a.Equal("/default/dir", result)
}

func TestShortDescription_shortUrl(t *testing.T) {
	a := require.New(t)

	result := ShortDescription("http://example.com/file.txt")
	a.NotEmpty(result)
}

func TestShortDescription_longUrl(t *testing.T) {
	a := require.New(t)

	// A URL longer than 3 + 8 + 1 + 5 = 17 chars after protocol
	longUrl := "http://example.com/very/long/path/12345678.hosts.txt"
	result := ShortDescription(longUrl)
	a.NotEmpty(result)
	a.Contains(result, "...")
}

func TestShortDescription_noProtocol(t *testing.T) {
	a := require.New(t)

	result := ShortDescription("/path/to/file.txt")
	a.Equal("/path/to/file.txt", result)
}

func TestDownloadBytesP_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.txt", "test content")

	result := DownloadBytesP(nil, "", fs, "/test.txt", nil, 0)
	a.NotEmpty(result)
	a.Equal([]byte("test content"), result)
}

func TestDownloadBytes_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.txt", "test content")

	result, err := DownloadBytes(nil, "", fs, "/test.txt", nil, 0)
	a.NoError(err)
	a.NotEmpty(result)
	a.Equal([]byte("test content"), result)
}

func TestDownloadTextP_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.txt", "test content")

	result := DownloadTextP(nil, "", fs, "/test.txt", nil, 0)
	a.NotEmpty(result)
	a.Equal("test content", result)
}

func TestDownloadText_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.txt", "test content")

	result, err := DownloadText(nil, "", fs, "/test.txt", nil, 0)
	a.NoError(err)
	a.NotEmpty(result)
	a.Equal("test content", result)
}

func TestMapFromYamlP_happy(t *testing.T) {
	a := require.New(t)

	yamlText := "key1: value1\nkey2: value2"
	result := MapFromYamlP(yamlText, false)
	a.NotNil(result)
	a.Equal("value1", result["key1"])
	a.Equal("value2", result["key2"])
}

func TestMapFromJson_happy(t *testing.T) {
	a := require.New(t)

	jsonText := `{"key1": "value1", "key2": "value2"}`
	result, err := MapFromJson(jsonText, false)
	a.NoError(err)
	a.NotNil(result)
	a.Equal("value1", result["key1"])
	a.Equal("value2", result["key2"])
}

func TestMapFromJsonP_happy(t *testing.T) {
	a := require.New(t)

	jsonText := `{"key1": "value1", "key2": "value2"}`
	result := MapFromJsonP(jsonText, false)
	a.NotNil(result)
	a.Equal("value1", result["key1"])
	a.Equal("value2", result["key2"])
}

func TestFromYamlFileP_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	yamlContent := "name: test\nage: 30"
	WriteFileTextP(fs, "/test.yaml", yamlContent)

	type Config struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	var result Config
	FromYamlFileP(fs, "/test.yaml", false, &result)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestFromYamlP_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	yamlText := "name: test\nage: 30"
	var result Config
	FromYamlP(yamlText, false, &result)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestFromJsonFileP_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	jsonContent := `{"name": "test", "age": 30}`
	WriteFileTextP(fs, "/test.json", jsonContent)

	type Config struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var result Config
	FromJsonFileP(fs, "/test.json", false, &result)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestFromJsonP_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonText := `{"name": "test", "age": 30}`
	var result Config
	FromJsonP(jsonText, false, &result)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}
