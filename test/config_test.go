package test

import (
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_DecodeWithYaml_happy(t *testing.T) {
	a := require.New(t)

	devault := map[string]any{
		"a": 100,
		"b": "BBB",
		"c": true,
	}

	yamlText := `
a: 1
b: B
`
	type Temp struct {
		A int    `mapstruct:"a"`
		B string `mapstruct:"b"`
		C bool   `mapstruct:"c"`
	}
	r, _ := comm.DecodeWithYamlP(yamlText, comm.StrictConfigConfig(), &Temp{}, devault)
	a.Equal(1, r.A)
	a.Equal("B", r.B)
	a.True(r.C)
}
