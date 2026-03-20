package comm

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOutput(t *testing.T) {
	result := DefaultOutput()
	a := require.New(t)
	a.NotNil(result)
}

func Test_EnvironMap_happy(t *testing.T) {
	a := require.New(t)

	actual := EnvironMapP(nil)
	a.True(len(actual) > 0)
	a.True(len(actual["PATH"]) > 0)

	overrides := map[string]string{
		"k1":   "v1",
		"k2":   "v2",
		"PATH": "overrided_path",
	}

	actual = EnvironMapP(overrides)
	a.Equal("v1", actual["k1"])
	a.Equal("v2", actual["k2"])
	a.Equal("overrided_path", actual["PATH"])
}

func Test_AbsPath_happy(t *testing.T) {
	a := require.New(t)

	if IsWindows() {
		a.Equal("\\home", AbsPath("/home", "."))
		a.Equal("\\home", AbsPath("/home", "./"))
		a.Equal("\\", AbsPath("/home", ".."))
		a.Equal("\\", AbsPath("/home", "../"))

		a.Equal("\\home\\1\\2", AbsPath("/home/1", "2"))
		a.Equal("\\home\\1\\2\\3", AbsPath("/home/1", "2/3"))
		a.Equal("\\home\\1\\2\\3", AbsPath("/home/1", "./2/3"))
		a.Equal("\\home\\2\\3", AbsPath("/home/1", "../2/3"))

		a.Equal("\\home\\1\\2\\3", AbsPath("/home/1", "2/./3"))
		a.Equal("\\home\\1\\3", AbsPath("/home/1", "2/../3"))
	} else {
		a.Equal("/home", AbsPath("/home", "."))
		a.Equal("/home", AbsPath("/home", "./"))
		a.Equal("/", AbsPath("/home", ".."))
		a.Equal("/", AbsPath("/home", "../"))

		a.Equal("/home/1/2", AbsPath("/home/1", "2"))
		a.Equal("/home/1/2/3", AbsPath("/home/1", "2/3"))
		a.Equal("/home/1/2/3", AbsPath("/home/1", "./2/3"))
		a.Equal("/home/2/3", AbsPath("/home/1", "../2/3"))

		a.Equal("/home/1/2/3", AbsPath("/home/1", "2/./3"))
		a.Equal("/home/1/3", AbsPath("/home/1", "2/../3"))
	}
}

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
