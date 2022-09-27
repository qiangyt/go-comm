package test

import (
	"testing"

	"github.com/fastgh/go-comm"
	"github.com/stretchr/testify/require"
)

func Test_VarsToPair_happy(t *testing.T) {
	a := require.New(t)

	a.Nil(comm.VarsToPair(nil))
	a.Nil(comm.VarsToPair(map[string]any{}))

	vars := map[string]any{
		"k1": "v1",
		"k2": "v2",
	}
	pairs := comm.VarsToPair(vars)
	a.Len(pairs, 2)

	if pairs[0] == "k1=v1" {
		a.Equal("k2=v2", pairs[1])
	} else {
		a.Equal("k1=v1", pairs[1])
	}
}

func Test_RunGoShellCommand_happy(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"YOU": "fastgh",
	}

	output := comm.RunGoShellCommand(vars, "", "echo hello")
	a.Equal("hello\n", output)

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
	a.Equal("Hi fastgh\n", output)

	output = comm.RunShellCommand(vars, "", "gosh", "echo gosh")
	a.Equal("gosh\n", output)
}
