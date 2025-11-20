package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Mock plugin for testing
type MockPlugin struct {
	name          string
	kind          PluginKind
	major         int
	minor         int
	startCalled   bool
	stopCalled    bool
	shouldPanic   bool
	panicMessage  string
}

func (m *MockPlugin) Name() string {
	return m.name
}

func (m *MockPlugin) Kind() PluginKind {
	return m.kind
}

func (m *MockPlugin) Version() (major int, minor int) {
	return m.major, m.minor
}

func (m *MockPlugin) Start(logger Logger) {
	if m.shouldPanic {
		panic(m.panicMessage)
	}
	m.startCalled = true
}

func (m *MockPlugin) Stop(logger Logger) {
	if m.shouldPanic {
		panic(m.panicMessage)
	}
	m.stopCalled = true
}

// Mock plugin loader for testing
type MockPluginLoader struct {
	namespace    string
	plugins      map[string]Plugin
	startCalled  bool
	stopCalled   bool
	shouldError  bool
}

func (m *MockPluginLoader) Namespace() string {
	return m.namespace
}

func (m *MockPluginLoader) Plugins() map[string]Plugin {
	return m.plugins
}

func (m *MockPluginLoader) Start(logger Logger) error {
	m.startCalled = true
	if m.shouldError {
		return LocalizeError("error.plugin.start_failed", map[string]interface{}{
			"PluginId": m.namespace + "/test",
			"Version":  "1.0",
			"Cause":    "test error",
		})
	}
	return nil
}

func (m *MockPluginLoader) Stop(logger Logger) error {
	m.stopCalled = true
	if m.shouldError {
		return LocalizeError("error.plugin.stop_failed", map[string]interface{}{
			"PluginId": m.namespace + "/test",
			"Version":  "1.0",
			"Cause":    "test error",
		})
	}
	return nil
}

func TestPluginId(t *testing.T) {
	a := require.New(t)

	id := PluginId("test-namespace", "test-plugin")
	a.Equal("test-namespace/test-plugin", id)
}

func TestStartPlugin_Success(t *testing.T) {
	a := require.New(t)

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	logger := NewDiscardLogger()
	err := StartPlugin("test-namespace", plugin, logger)

	a.NoError(err)
	a.True(plugin.startCalled)
}

func TestStartPlugin_Panic(t *testing.T) {
	a := require.New(t)

	plugin := &MockPlugin{
		name:         "test-plugin",
		kind:         "test-kind",
		major:        1,
		minor:        0,
		shouldPanic:  true,
		panicMessage: "test panic",
	}

	logger := NewDiscardLogger()
	err := StartPlugin("test-namespace", plugin, logger)

	a.Error(err)
	a.Contains(err.Error(), "test-namespace/test-plugin")
	a.Contains(err.Error(), "1/0")
}

func TestStopPlugin_Success(t *testing.T) {
	a := require.New(t)

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	logger := NewDiscardLogger()
	err := StopPlugin("test-namespace", plugin, logger)

	a.NoError(err)
	a.True(plugin.stopCalled)
}

func TestStopPlugin_Panic(t *testing.T) {
	a := require.New(t)

	plugin := &MockPlugin{
		name:         "test-plugin",
		kind:         "test-kind",
		major:        1,
		minor:        0,
		shouldPanic:  true,
		panicMessage: "test panic",
	}

	logger := NewDiscardLogger()
	err := StopPlugin("test-namespace", plugin, logger)

	a.Error(err)
	a.Contains(err.Error(), "test-namespace/test-plugin")
	a.Contains(err.Error(), "1/0")
}

func TestNewPluginRegistry(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "kind1", "kind2")

	a.NotNil(registry)
	a.Equal(1, registry.SupportedMajorVersion())
	a.True(registry.IsSupportedPluginKind("kind1"))
	a.True(registry.IsSupportedPluginKind("kind2"))
	a.False(registry.IsSupportedPluginKind("kind3"))
}

func TestPluginRegistry_ValidatePlugin_Success(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	err := registry.ValidatePlugin("test-namespace", plugin)
	a.NoError(err)
}

func TestPluginRegistry_ValidatePlugin_VersionMismatch(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	registry := NewPluginRegistry(2, "test-kind")

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	err := registry.ValidatePlugin("test-namespace", plugin)
	a.Error(err)
	a.Contains(err.Error(), "test-namespace/test-plugin")
	a.Contains(err.Error(), "version")
}

func TestPluginRegistry_ValidatePlugin_UnsupportedKind(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	registry := NewPluginRegistry(1, "other-kind")

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	err := registry.ValidatePlugin("test-namespace", plugin)
	a.Error(err)
	a.Contains(err.Error(), "test-namespace/test-plugin")
	a.Contains(err.Error(), "unsupported")
}

func TestPluginRegistry_Register_Success(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	plugin := &MockPlugin{
		name:  "test-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	loader := &MockPluginLoader{
		namespace: "test-namespace",
		plugins: map[string]Plugin{
			"test-plugin": plugin,
		},
	}

	registry.Register(loader)

	a.True(registry.HasNamespace("test-namespace"))

	byKind := registry.ByKind("test-kind")
	a.NotNil(byKind)
	a.Equal(plugin, byKind["test-plugin"])
}

func TestPluginRegistry_Register_EmptyNamespace(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	loader := &MockPluginLoader{
		namespace: "",
		plugins:   map[string]Plugin{},
	}

	a.Panics(func() {
		registry.Register(loader)
	})
}

func TestPluginRegistry_Register_DuplicateNamespace(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	plugin1 := &MockPlugin{
		name:  "plugin1",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	loader1 := &MockPluginLoader{
		namespace: "test-namespace",
		plugins: map[string]Plugin{
			"plugin1": plugin1,
		},
	}

	registry.Register(loader1)

	// Try to register the same namespace again
	plugin2 := &MockPlugin{
		name:  "plugin2",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	loader2 := &MockPluginLoader{
		namespace: "test-namespace",
		plugins: map[string]Plugin{
			"plugin2": plugin2,
		},
	}

	a.Panics(func() {
		registry.Register(loader2)
	})
}

func TestPluginRegistry_Register_DuplicateKind(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	plugin1 := &MockPlugin{
		name:  "same-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	loader1 := &MockPluginLoader{
		namespace: "namespace1",
		plugins: map[string]Plugin{
			"same-plugin": plugin1,
		},
	}

	registry.Register(loader1)

	// Try to register a plugin with the same name and kind in a different namespace
	plugin2 := &MockPlugin{
		name:  "same-plugin",
		kind:  "test-kind",
		major: 1,
		minor: 0,
	}

	loader2 := &MockPluginLoader{
		namespace: "namespace2",
		plugins: map[string]Plugin{
			"same-plugin": plugin2,
		},
	}

	a.Panics(func() {
		registry.Register(loader2)
	})
}

func TestPluginRegistry_Init(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	loader := &MockPluginLoader{
		namespace: "test-namespace",
		plugins:   map[string]Plugin{},
	}

	registry.Register(loader)

	logger := NewDiscardLogger()
	registry.Init(logger)

	a.True(loader.startCalled)
}

func TestPluginRegistry_Destroy(t *testing.T) {
	a := require.New(t)

	registry := NewPluginRegistry(1, "test-kind")

	loader := &MockPluginLoader{
		namespace: "test-namespace",
		plugins:   map[string]Plugin{},
	}

	registry.Register(loader)

	logger := NewDiscardLogger()
	registry.Destroy(logger)

	a.True(loader.stopCalled)
	a.False(registry.HasNamespace("test-namespace"))
}
