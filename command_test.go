package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommandOutput_text(t *testing.T) {
	a := require.New(t)

	result, err := ParseCommandOutput("hello world")
	a.NoError(err)
	a.NotNil(result)
	a.Equal(COMMAND_OUTPUT_KIND_TEXT, result.Kind)
	a.Equal("hello world", result.Text)
}

func TestParseCommandOutputP_text(t *testing.T) {
	a := require.New(t)

	result := ParseCommandOutputP("test output")
	a.NotNil(result)
	a.Equal(COMMAND_OUTPUT_KIND_TEXT, result.Kind)
}

func TestParseCommandOutput_json(t *testing.T) {
	a := require.New(t)

	jsonOutput := "$json$\n\n{\"key\":\"value\"}"
	result, err := ParseCommandOutput(jsonOutput)
	a.NoError(err)
	a.NotNil(result)
	a.Equal(COMMAND_OUTPUT_KIND_JSON, result.Kind)
	a.NotNil(result.Json)
}

func TestParseCommandOutput_vars(t *testing.T) {
	a := require.New(t)

	varsOutput := "$vars$\n\nKEY1=value1\nKEY2=value2"
	result, err := ParseCommandOutput(varsOutput)
	a.NoError(err)
	a.NotNil(result)
	a.Equal(COMMAND_OUTPUT_KIND_VARS, result.Kind)
	a.Equal("value1", result.Vars["KEY1"])
	a.Equal("value2", result.Vars["KEY2"])
}

func TestParseCommandOutput_invalidJson(t *testing.T) {
	a := require.New(t)

	invalidJson := "$json$\n\n{invalid}"
	_, err := ParseCommandOutput(invalidJson)
	a.Error(err)
}

func TestVars2Pair_happy(t *testing.T) {
	a := require.New(t)

	vars := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result := Vars2Pair(vars)
	a.Equal(2, len(result))
}

func TestVars2Pair_empty(t *testing.T) {
	a := require.New(t)

	result := Vars2Pair(nil)
	a.Nil(result)

	result = Vars2Pair(map[string]string{})
	a.Nil(result)
}

func TestText2Vars_happy(t *testing.T) {
	a := require.New(t)

	text := "KEY1=value1\nKEY2=value2"
	result := Text2Vars(text)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}

func TestPair2Vars_happy(t *testing.T) {
	a := require.New(t)

	pairs := []string{"KEY1=value1", "KEY2=value2"}
	result := Pair2Vars(pairs)
	a.Equal("value1", result["KEY1"])
	a.Equal("value2", result["KEY2"])
}

func TestPair2Vars_empty(t *testing.T) {
	a := require.New(t)

	result := Pair2Vars(nil)
	a.NotNil(result)
	a.Empty(result)

	result = Pair2Vars([]string{})
	a.NotNil(result)
	a.Empty(result)
}

func TestPair2Vars_emptyValue(t *testing.T) {
	a := require.New(t)

	pairs := []string{"KEY="}
	result := Pair2Vars(pairs)
	a.Equal("", result["KEY"])
}

func TestPair2Vars_invalidPair(t *testing.T) {
	a := require.New(t)

	// Invalid pairs should be ignored
	pairs := []string{"invalid", "=nokey", "KEY=value"}
	result := Pair2Vars(pairs)
	a.Equal(1, len(result))
	a.Equal("value", result["KEY"])
}

func TestIsSudoCommand_true(t *testing.T) {
	a := require.New(t)

	a.True(IsSudoCommand("sudo apt-get install"))
	a.True(IsSudoCommand("sudo systemctl start"))
}

func TestIsSudoCommand_false(t *testing.T) {
	a := require.New(t)

	a.False(IsSudoCommand("apt-get install"))
	a.False(IsSudoCommand("echo hello"))
}

func TestInputSudoPassword_nil(t *testing.T) {
	a := require.New(t)

	result := InputSudoPassword(nil)
	a.Equal("", result)
}

func TestInputSudoPassword_empty(t *testing.T) {
	a := require.New(t)

	fn := func() string { return "" }
	result := InputSudoPassword(fn)
	a.Equal("", result)
}

func TestInputSudoPassword_value(t *testing.T) {
	a := require.New(t)

	fn := func() string { return "password" }
	result := InputSudoPassword(fn)
	a.Equal("password", result)
}

func TestInstrumentSudoCommand_notSudo(t *testing.T) {
	a := require.New(t)

	result := InstrumentSudoCommand("echo hello")
	a.Equal("echo hello", result)
}

func TestInstrumentSudoCommand_alreadyInstrumented(t *testing.T) {
	a := require.New(t)

	result := InstrumentSudoCommand("sudo -S apt-get install")
	a.Equal("sudo -S apt-get install", result)

	result = InstrumentSudoCommand("sudo --stdin apt-get install")
	a.Equal("sudo --stdin apt-get install", result)
}

func TestInstrumentSudoCommand_needsInstrumentation(t *testing.T) {
	a := require.New(t)

	result := InstrumentSudoCommand("sudo apt-get install")
	a.Equal("sudo --stdin apt-get install", result)
}

func TestRunCommandNoInput_echo(t *testing.T) {
	a := require.New(t)

	result, err := RunCommandNoInput(nil, "", "echo", "hello")
	a.NoError(err)
	a.NotNil(result)
	a.Equal(COMMAND_OUTPUT_KIND_TEXT, result.Kind)
	a.Contains(result.Text, "hello")
}

func TestRunCommandNoInputP_echo(t *testing.T) {
	a := require.New(t)

	result := RunCommandNoInputP(nil, "", "echo", "test")
	a.NotNil(result)
	a.Contains(result.Text, "test")
}
