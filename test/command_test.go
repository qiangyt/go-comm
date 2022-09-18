package test

import (
	"testing"

	"github.com/fastgh/go-comm"
	"github.com/stretchr/testify/require"
)

func Test_VarsToPair_happy(t *testing.T) {
	a := require.New(t)

	a.Nil(comm.VarsToPair(nil))
	a.Nil(comm.VarsToPair(map[string]string{}))

	vars := map[string]string{
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

	output := comm.RunGoShellCommand("", "echo hello")
	a.Equal("hello\n", output)

	a.Panics(func() {
		comm.RunGoShellCommand("", "fail.sh")
	})
}

func Test_RunGoShellCommand_gosh(t *testing.T) {
	a := require.New(t)

	output := comm.RunShellCommand("", "", "echo default")
	a.Equal("default\n", output)

	output = comm.RunShellCommand("", "gosh", "echo gosh")
	a.Equal("gosh\n", output)
}
