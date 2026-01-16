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

func TestRunShellCommand_gosh(t *testing.T) {
	a := require.New(t)

	// Test with gosh shell
	result, err := RunShellCommand(nil, "", "gosh", "echo hello", nil)
	a.NoError(err)
	a.NotNil(result)
	// gosh should return output containing our text
}

func TestRunShellCommand_sudo(t *testing.T) {
	a := require.New(t)

	// Test with sudo command (should instrument)
	result, err := RunShellCommand(nil, "", "", "sudo echo test", nil)
	// This might fail on systems without sudo, so we check for expected behavior
	if err != nil {
		t.Logf("Sudo command not available: %v", err)
	} else {
		a.NotNil(result)
	}
}

func TestRunShellCommandP(t *testing.T) {
	a := require.New(t)

	// Test panic version
	result := RunShellCommandP(nil, "", "gosh", "echo test", nil)
	a.NotNil(result)
}

func TestRunUserCommandP(t *testing.T) {
	// Test panic version with invalid command - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("RunUserCommandP should panic on error")
		}
	}()
	RunUserCommandP(nil, "", "invalid-command-xyz-123")
}

func TestRunUserCommand_invalid(t *testing.T) {
	a := require.New(t)

	// Test with invalid command
	_, err := RunUserCommand(nil, "", "invalid-command-xyz-123")
	a.Error(err)
}

func TestRunCommandWithInput_echo(t *testing.T) {
	a := require.New(t)

	cmdFn := RunCommandWithInput(nil, "", "echo", "arg1")
	result, err := cmdFn()
	a.NoError(err)
	a.NotNil(result)
	a.Contains(result.Text, "arg1")
}

func TestRunCommandWithInput_sudo(t *testing.T) {
	a := require.New(t)

	// Test sudo command instrumentation
	cmdFn := RunCommandWithInput(nil, "", "sudo", "echo", "test")
	result, err := cmdFn("password")

	// This might fail on systems without sudo, check for expected behavior
	if err != nil {
		t.Logf("Sudo command not available: %v", err)
	} else {
		a.NotNil(result)
	}
}

func TestRunCommandWithInput_withArgs(t *testing.T) {
	a := require.New(t)

	// Test with multiple arguments
	cmdFn := RunCommandWithInput(nil, "", "echo", "arg1", "arg2", "arg3")
	result, err := cmdFn()
	a.NoError(err)
	a.NotNil(result)
}

func TestNewExecCommand_happy(t *testing.T) {
	a := require.New(t)

	// Test creating exec command with vars
	vars := map[string]string{
		"TEST_VAR": "test_value",
	}

	cmd, err := newExecCommand(vars, "", "echo", "hello")
	a.NoError(err)
	a.NotNil(cmd)
	// On Windows, cmd.Path is the full path to echo.exe
	a.True(cmd.Path != "")
}
