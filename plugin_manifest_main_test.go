package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestPluginManifestWithMap_happy(t *testing.T) {
	a := require.New(t)

	manifestMap := map[string]any{
		"kind":          "go_external",
		"name":          "TestPlugin",
		"version_major": 1,
		"version_minor": 0,
	}

	manifest := PluginManifestWithMap(manifestMap)
	a.NotNil(manifest)
	a.Equal("go_external", manifest.Kind)
	a.Equal("testplugin", manifest.Name) // Should be lowercased
	a.Equal(1, manifest.VersionMajor)
	a.Equal(0, manifest.VersionMinor)
}

func TestPluginManifestWithMap_noKind(t *testing.T) {
	a := require.New(t)

	manifestMap := map[string]any{
		"name":          "TestPlugin",
		"version_major": 1,
		"version_minor": 0,
	}

	manifest := PluginManifestWithMap(manifestMap)
	a.NotNil(manifest)
	a.Equal("testplugin", manifest.Name)
}

func TestPluginManifestWithMap_nameLowercased(t *testing.T) {
	a := require.New(t)

	manifestMap := map[string]any{
		"name": "MYPLUGIN",
	}

	manifest := PluginManifestWithMap(manifestMap)
	a.Equal("myplugin", manifest.Name)
}

func TestPluginManifestWithJsonFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	jsonContent := `{"kind": "go_external", "name": "TestPlugin", "version_major": 1, "version_minor": 0}`
	WriteFileTextP(fs, "/manifest.json", jsonContent)

	manifest := PluginManifestWithJsonFile(fs, "/manifest.json")
	a.NotNil(manifest)
	a.Equal("testplugin", manifest.Name)
	a.Equal(1, manifest.VersionMajor)
}

func TestPluginManifestWithYamlFile_happy(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	yamlContent := `kind: go_external
name: TestPlugin
version_major: 1
version_minor: 0`
	WriteFileTextP(fs, "/manifest.yaml", yamlContent)

	manifest := PluginManifestWithYamlFile(fs, "/manifest.yaml")
	a.NotNil(manifest)
	a.Equal("testplugin", manifest.Name)
	a.Equal(1, manifest.VersionMajor)
}
