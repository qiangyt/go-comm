package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestNewLocalPluginLoader(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	logger := NewDiscardLogger()

	loader := NewLocalPluginLoader(logger, fs, "/plugins")
	a.NotNil(loader)
	a.Equal("local", loader.Namespace())
}

func TestNewRemotePluginLoader(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	logger := NewDiscardLogger()

	loader := NewRemotePluginLoader(logger, fs, "/plugins")
	a.NotNil(loader)
	a.Equal("remote", loader.Namespace())
}

func TestNewFsPluginLoader(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	logger := NewDiscardLogger()

	loader := NewFsPluginLoader(logger, fs, "/plugins", "test-namespace")
	a.NotNil(loader)
	a.Equal("test-namespace", loader.Namespace())
}
