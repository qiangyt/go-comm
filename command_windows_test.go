package comm

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSudoCommand_Windows(t *testing.T) {
	a := require.New(t)
	_ = a

	if runtime.GOOS != "windows" {
		t.Skip("Skipping on non-Windows")
	}

	// RunSudoCommand on Windows returns "todo" error
	// We just want to call it to get coverage
	vars := map[string]string{}
	dir := ""
	cmd := "sudo some command"
	var passwordInput FnInput = nil

	_, err := RunSudoCommand(vars, dir, cmd, passwordInput)
	// Should return "todo" error
	a.Error(err)
	a.Contains(err.Error(), "todo")
}
