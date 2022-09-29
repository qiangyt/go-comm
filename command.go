package comm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"mvdan.cc/sh/v3/interp"
)

type CommandOutputKind byte

const (
	COMMAND_OUTPUT_KIND_TEXT CommandOutputKind = iota
	COMMAND_OUTPUT_KIND_VARS
	COMMAND_OUTPUT_KIND_JSON
)

type CommandOutputT struct {
	Kind CommandOutputKind
	Vars map[string]any
	Text string
	Json any
}

type CommandOutput = *CommandOutputT

func ParseCommandOutput(outputText string) CommandOutput {
	r := &CommandOutputT{Kind: COMMAND_OUTPUT_KIND_TEXT, Text: outputText}

	if strings.HasPrefix(outputText, "$json$\n\n") {
		jsonBody := outputText[len("$json$\n\n"):]

		err := json.Unmarshal([]byte(jsonBody), &r.Json)
		if err != nil {
			panic(errors.New("invalid json: " + jsonBody))
		}

		r.Kind = COMMAND_OUTPUT_KIND_JSON
		return r
	}

	if strings.HasPrefix(outputText, "$vars$\n\n") {
		varsBody := outputText[len("$vars$\n\n"):]

		r.Vars = Text2Vars(varsBody)
		r.Kind = COMMAND_OUTPUT_KIND_VARS
		return r
	}

	return r
}

var errUnsupportedOS = errors.New("unsupported OS")

func Vars2Pair(vars map[string]any) []string {
	if len(vars) == 0 {
		return nil
	}

	r := make([]string, 0, len(vars))
	for k, v := range vars {
		r = append(r, k+"="+cast.ToString(v))
	}
	return r
}

func Text2Vars(text string) map[string]any {
	pairs := Text2Lines(text)
	return Pair2Vars(pairs)
}

func Pair2Vars(pairs []string) map[string]any {
	if len(pairs) == 0 {
		return map[string]any{}
	}

	r := map[string]any{}
	for _, pair := range pairs {
		pair = strings.TrimLeft(pair, "\t \r")
		pos := strings.IndexByte(pair, '=')
		if pos <= 0 {
			continue
		}
		k := pair[:pos]
		if pos == len(pair)-1 {
			r[k] = ""
		} else {
			r[k] = pair[pos+1:]
		}
	}
	return r
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return devNull{}, nil
	}
	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}

func RunShellCommand(vars map[string]any, dir string, sh string, cmd string) CommandOutput {
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

func RunAdminCommand(vars map[string]any, adminPassword string, dir string, cmd string) CommandOutput {
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

func RunUserCommand(vars map[string]any, dir string, cmd string) CommandOutput {
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
func RunAppleScript(vars map[string]any, adminPassword string, dir string, script string) CommandOutput {
	subArgs := []string{fmt.Sprintf(`do shell script "%s"`, script)}

	if len(adminPassword) > 0 {
		subArgs = append(subArgs, fmt.Sprintf(`password "%s"`, adminPassword))
	}
	subArgs = append(subArgs, "with administrator privileges")

	return RunCommandWithoutInput(vars, dir, "osascript", "-e", strings.Join(subArgs, " "))
}

func RunSudoCommand(vars map[string]any, sudoerPassword string, dir string, command string) CommandOutput {
	if len(sudoerPassword) > 0 {
		return RunCommandWithInput(vars, dir, "sudo", "sh", command)(sudoerPassword)
	}

	return RunCommandWithoutInput(vars, dir, "sudo", "sh", command)
}

func newExecCommand(vars map[string]any, dir string, cmd string, args ...string) *exec.Cmd {
	r := exec.Command(cmd, args...)
	r.Env = EnvironList(vars)
	r.Dir = dir
	return r
}

func RunCommandWithoutInput(vars map[string]any, dir string, cmd string, args ...string) CommandOutput {
	_cmd := newExecCommand(vars, dir, cmd, args...)
	b, err := _cmd.Output()
	if err != nil {
		cli := strings.Join(append([]string{cmd}, args...), " ")
		panic(errors.Wrapf(err, "failed to get output for command '%s'", cli))
	}

	return ParseCommandOutput(cast.ToString(b))
}

func RunCommandWithInput(vars map[string]any, dir string, cmd string, args ...string) func(...string) CommandOutput {
	return func(input ...string) CommandOutput {
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

		return ParseCommandOutput(cast.ToString(b))
	}
}
