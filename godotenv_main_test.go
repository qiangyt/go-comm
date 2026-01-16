package comm

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestLoadEnv_noFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	err := LoadEnv(fs)
	a.Error(err) // .env doesn't exist
}

func TestLoadEnv_withFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	err := LoadEnv(fs)
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
	a.Equal("value2", os.Getenv("KEY2"))
}

func TestLoadEnv_multipleFiles(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env1", "KEY1=value1")
	WriteFileTextP(fs, ".env2", "KEY2=value2")

	err := LoadEnv(fs, ".env1", ".env2")
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
	a.Equal("value2", os.Getenv("KEY2"))
}

func TestOverloadEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1")

	err := OverloadEnv(fs)
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
}

func TestReadEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	result, err := ReadEnv(fs)
	a.NoError(err)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}

func TestReadEnv_noFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, err := ReadEnv(fs)
	a.Error(err)
}

func TestExecEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "TEST_VAR=test_value")

	// Use echo command which should work on all platforms
	err := ExecEnv(fs, []string{".env"}, "echo", []string{"hello"})
	// The command should execute successfully
	a.NoError(err)
}

func TestWriteEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	envMap := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	err := WriteEnv(fs, envMap, ".env")
	a.NoError(err)

	content := ReadFileTextP(fs, ".env")
	a.Contains(content, "KEY1")
	a.Contains(content, "value1")
}

func TestWriteEnv_numericValue(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	envMap := map[string]string{
		"PORT": "8080",
	}

	err := WriteEnv(fs, envMap, ".env")
	a.NoError(err)

	content := ReadFileTextP(fs, ".env")
	a.Contains(content, "PORT=8080")
}

func TestMarshalEnv_happy(t *testing.T) {
	a := require.New(t)

	envMap := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result, err := MarshalEnv(envMap)
	a.NoError(err)
	a.NotEmpty(result)
	// Should be sorted
	lines := strings.Split(result, "\n")
	a.True(lines[0] < lines[1])
}

func TestMarshalEnv_numeric(t *testing.T) {
	a := require.New(t)

	envMap := map[string]string{
		"PORT": "8080",
	}

	result, err := MarshalEnv(envMap)
	a.NoError(err)
	a.Contains(result, "PORT=8080")
}

func TestEnvFilenamesOrDefault_empty(t *testing.T) {
	result := envFilenamesOrDefault([]string{})
	a := require.New(t)
	a.Equal([]string{".env"}, result)
}

func TestEnvFilenamesOrDefault_provided(t *testing.T) {
	a := require.New(t)

	input := []string{".env1", ".env2"}
	result := envFilenamesOrDefault(input)
	a.Equal(input, result)
}

func TestReadEnvFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	result, err := readEnvFile(fs, ".env")
	a.NoError(err)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}

func TestReadEnvFile_notFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, err := readEnvFile(fs, ".env")
	a.Error(err)
}

func TestEnvDoubleQuoteEscape_happy(t *testing.T) {
	a := require.New(t)

	// Test escaping of special characters
	result := envDoubleQuoteEscape("test\nvalue")
	a.Equal(`test\nvalue`, result)
}
