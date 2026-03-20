package qfile

import (
	"os"
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestSysEnvFileNames_emptyShell(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "")
	a.NotNil(result)
}

func TestSysEnvFileNames_withZsh(t *testing.T) {
	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "zsh")
	a := require.New(t)
	a.NotNil(result)
}

func TestSysEnvFileNames_withBash(t *testing.T) {
	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "bash")
	a := require.New(t)
	a.NotNil(result)
}

func TestLoadEnvScripts_noFilenames(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}

	result, err := LoadEnvScripts(fs, vars)
	a.NoError(err)
	a.NotNil(result)
}

func TestLoadEnvScripts_withFilenames(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}
	// Use non-existent files - RunGoshCommand will error, but LoadEnvScripts should still return with error
	result, err := LoadEnvScripts(fs, vars, "/nonexistent/file1", "/nonexistent/file2")
	a.Error(err) // Expected to have errors
	a.NotNil(result)
}

func TestLoadEnvScript_etcPaths(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/etc")
	WriteFileLinesP(fs, "/etc/paths", "/usr/bin", "/usr/local/bin")

	vars := map[string]string{"PATH": "/existing/path"}
	result, err := LoadEnvScript(fs, vars, "/etc/paths")
	a.NoError(err)
	a.NotNil(result)
	a.Contains(result["PATH"], "/existing/path")
	a.Contains(result["PATH"], "/usr/bin")
}

func TestLoadEnvScript_nonExistent(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}

	result, err := LoadEnvScript(fs, vars, "/nonexistent/file")
	a.Error(err)
	a.NotNil(result)
}

func TestLoadEnv_noFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	err := comm.LoadEnv(fs)
	a.Error(err) // .env doesn't exist
}

func TestLoadEnv_withFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	err := comm.LoadEnv(fs)
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
	a.Equal("value2", os.Getenv("KEY2"))
}

func TestLoadEnv_multipleFiles(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env1", "KEY1=value1")
	WriteFileTextP(fs, ".env2", "KEY2=value2")

	err := comm.LoadEnv(fs, ".env1", ".env2")
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
	a.Equal("value2", os.Getenv("KEY2"))
}

func TestExecEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "TEST_VAR=test_value")

	// Use echo command which should work on all platforms
	err := comm.ExecEnv(fs, []string{".env"}, "echo", []string{"hello"})
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

	err := comm.WriteEnv(fs, envMap, ".env")
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

	err := comm.WriteEnv(fs, envMap, ".env")
	a.NoError(err)

	content := ReadFileTextP(fs, ".env")
	a.Contains(content, "PORT=8080")
}

func TestReadEnvFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	result, err := comm.ReadEnvFile(fs, ".env")
	a.NoError(err)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}

func TestReadEnvFile_notFound(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, err := comm.ReadEnvFile(fs, ".env")
	a.Error(err)
}

func TestOverloadEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1")

	err := comm.OverloadEnv(fs)
	a.NoError(err)
	a.Equal("value1", os.Getenv("KEY1"))
}

func TestReadEnv_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, ".env", "KEY1=value1\nKEY2=value2")

	result, err := comm.ReadEnv(fs)
	a.NoError(err)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}
