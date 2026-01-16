package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunGoshCommandP_happy(t *testing.T) {
	a := require.New(t)

	// Test simple echo command
	result := RunGoshCommandP(nil, "", "echo hello", nil)
	a.NotNil(result)
	a.NotNil(result.Text)
}
