package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPluginLoader(t *testing.T) {
	loader := NewPluginLoader("test-namespace")
	a := require.New(t)

	a.NotNil(loader)
	a.Equal("test-namespace", loader.Namespace())
	a.NotNil(loader.Plugins())
	a.Empty(loader.Plugins())
}

func TestPluginLoader_Register(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin := &mockPlugin{name: "plugin1"}
	loader.Register(plugin)

	a.Len(loader.Plugins(), 1)
	a.Equal("plugin1", loader.Plugins()["plugin1"].Name())
}

func TestPluginLoader_Register_duplicate(t *testing.T) {
	loader := NewPluginLoader("test-ns")

	plugin1 := &mockPlugin{name: "plugin1"}
	plugin2 := &mockPlugin{name: "plugin1"}

	loader.Register(plugin1)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Register should panic on duplicate plugin")
		}
	}()

	loader.Register(plugin2)
}

func TestPluginLoader_RegisterThenStart(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin := &mockPlugin{name: "plugin1"}
	logger := NewDiscardLogger()

	loader.RegisterThenStart(logger, plugin)

	a.True(plugin.started)
}

func TestPluginLoader_Start(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin1 := &mockPlugin{name: "plugin1"}
	plugin2 := &mockPlugin{name: "plugin2"}
	loader.Register(plugin1)
	loader.Register(plugin2)

	logger := NewDiscardLogger()
	err := loader.Start(logger)
	a.NoError(err)
	a.True(plugin1.started)
	a.True(plugin2.started)
}

func TestPluginLoader_Start_alreadyStarted(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin := &mockPlugin{name: "plugin1"}
	loader.Register(plugin)

	logger := NewDiscardLogger()
	err := loader.Start(logger)
	a.NoError(err)

	// Start again should not error
	err = loader.Start(logger)
	a.NoError(err)
}

func TestPluginLoader_Stop(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin1 := &mockPlugin{name: "plugin1"}
	plugin2 := &mockPlugin{name: "plugin2"}
	loader.Register(plugin1)
	loader.Register(plugin2)

	logger := NewDiscardLogger()
	err := loader.Start(logger)
	a.NoError(err)

	err = loader.Stop(logger)
	a.NoError(err)
	a.True(plugin1.stopped)
	a.True(plugin2.stopped)
}

func TestPluginLoader_Stop_notStarted(t *testing.T) {
	loader := NewPluginLoader("test-ns")
	a := require.New(t)

	plugin := &mockPlugin{name: "plugin1"}
	loader.Register(plugin)

	logger := NewDiscardLogger()
	err := loader.Stop(logger)
	a.NoError(err)
}

// Mock plugin for testing
type mockPlugin struct {
	name    string
	started bool
	stopped bool
}

func (m *mockPlugin) Name() string {
	return m.name
}

func (m *mockPlugin) Kind() PluginKind {
	return "go_external"
}

func (m *mockPlugin) Version() (major int, minor int) {
	return 1, 0
}

func (m *mockPlugin) Start(logger Logger) {
	m.started = true
}

func (m *mockPlugin) Stop(logger Logger) {
	m.stopped = true
}

func (m *mockPlugin) IsStarted() bool {
	return m.started
}
