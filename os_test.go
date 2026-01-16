package comm

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsWindows(t *testing.T) {
	a := require.New(t)

	result := IsWindows()
	if runtime.GOOS == "windows" {
		a.True(result)
	} else {
		a.False(result)
	}
}

func TestIsDarwin(t *testing.T) {
	a := require.New(t)

	result := IsDarwin()
	if runtime.GOOS == "darwin" {
		a.True(result)
	} else {
		a.False(result)
	}
}

func TestIsLinux(t *testing.T) {
	a := require.New(t)

	result := IsLinux()
	if runtime.GOOS == "linux" {
		a.True(result)
	} else {
		a.False(result)
	}
}

func TestEnvironMap_happy(t *testing.T) {
	a := require.New(t)

	result, err := EnvironMap(nil)
	a.NoError(err)
	a.NotNil(result)
	a.True(len(result) > 0)
}

func TestEnvironMapP_happy(t *testing.T) {
	a := require.New(t)

	result := EnvironMapP(nil)
	a.NotNil(result)
}

func TestEnvironMap_withOverrides(t *testing.T) {
	a := require.New(t)

	overrides := map[string]string{
		"MY_TEST_VAR": "test_value",
	}

	result, err := EnvironMap(overrides)
	a.NoError(err)
	a.Equal("test_value", result["MY_TEST_VAR"])
}

func TestEnvironList_happy(t *testing.T) {
	a := require.New(t)

	result, err := EnvironList(nil)
	a.NoError(err)
	a.NotNil(result)
	a.True(len(result) > 0)
}

func TestEnvironListP_happy(t *testing.T) {
	a := require.New(t)

	result := EnvironListP(nil)
	a.NotNil(result)
}

func TestEnvSubst_happy(t *testing.T) {
	a := require.New(t)

	env := map[string]string{
		"MY_VAR": "hello",
	}

	result, err := EnvSubst("value is $MY_VAR", env)
	a.NoError(err)
	a.Equal("value is hello", result)
}

func TestEnvSubstP_happy(t *testing.T) {
	a := require.New(t)

	env := map[string]string{
		"NAME": "world",
	}

	result := EnvSubstP("hello $NAME", env)
	a.Equal("hello world", result)
}

func TestEnvSubstSlice_happy(t *testing.T) {
	a := require.New(t)

	env := map[string]string{
		"A": "1",
		"B": "2",
	}

	inputs := []string{"$A", "$B", "$A+$B"}
	result, err := EnvSubstSlice(inputs, env)
	a.NoError(err)
	a.Equal([]string{"1", "2", "1+2"}, result)
}

func TestEnvSubstSliceP_happy(t *testing.T) {
	a := require.New(t)

	env := map[string]string{
		"X": "test",
	}

	inputs := []string{"$X"}
	result := EnvSubstSliceP(inputs, env)
	a.Equal([]string{"test"}, result)
}

func TestIsTerminal(t *testing.T) {
	// Just test that it doesn't panic
	_ = IsTerminal()
}

func TestExecutable_happy(t *testing.T) {
	a := require.New(t)

	result, err := Executable()
	a.NoError(err)
	a.NotEmpty(result)
}

func TestExecutableP_happy(t *testing.T) {
	a := require.New(t)

	result := ExecutableP()
	a.NotEmpty(result)
}

func TestWorkingDirectory_happy(t *testing.T) {
	a := require.New(t)

	result, err := WorkingDirectory()
	a.NoError(err)
	a.NotEmpty(result)
}

func TestWorkingDirectoryP_happy(t *testing.T) {
	a := require.New(t)

	result := WorkingDirectoryP()
	a.NotEmpty(result)
}

func TestAbsPath_relative(t *testing.T) {
	a := require.New(t)

	// Use platform-specific paths
	base := "/home/user"
	if IsWindows() {
		base = `C:\home\user`
	}
	result := AbsPath(base, "mydir/file.txt")
	a.Contains(result, "mydir")
	a.Contains(result, "file.txt")
}

func TestAbsPath_absolute(t *testing.T) {
	a := require.New(t)

	base := "/home/user"
	absolutePath := "/etc/config.conf"
	if IsWindows() {
		base = `C:\home\user`
		absolutePath = `D:\etc\config.conf`
	}
	result := AbsPath(base, absolutePath)
	a.Equal(absolutePath, result)
}

func TestAbsPath_current(t *testing.T) {
	a := require.New(t)

	base := "/home/user"
	if IsWindows() {
		base = `C:\home\user`
	}
	result := AbsPath(base, "./file.txt")
	a.Contains(result, "file.txt")
}
