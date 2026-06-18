package qlang

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RenderAsTemplate_happy(t *testing.T) {
	a := require.New(t)

	data := map[string]any{
		"x1": "A",
		"x2": "B",
	}
	actual := RenderAsTemplateP("begin {{.x1}} {{.x2}} end", data)

	a.Equal("begin A B end", actual)
}

func Test_RenderWithTemplate_happy(t *testing.T) {
	a := require.New(t)

	actual := strings.Builder{}
	data := map[string]any{
		"x1": "A",
		"x2": "B",
	}
	RenderWithTemplateP(&actual, "test", "begin {{.x1}} {{.x2}} end", data)

	a.Equal("begin A B end", actual.String())
}

func Test_JoinedLinesAsBytes_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]byte{}, JoinedLinesAsBytes())
	a.Equal([]byte("1"), JoinedLinesAsBytes("1"))
	a.Equal([]byte("1\n2"), JoinedLinesAsBytes("1", "2"))
}

func Test_SubstVars_noLocalVars(t *testing.T) {
	a := require.New(t)

	vars := map[string]any{
		"k0": "v0",
		"k":  "vParent",
	}

	actualMap := SubstVarsP(true, map[string]any{
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

	actualMap := SubstVarsP(true, map[string]any{
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
	actualMap := SubstVarsP(true, map[string]any{
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
	actualMap := SubstVarsP(true, map[string]any{
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

	a.Equal([]string{}, TextLine2Array(" \n \r \t  "))
	a.Equal([]string{"1", "2"}, TextLine2Array("1,2"))
	a.Equal([]string{"1", "2", "3", "4"}, TextLine2Array(" 1, 2 \t,\r3\n ,4\n"))

	a.Equal([]string{"1", "2"}, TextLine2Array("1\t2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1\n2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1\r2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1\r\n2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1\n\r2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1;2"))
	a.Equal([]string{"1", "2"}, TextLine2Array("1|2"))
}

func TestRenderWithTemplate_happy(t *testing.T) {
	a := require.New(t)

	var buf bytes.Buffer
	err := RenderWithTemplate(&buf, "test", "Hello {{.Name}}!", map[string]any{"Name": "World"})
	a.NoError(err)
	a.Equal("Hello World!", buf.String())
}

func TestRenderWithTemplateP_happy(t *testing.T) {
	a := require.New(t)

	var buf bytes.Buffer
	RenderWithTemplateP(&buf, "test", "Value: {{.Value}}", map[string]any{"Value": 42})
	a.Equal("Value: 42", buf.String())
}

func TestRenderWithTemplate_error(t *testing.T) {
	a := require.New(t)

	var buf bytes.Buffer
	err := RenderWithTemplate(&buf, "test", "{{.Invalid", map[string]any{})
	a.Error(err)
}

func TestRenderAsTemplate_happy(t *testing.T) {
	a := require.New(t)

	result, err := RenderAsTemplate("Hello {{.Name}}!", map[string]any{"Name": "Test"})
	a.NoError(err)
	a.Equal("Hello Test!", result)
}

func TestRenderAsTemplateP_happy(t *testing.T) {
	a := require.New(t)

	result := RenderAsTemplateP("Count: {{.Count}}", map[string]any{"Count": 100})
	a.Equal("Count: 100", result)
}

func TestRenderAsTemplateArray_happy(t *testing.T) {
	a := require.New(t)

	tmplArray := []string{"Hello {{.Name}}", "Goodbye {{.Name}}"}
	result, err := RenderAsTemplateArray(tmplArray, map[string]any{"Name": "User"})
	a.NoError(err)
	a.Equal([]string{"Hello User", "Goodbye User"}, result)
}

func TestRenderAsTemplateArrayP_happy(t *testing.T) {
	a := require.New(t)

	tmplArray := []string{"A: {{.A}}", "B: {{.B}}"}
	result := RenderAsTemplateArrayP(tmplArray, map[string]any{"A": 1, "B": 2})
	a.Equal([]string{"A: 1", "B: 2"}, result)
}

func TestJoinedLines_happy(t *testing.T) {
	a := require.New(t)

	result := JoinedLines("line1", "line2", "line3")
	a.Equal("line1\nline2\nline3", result)
}

func TestJoinedLinesAsBytes_happy(t *testing.T) {
	a := require.New(t)

	result := JoinedLinesAsBytes("a", "b")
	a.Equal([]byte("a\nb"), result)
}

func TestToYaml_happy(t *testing.T) {
	a := require.New(t)

	data := map[string]any{"key": "value"}
	result, err := ToYaml("test data", data)
	a.NoError(err)
	a.Contains(result, "key: value")
}

func TestToYamlP_happy(t *testing.T) {
	a := require.New(t)

	data := map[string]any{"num": 42}
	result := ToYamlP("", data)
	a.Contains(result, "num: 42")
}

func TestToYaml_withHint(t *testing.T) {
	a := require.New(t)

	data := map[string]any{"name": "test"}
	result, err := ToYaml("config", data)
	a.NoError(err)
	a.Contains(result, "name: test")
}

func TestToYaml_noHint(t *testing.T) {
	a := require.New(t)

	data := map[string]int{"count": 5}
	result, err := ToYaml("", data)
	a.NoError(err)
	a.Contains(result, "count: 5")
}

func TestSubstVars_happy(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"vars": map[string]any{
			"name": "test",
		},
		"key": "$name",
	}

	result, err := SubstVars(false, m, nil)
	a.NoError(err)
	a.NotNil(result)
}

func TestSubstVarsP_happy(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key": "value",
	}

	result := SubstVarsP(false, m, nil)
	a.NotNil(result)
}

func TestSubstVars_withGoTemplate(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"vars": map[string]any{
			"name": "test",
		},
		"key": "{{.name}}",
	}

	result, err := SubstVars(true, m, nil)
	a.NoError(err)
	a.NotNil(result)
}

func TestSubstVars_withParentVars(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key": "$parent_var",
	}
	parentVars := map[string]any{
		"parent_var": "parent_value",
	}

	result, err := SubstVars(false, m, parentVars)
	a.NoError(err)
	a.NotNil(result)
}

func TestSubstVars_withKeysToSkip(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"keep":    "original",
		"process": "value",
	}

	result, err := SubstVars(false, m, nil, "keep")
	a.NoError(err)
	a.Equal("original", result["keep"])
}

func TestTextLine2Array_comma(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a,b,c")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_tab(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a\tb\tc")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_newline(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a\nb\nc")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_carriage(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a\rb\rc")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_semicolon(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a;b;c")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_pipe(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a|b|c")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_space(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("a b c")
	a.Equal([]string{"a", "b", "c"}, result)
}

func TestTextLine2Array_empty(t *testing.T) {
	a := require.New(t)

	result := TextLine2Array("")
	a.Equal([]string{}, result)

	result = TextLine2Array("   ")
	a.Equal([]string{}, result)
}

func TestText2Lines_happy(t *testing.T) {
	a := require.New(t)

	result := Text2Lines("line1\nline2\nline3")
	a.Equal([]string{"line1", "line2", "line3"}, result)
}

func TestSliceLines(t *testing.T) {
	t.Run("无 line 和 limit 参数返回原文", func(t *testing.T) {
		a := require.New(t)
		result := SliceLines("line1\nline2\nline3", nil, nil)
		a.Equal("line1\nline2\nline3", result)
	})

	t.Run("从指定行开始", func(t *testing.T) {
		a := require.New(t)
		line := 2
		result := SliceLines("line1\nline2\nline3", &line, nil)
		a.Equal("line2\nline3", result)
	})

	t.Run("限制行数", func(t *testing.T) {
		a := require.New(t)
		limit := 2
		result := SliceLines("line1\nline2\nline3", nil, &limit)
		a.Equal("line1\nline2", result)
	})

	t.Run("从指定行开始并限制行数", func(t *testing.T) {
		a := require.New(t)
		line := 2
		limit := 1
		result := SliceLines("line1\nline2\nline3", &line, &limit)
		a.Equal("line2", result)
	})

	t.Run("limit 超过剩余行数", func(t *testing.T) {
		a := require.New(t)
		line := 2
		limit := 10
		result := SliceLines("line1\nline2\nline3", &line, &limit)
		a.Equal("line2\nline3", result)
	})

	t.Run("line 超过总行数返回空", func(t *testing.T) {
		a := require.New(t)
		line := 10
		result := SliceLines("line1\nline2\nline3", &line, nil)
		a.Equal("", result)
	})

	t.Run("负数 line 从第一行开始", func(t *testing.T) {
		a := require.New(t)
		line := 0
		limit := 2
		result := SliceLines("line1\nline2\nline3", &line, &limit)
		a.Equal("line1\nline2", result)
	})

	t.Run("空文本返回空", func(t *testing.T) {
		a := require.New(t)
		line := 1
		limit := 2
		result := SliceLines("", &line, &limit)
		a.Equal("", result)
	})
}

func TestSliceLinesP(t *testing.T) {
	t.Run("基本切分", func(t *testing.T) {
		a := require.New(t)
		result := SliceLinesP("line1\nline2\nline3", 1, 2)
		a.Equal("line1\nline2", result)
	})

	t.Run("从中间开始", func(t *testing.T) {
		a := require.New(t)
		result := SliceLinesP("line1\nline2\nline3", 2, 2)
		a.Equal("line2\nline3", result)
	})

	t.Run("负数 line 表示从头开始", func(t *testing.T) {
		a := require.New(t)
		result := SliceLinesP("line1\nline2\nline3", 0, 2)
		a.Equal("line1\nline2", result)
	})

	t.Run("负数 limit 表示不限制", func(t *testing.T) {
		a := require.New(t)
		result := SliceLinesP("line1\nline2\nline3", 2, 0)
		a.Equal("line2\nline3", result)
	})

	t.Run("都为负数返回全文", func(t *testing.T) {
		a := require.New(t)
		result := SliceLinesP("line1\nline2\nline3", -1, -1)
		a.Equal("line1\nline2\nline3", result)
	})
}

func TestShorten(t *testing.T) {
	t.Run("字符串不超过长度返回原文", func(t *testing.T) {
		a := require.New(t)
		result := Shorten("hello", 10)
		a.Equal("hello", result)
	})

	t.Run("字符串超过长度时截断", func(t *testing.T) {
		a := require.New(t)
		result := Shorten("hello world", 8)
		a.Equal("hello...", result)
	})

	t.Run("空字符串", func(t *testing.T) {
		a := require.New(t)
		result := Shorten("", 10)
		a.Equal("", result)
	})

	t.Run("maxLen 小于等于 3", func(t *testing.T) {
		a := require.New(t)
		result := Shorten("hello", 3)
		a.Equal("...", result)
		result = Shorten("hello", 2)
		a.Equal("...", result)
	})
}

func TestShortenWithSuffix(t *testing.T) {
	t.Run("使用自定义后缀截断", func(t *testing.T) {
		a := require.New(t)
		result := ShortenWithSuffix("hello world, this is a long text", 20, "[truncated]")
		a.Equal("hello wor[truncated]", result)
		a.Equal(20, len(result))
	})

	t.Run("字符串不超过长度返回原文", func(t *testing.T) {
		a := require.New(t)
		result := ShortenWithSuffix("hi", 10, "...")
		a.Equal("hi", result)
	})

	t.Run("maxLen 小于后缀长度", func(t *testing.T) {
		a := require.New(t)
		result := ShortenWithSuffix("hello", 3, "[truncated]")
		a.Equal("[truncated]", result)
	})
}
