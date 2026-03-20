package qshell

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RunGoshCommand_happy(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{
		"YOU": "qiangyt",
	}

	output := RunGoshCommandP(vars, "", "echo '$vars$\n\nkey=value'\n", nil)
	a.Equal(COMMAND_OUTPUT_KIND_VARS, output.Kind)
	a.Equal("value", output.Vars["key"])

	a.Panics(func() {
		RunGoshCommandP(vars, "", "fail.sh", nil)
	})
}

func Test_RunShellCommand_gosh(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{
		"YOU": "qiangyt",
	}

	output := RunShellCommandP(vars, "", "", "echo Hi ${YOU}", nil)
	a.Equal(COMMAND_OUTPUT_KIND_TEXT, output.Kind)
	a.Equal("Hi qiangyt\n", output.Text)

	output = RunShellCommandP(vars, "", "gosh", "echo '$json$\n\ntrue'", nil)
	a.Equal(COMMAND_OUTPUT_KIND_JSON, output.Kind)
	a.Equal(true, output.Json)
}
