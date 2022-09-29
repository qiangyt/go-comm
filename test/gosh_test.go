package test

import (
	"testing"

	"github.com/fastgh/go-comm"
	"github.com/stretchr/testify/require"
)

func Test_RunGoShellCommand_happy(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"YOU": "fastgh",
	}

	output := comm.RunGoShellCommand(vars, "", "echo '$vars$\n\nkey=value'\n")
	a.Equal(comm.COMMAND_OUTPUT_KIND_VARS, output.Kind)
	a.Equal("value", output.Vars["key"])

	a.Panics(func() {
		comm.RunGoShellCommand(vars, "", "fail.sh")
	})
}

func Test_RunGoShellCommand_gosh(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"YOU": "fastgh",
	}

	output := comm.RunShellCommand(vars, "", "", "echo Hi ${YOU}")
	a.Equal(comm.COMMAND_OUTPUT_KIND_TEXT, output.Kind)
	a.Equal("Hi fastgh\n", output.Text)

	output = comm.RunShellCommand(vars, "", "gosh", "echo '$json$\n\ntrue'")
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, output.Kind)
	a.Equal(true, output.Json)
}
