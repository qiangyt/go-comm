package qfile

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

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

func Test_MapFromJsonFileP_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	WriteFileTextP(fs, "test.json", `{"k": "v"}`)

	configMap := MapFromJsonFileP(fs, "test.json", false)

	a.Len(configMap, 1)
	a.Equal("v", configMap["k"])
}

// ==================== MapFromJsonFile 覆盖率测试 ====================

func TestMapFromJsonFile_fileNotFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, err := MapFromJsonFile(fs, "/nonexistent.json", false)
	a.Error(err)
}

func TestMapFromJsonFile_invalidJson(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/invalid.json", `{invalid}`)

	_, err := MapFromJsonFile(fs, "/invalid.json", false)
	a.Error(err)
}

func TestMapFromJsonFileP_panicsOnError(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	defer func() {
		r := recover()
		a.NotNil(r)
	}()

	MapFromJsonFileP(fs, "/nonexistent.json", false)
}

// ==================== FromJsonFile 覆盖率测试 ====================

func TestFromJsonFile_fileNotFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	var result map[string]any
	err := FromJsonFile(fs, "/nonexistent.json", false, &result)
	a.Error(err)
}

func TestFromJsonFileP_panicsOnError(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	defer func() {
		r := recover()
		a.NotNil(r)
	}()

	var result map[string]any
	FromJsonFileP(fs, "/nonexistent.json", false, &result)
}

func TestFromJsonFile_invalidJson(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/invalid.json", `{invalid}`)

	var result map[string]any
	err := FromJsonFile(fs, "/invalid.json", false, &result)
	a.Error(err)
}

func TestFromJsonFile_withEnvsubst(t *testing.T) {
	a := require.New(t)

	t.Setenv("FILE_VAR", "file_value")

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.json", `{"key": "${FILE_VAR}"}`)

	var result map[string]any
	err := FromJsonFile(fs, "/test.json", true, &result)
	a.NoError(err)
	a.Equal("file_value", result["key"])
}
