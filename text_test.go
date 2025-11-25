package comm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		"keep": "original",
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
