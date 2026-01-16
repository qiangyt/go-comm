package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// mockExternalPluginContext is a mock implementation of ExternalPluginContext
type mockExternalPluginContext struct {
	startCalled bool
	stopCalled  bool
}

func (m *mockExternalPluginContext) Init(logger Logger, fs afero.Fs, codeFile string) {
	// Mock implementation
}

func (m *mockExternalPluginContext) Start() any {
	m.startCalled = true
	return "started"
}

func (m *mockExternalPluginContext) Stop() any {
	m.stopCalled = true
	return "stopped"
}

func TestExternalPlugin_Name(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		name: "test-plugin",
	}
	a.Equal("test-plugin", plugin.Name())
}

func TestExternalPlugin_IsStarted(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		started: false,
	}
	a.False(plugin.IsStarted())

	plugin.started = true
	a.True(plugin.IsStarted())
}

func TestExternalPlugin_Start(t *testing.T) {
	a := require.New(t)

	ctx := &mockExternalPluginContext{}
	plugin := &ExternalPluginT{
		started: false,
		context: ctx,
	}

	plugin.Start(nil)
	a.True(plugin.started)
	a.True(ctx.startCalled)
}

func TestExternalPlugin_Start_alreadyStarted(t *testing.T) {
	a := require.New(t)

	ctx := &mockExternalPluginContext{}
	plugin := &ExternalPluginT{
		started: true,
		context: ctx,
	}

	plugin.Start(nil)
	a.True(plugin.started)
	// Should not call context.Start() again
	a.False(ctx.startCalled)
}

func TestExternalPlugin_Kind(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		kind: "test-kind",
	}
	a.Equal("test-kind", plugin.Kind())
}

func TestExternalPlugin_Stop(t *testing.T) {
	a := require.New(t)

	ctx := &mockExternalPluginContext{}
	plugin := &ExternalPluginT{
		started: true,
		context: ctx,
	}

	plugin.Stop(nil)
	a.False(plugin.started)
	a.True(ctx.stopCalled)
}

func TestExternalPlugin_Stop_notStarted(t *testing.T) {
	a := require.New(t)

	ctx := &mockExternalPluginContext{}
	plugin := &ExternalPluginT{
		started: false,
		context: ctx,
	}

	plugin.Stop(nil)
	a.False(plugin.started)
	// Should not call context.Stop()
	a.False(ctx.stopCalled)
}

func TestExternalPlugin_Version(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		versionMajor: 1,
		versionMinor: 2,
	}
	major, minor := plugin.Version()
	a.Equal(1, major)
	a.Equal(2, minor)
}

func TestExternalPlugin_Language(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		language: "go",
	}
	a.Equal("go", plugin.Language())
}

func TestExternalPlugin_Dir(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		dir: "/plugins/test",
	}
	a.Equal("/plugins/test", plugin.Dir())
}

func TestExternalPlugin_CodeFile(t *testing.T) {
	a := require.New(t)

	plugin := &ExternalPluginT{
		codeFile: "/plugins/test/plugin.go",
	}
	a.Equal("/plugins/test/plugin.go", plugin.CodeFile())
}

func TestResolveExternalPlugin_noManifest(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/plugins/test")
	logger := NewDiscardLogger()

	result := ResolveExternalPlugin(logger, fs, "/plugins/test")
	a.Nil(result)
}

func TestResolveExternalPlugin_withYamlManifest(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/plugins/test")
	WriteFileTextP(fs, "/plugins/test/plugin.manifest.yml", "kind: go_external\nname: test\nversion_major: 1\nversion_minor: 0")
	WriteFileTextP(fs, "/plugins/test/plugin.go", "package plugin\n\nfunc PluginStart() {}\nfunc PluginStop() {}")

	logger := NewDiscardLogger()

	result := ResolveExternalPlugin(logger, fs, "/plugins/test")
	a.NotNil(result)
	a.Equal("test", result.Name())
	major, minor := result.Version()
	a.Equal(1, major)
	a.Equal(0, minor)
}

func TestResolveExternalPlugin_noCodeFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/plugins/test")
	WriteFileTextP(fs, "/plugins/test/plugin.manifest.yml", "kind: go_external\nname: test\nversion_major: 1\nversion_minor: 0")
	// Missing plugin.go

	logger := NewDiscardLogger()

	result := ResolveExternalPlugin(logger, fs, "/plugins/test")
	a.Nil(result)
}

func TestListExternalPlugins_emptyDir(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/plugins")
	logger := NewDiscardLogger()

	result := ListExternalPlugins(logger, fs, "/plugins")
	a.NotNil(result)
	a.Empty(result)
}

func TestListExternalPlugins_withValidPlugins(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/plugins/test1")
	MkdirP(fs, "/plugins/test2")

	WriteFileTextP(fs, "/plugins/test1/plugin.manifest.yml", "kind: go_external\nname: plugin1\nversion_major: 1\nversion_minor: 0")
	WriteFileTextP(fs, "/plugins/test1/plugin.go", "package plugin\n\nfunc PluginStart() {}\nfunc PluginStop() {}")

	WriteFileTextP(fs, "/plugins/test2/plugin.manifest.yml", "kind: go_external\nname: plugin2\nversion_major: 1\nversion_minor: 0")
	WriteFileTextP(fs, "/plugins/test2/plugin.go", "package plugin\n\nfunc PluginStart() {}\nfunc PluginStop() {}")

	logger := NewDiscardLogger()

	result := ListExternalPlugins(logger, fs, "/plugins")
	a.Len(result, 2)
}
