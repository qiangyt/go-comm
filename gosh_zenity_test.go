package comm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecZenity_UnknownSubCommand(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	_ = ctx
	_ = a
	t.Skip("Requires proper HandlerContext")
}

func TestZenityError_WithText(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	args := []string{"Error message"}
	_ = ctx
	_ = args
	_ = a
	t.Skip("zenity requires GUI")
}

func TestZenityInfo_WithText(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	args := []string{"Info message"}
	_ = ctx
	_ = args
	_ = a
	t.Skip("zenity requires GUI")
}

func TestZenityWarning_WithText(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	args := []string{"Warning message"}
	_ = ctx
	_ = args
	_ = a
	t.Skip("zenity requires GUI")
}

func TestZenityQuestion_WithText(t *testing.T) {
	a := require.New(t)

	ctx := context.Background()
	args := []string{"Question?"}
	_ = ctx
	_ = args
	_ = a
	t.Skip("zenity requires GUI")
}
