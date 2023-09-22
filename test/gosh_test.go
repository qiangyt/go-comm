package test

import (
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_RunGoshCommand_happy(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{
		"YOU": "fastgh",
	}

	output := comm.RunGoshCommandP(vars, "", "echo '$vars$\n\nkey=value'\n", nil)
	a.Equal(comm.COMMAND_OUTPUT_KIND_VARS, output.Kind)
	a.Equal("value", output.Vars["key"])

	a.Panics(func() {
		comm.RunGoshCommandP(vars, "", "fail.sh", nil)
	})
}

func Test_RunShellCommand_gosh(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{
		"YOU": "fastgh",
	}

	output := comm.RunShellCommandP(vars, "", "", "echo Hi ${YOU}", nil)
	a.Equal(comm.COMMAND_OUTPUT_KIND_TEXT, output.Kind)
	a.Equal("Hi fastgh\n", output.Text)

	output = comm.RunShellCommandP(vars, "", "gosh", "echo '$json$\n\ntrue'", nil)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, output.Kind)
	a.Equal(true, output.Json)
}
