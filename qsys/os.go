package qsys

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	plog "github.com/phuslu/log"
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qerr"
)

func DefaultOutput() io.Writer {
	if IsTerminal() {
		return os.Stdout
	} else {
		return io.Discard
	}
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsTerminal() bool {
	return plog.IsTerminal(os.Stdout.Fd())
}

func ExecutableP() string {
	r, err := Executable()
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func Executable() (string, error) {
	r, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, "get the path name of the executable file")
	}
	r, err = filepath.EvalSymlinks(r)
	if err != nil {
		return "", errors.Wrapf(err, "evaluate the symbol linke of the executable file: %s", r)
	}
	return r, nil
}

func WorkingDirectoryP() string {
	r, err := WorkingDirectory()
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func WorkingDirectory() (string, error) {
	r, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get working directory")
	}
	return r, nil
}

func AbsPath(cwd string, _path string) string {
	r := filepath.Clean(_path)
	if filepath.IsAbs(r) {
		return r
	}

	return filepath.Join(filepath.Clean(cwd), _path)
}
