package test

import (
	"testing"

	"github.com/fastgh/go-comm"
	"github.com/stretchr/testify/require"
)

func Test_Vars2Pair_happy(t *testing.T) {
	a := require.New(t)

	a.Nil(comm.Vars2Pair(nil))
	a.Nil(comm.Vars2Pair(map[string]any{}))

	vars := map[string]any{
		"k1": "v1",
		"k2": "v2",
	}
	pairs := comm.Vars2Pair(vars)
	a.Len(pairs, 2)

	if pairs[0] == "k1=v1" {
		a.Equal("k2=v2", pairs[1])
	} else {
		a.Equal("k1=v1", pairs[1])
	}
}

func Test_Text2Vars_happy(t *testing.T) {
	a := require.New(t)

	m := comm.Text2Vars(`
L1
L2=
L3=V3

=
=V6

      L8=V8
L9= V9
L10=V10-1,V10-2
L11=V11    `)

	a.Len(m, 6)
	a.Equal(m["L2"], "")
	a.Equal(m["L3"], "V3")
	a.Equal(m["L8"], "V8")
	a.Equal(m["L9"], " V9")
	a.Equal(m["L10"], "V10-1,V10-2")
	a.Equal(m["L11"], "V11    ")
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
