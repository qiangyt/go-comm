package comm

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func RunGoShellCommand(vars map[string]any, dir string, cmd string) CommandOutput {
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
		interp.ExecHandler(GoshExecHandler(6 * time.Second)),
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

	return ParseCommandOutput(out.String())
}
