package test

import (
	"strings"
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_RenderAsTemplate_happy(t *testing.T) {
	a := require.New(t)

	data := map[string]any{
		"x1": "A",
		"x2": "B",
	}
	actual := comm.RenderAsTemplateP("begin {{.x1}} {{.x2}} end", data)

	a.Equal("begin A B end", actual)
}

func Test_RenderWithTemplate_happy(t *testing.T) {
	a := require.New(t)

	actual := strings.Builder{}
	data := map[string]any{
		"x1": "A",
		"x2": "B",
	}
	comm.RenderWithTemplateP(&actual, "test", "begin {{.x1}} {{.x2}} end", data)

	a.Equal("begin A B end", actual.String())
}

func Test_JoinedLinesAsBytes_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]byte{}, comm.JoinedLinesAsBytes())
	a.Equal([]byte("1"), comm.JoinedLinesAsBytes("1"))
	a.Equal([]byte("1\n2"), comm.JoinedLinesAsBytes("1", "2"))
}

func Test_SubstVars_noLocalVars(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"k0": "v0",
		"k":  "vParent",
	}

	actualMap := comm.SubstVarsP(true, map[string]any{
		"template": "prefix-{{.k0}}-{{.k}}-suffix",
	}, vars)
	a.Len(actualMap, 2)
	a.Equal("prefix-v0-vParent-suffix", actualMap["template"])

	actualVars := actualMap["vars"].(map[string]any)
	a.Len(actualVars, 2)
	a.Equal("v0", actualVars["k0"])
	a.Equal("vParent", actualVars["k"])
}

func Test_SubstVars_hasDifferentLocalVars(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"k0": "v0",
		"k":  "vParent",
	}

	actualMap := comm.SubstVarsP(true, map[string]any{
		"vars": map[string]any{
			"k1": "v1",
		},
		"template": "prefix-{{.k0}}-{{.k}}-{{.k1}}-suffix",
	}, vars)
	a.Len(actualMap, 2)
	a.Equal("prefix-v0-vParent-v1-suffix", actualMap["template"])

	actualVars := actualMap["vars"].(map[string]any)
	a.Len(actualVars, 3)
	a.Equal("v0", actualVars["k0"])
	a.Equal("vParent", actualVars["k"])
	a.Equal("v1", actualVars["k1"])
}

func Test_SubstVars_hasOverwrittenLocalVars(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"k0": "v0",
		"k":  "vParent",
	}
	// has overwritten local vars
	actualMap := comm.SubstVarsP(true, map[string]any{
		"vars": map[string]any{
			"k":  "vChild",
			"k1": "v1",
		},
		"template": "prefix-{{.k0}}-{{.k}}-{{.k1}}-suffix",
	}, vars)
	a.Len(actualMap, 2)
	a.Equal("prefix-v0-vChild-v1-suffix", actualMap["template"])

	actualVars := actualMap["vars"].(map[string]any)
	a.Len(actualVars, 3)
	a.Equal("v0", actualVars["k0"])
	a.Equal("vChild", actualVars["k"])
	a.Equal("v1", actualVars["k1"])
}

func Test_SubstVars_skip(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"k0": "v0",
		"k":  "vParent",
	}
	// has overwritten local vars
	actualMap := comm.SubstVarsP(true, map[string]any{
		"vars": map[string]any{
			"k":  "vChild",
			"k1": "v1",
		},
		"template":       "prefix-{{.k0}}-{{.k}}-{{.k1}}-suffix",
		"templateToSkip": "prefix-{{.k0}}-{{.k}}-{{.k1}}-suffix",
	}, vars, "templateToSkip")
	a.Len(actualMap, 3)
	a.Equal("prefix-v0-vChild-v1-suffix", actualMap["template"])
	a.Equal("prefix-{{.k0}}-{{.k}}-{{.k1}}-suffix", actualMap["templateToSkip"])

	actualVars := actualMap["vars"].(map[string]any)
	a.Len(actualVars, 3)
	a.Equal("v0", actualVars["k0"])
	a.Equal("vChild", actualVars["k"])
	a.Equal("v1", actualVars["k1"])
}

func Test_TextLine2Array_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]string{}, comm.TextLine2Array(" \n \r \t  "))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1,2"))
	a.Equal([]string{"1", "2", "3", "4"}, comm.TextLine2Array(" 1, 2 \t,\r3\n ,4\n"))

	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1\t2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1\n2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1\r2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1\r\n2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1\n\r2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1;2"))
	a.Equal([]string{"1", "2"}, comm.TextLine2Array("1|2"))
}
