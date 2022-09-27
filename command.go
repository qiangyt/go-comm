package comm

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var errUnsupportedOS = errors.New("unsupported OS")

func VarsToPair[T any](vars map[string]T) []string {
	if len(vars) == 0 {
		return nil
	}

	r := make([]string, 0, len(vars))
	for k, v := range vars {
		r = append(r, k+"="+cast.ToString(v))
	}
	return r
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devNull{}, nil
	}
	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}

func RunGoShellCommand[T any](vars map[string]T, dir string, cmd string) string {
	var err error

	var sf *syntax.File
	if sf, err = syntax.NewParser().Parse(strings.NewReader(cmd), ""); err != nil {
		panic(errors.Wrapf(err, "failed to parse command: \n%s", cmd))
	}

	out := strings.Builder{}

	environ := append(os.Environ(), EnvironList(vars)...)

	opts := []interp.RunnerOption{
		interp.Params("-e"),
		interp.Env(expand.ListEnviron(environ...)),
		interp.ExecHandler(interp.DefaultExecHandler(6 * time.Second)),
		interp.OpenHandler(openHandler),
		interp.StdIO(nil, &out, &out),
	}
	if len(dir) > 0 {
		opts = append(opts, interp.Dir(dir))
	}

	var runner *interp.Runner
	if runner, err = interp.New(opts...); err != nil {
		panic(errors.Wrapf(err, "failed to create runner for command: \n%s", cmd))
	}

	if err = runner.Run(context.TODO(), sf); err != nil {
		panic(errors.Wrapf(err, "failed to run command: \n%s", cmd))
	}

	return out.String()
}

func RunShellCommand[T any](vars map[string]T, dir string, sh string, cmd string) string {
	if len(sh) == 0 || sh == "gosh" {
		return RunGoShellCommand(vars, dir, cmd)
	}

	switch DefaultOSType() {
	case Windows:
		panic(errUnsupportedOS)
	case Linux, Darwin:
		return RunCommandWithoutInput(vars, dir, sh, cmd)
	default:
		panic(errUnsupportedOS)
	}
}

/*
func RunShellScriptFile(afs afero.Fs, url string, credentials ufs.Credentials, timeout time.Duration,
	dir string, sh string) string {

	scriptContent := ufs.DownloadText(afs, url, credentials, timeout)
	return RunShellCommand(dir, sh, scriptContent)
}*/

func RunAdminCommand[T any](vars map[string]T, adminPassword string, dir string, cmd string) string {
	switch DefaultOSType() {
	case Windows:
		panic(errUnsupportedOS)
	case Linux:
		return RunSudoCommand(vars, adminPassword, dir, cmd)
	case Darwin:
		return RunAppleScript(vars, adminPassword, dir, cmd)
	default:
		panic(errUnsupportedOS)
	}
}

func RunUserCommand[T any](vars map[string]T, dir string, cmd string) string {
	switch DefaultOSType() {
	case Windows:
		panic(errUnsupportedOS)
	case Linux:
		return RunCommandWithoutInput(vars, dir, "sh", cmd)
	case Darwin:
		return RunCommandWithoutInput(vars, dir, "open", cmd)
	default:
		panic(errUnsupportedOS)
	}
}

// RunApplacript 运行 applacript
func RunAppleScript[T any](vars map[string]T, adminPassword string, dir string, script string) string {
	subArgs := []string{fmt.Sprintf(`do shell script "%s"`, script)}

	if len(adminPassword) > 0 {
		subArgs = append(subArgs, fmt.Sprintf(`password "%s"`, adminPassword))
	}
	subArgs = append(subArgs, "with administrator privileges")

	return RunCommandWithoutInput(vars, dir, "osascript", "-e", strings.Join(subArgs, " "))
}

func RunSudoCommand[T any](vars map[string]T, sudoerPassword string, dir string, command string) string {
	if len(sudoerPassword) > 0 {
		return RunCommandWithInput(vars, dir, "sudo", "sh", command)(sudoerPassword)
	}

	return RunCommandWithoutInput(vars, dir, "sudo", "sh", command)
}

func newExecCommand[T any](vars map[string]T, dir string, cmd string, args ...string) *exec.Cmd {
	r := exec.Command(cmd, args...)
	r.Env = EnvironList(vars)
	r.Dir = dir
	return r
}

func RunCommandWithoutInput[T any](vars map[string]T, dir string, cmd string, args ...string) string {
	_cmd := newExecCommand(vars, dir, cmd, args...)
	b, err := _cmd.Output()
	if err != nil {
		cli := strings.Join(append([]string{cmd}, args...), " ")
		panic(errors.Wrapf(err, "failed to get output for command '%s'", cli))
	}

	return strings.TrimSpace(string(b))
}

func RunCommandWithInput[T any](vars map[string]T, dir string, cmd string, args ...string) func(...string) string {
	return func(input ...string) string {
		cli := cmd + " " + strings.Join(args, " ")

		_cmd := newExecCommand(vars, dir, cmd, args...)

		stdin, err := _cmd.StdinPipe()
		if err != nil {
			panic(errors.Wrapf(err, "failed to open stdin for command '%s'", cli))
		}
		defer func() {
			if stdin != nil {
				stdin.Close()
				stdin = nil
			}
		}()

		io.WriteString(stdin, strings.Join(input, " "))
		stdin.Close()
		stdin = nil

		b, err := _cmd.Output()
		if err != nil {
			panic(errors.Wrapf(err, "failed to get output for command '%s'", cli))
		}

		return strings.TrimSpace(string(b))
	}
}
