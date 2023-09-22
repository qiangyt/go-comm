package test

import (
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func Test_Vars2Pair_happy(t *testing.T) {
	a := require.New(t)

	a.Nil(comm.Vars2Pair(nil))
	a.Nil(comm.Vars2Pair(map[string]string{}))

	vars := map[string]string{
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

func Test_ParseCommandOutput_json(t *testing.T) {
	a := require.New(t)

	text := `$json$

12`
	r := comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal(12, cast.ToInt(r.Json))
	a.Len(r.Vars, 0)

	text = `$json$

"12"`
	r = comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal("12", r.Json)

	text = `$json$

""`
	r = comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal("", r.Json)
	a.Len(r.Vars, 0)

	text = `$json$

["12"]`
	r = comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal([]any{"12"}, r.Json)
	a.Len(r.Vars, 0)

	text = `$json$

true`
	r = comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal(true, r.Json)
	a.Len(r.Vars, 0)

	text = `$json$

{"key": "value"}`
	r = comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_JSON, r.Kind)
	a.Equal(map[string]any{"key": "value"}, r.Json)
	a.Len(r.Vars, 0)

	text = `$json$

xyz`
	a.Panics(func() { comm.ParseCommandOutputP(text) }, "json: xyz")
}

func Test_ParseCommandOutput_vars(t *testing.T) {
	a := require.New(t)

	text := `$vars$

Key=Value`
	r := comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_VARS, r.Kind)
	a.Len(r.Vars, 1)
	a.Equal("Value", r.Vars["Key"])
	a.Nil(r.Json)
}

func Test_ParseCommandOutput_text(t *testing.T) {
	a := require.New(t)

	text := "something"
	r := comm.ParseCommandOutputP(text)
	a.Equal(text, r.Text)
	a.Equal(comm.COMMAND_OUTPUT_KIND_TEXT, r.Kind)
	a.Len(r.Vars, 0)
	a.Nil(r.Json)
}
