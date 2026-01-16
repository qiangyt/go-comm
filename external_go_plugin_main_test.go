package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewExternalGoPluginContext(t *testing.T) {
	a := require.New(t)

	ctx := NewExternalGoPluginContext()
	a.NotNil(ctx)
	a.Nil(ctx.interpreter)
	a.Nil(ctx.startFunc)
	a.Nil(ctx.stopFunc)
}

func TestExternalGoPluginContext_GetStartFunc(t *testing.T) {
	a := require.New(t)

	ctx := NewExternalGoPluginContext()
	result := ctx.GetStartFunc()
	a.Nil(result)
}

func TestExternalGoPluginContext_GetStopFunc(t *testing.T) {
	a := require.New(t)

	ctx := NewExternalGoPluginContext()
	result := ctx.GetStopFunc()
	a.Nil(result)
}

func TestExternalGoPluginContext_Start_noFunc(t *testing.T) {
	a := require.New(t)

	ctx := NewExternalGoPluginContext()
	result := ctx.Start()
	a.Equal("", result)
}

func TestExternalGoPluginContext_Stop_noFunc(t *testing.T) {
	a := require.New(t)

	ctx := NewExternalGoPluginContext()
	result := ctx.Stop()
	a.Equal("", result)
}

func TestExternalGoPluginContext_Init_invalidCode(t *testing.T) {
	fs := afero.NewMemMapFs()
	WriteFileTextP(fs, "/test.go", "invalid go code !!!")

	ctx := NewExternalGoPluginContext()
	logger := NewDiscardLogger()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Init should panic on invalid code")
		}
	}()

	ctx.Init(logger, fs, "/test.go")
}

func TestExternalGoPluginContext_Init_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	// Simple valid Go code
	code := `package plugin

func PluginStart() string {
	return "started"
}

func PluginStop() string {
	return "stopped"
}`
	WriteFileTextP(fs, "/test.go", code)

	ctx := NewExternalGoPluginContext()
	logger := NewDiscardLogger()

	// This should not panic
	ctx.Init(logger, fs, "/test.go")
	a.NotNil(ctx.interpreter)
}
