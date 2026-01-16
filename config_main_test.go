package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestStrictConfigConfig(t *testing.T) {
	cfg := StrictConfigConfig()
	a := require.New(t)

	a.True(cfg.ErrorUnused)
	a.True(cfg.ErrorUnset)
	a.True(cfg.ZeroFields)
	a.False(cfg.WeaklyTypedInput)
	a.False(cfg.Squash)
	a.True(cfg.IgnoreUntaggedFields)
	a.False(cfg.DoValidate)
}

func TestDynamicConfigConfig(t *testing.T) {
	cfg := DynamicConfigConfig()
	a := require.New(t)

	a.False(cfg.ErrorUnused)
	a.False(cfg.ErrorUnset)
	a.False(cfg.ZeroFields)
	a.True(cfg.WeaklyTypedInput)
	a.True(cfg.Squash)
	a.True(cfg.IgnoreUntaggedFields)
	a.False(cfg.DoValidate)
}

func TestConfigConfig_ToMapstruct(t *testing.T) {
	cfg := &ConfigConfig{
		ErrorUnused:          true,
		ErrorUnset:           true,
		ZeroFields:           true,
		WeaklyTypedInput:     false,
		Squash:               false,
		IgnoreUntaggedFields: true,
	}
	msCfg := cfg.ToMapstruct()
	a := require.New(t)

	a.True(msCfg.ErrorUnused)
	a.True(msCfg.ErrorUnset)
	a.True(msCfg.ZeroFields)
	a.False(msCfg.WeaklyTypedInput)
	a.False(msCfg.Squash)
	a.True(msCfg.IgnoreUntaggedFields)
}

func TestDecodeWithYaml_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	var result Config
	yamlText := "name: test\nage: 30"

	r, metadata, err := DecodeWithYaml(yamlText, DynamicConfigConfig(), &result, nil)
	a.NoError(err)
	a.NotNil(r)
	a.NotNil(metadata)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestDecodeWithYamlP_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `mapstructure:"name"`
	}

	var result Config
	yamlText := "name: test"

	r, metadata := DecodeWithYamlP[Config](yamlText, DynamicConfigConfig(), &result, nil)
	a.NotNil(r)
	a.NotNil(metadata)
	a.Equal("test", result.Name)
}

func TestDecodeWithYamlP_invalidYaml(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("DecodeWithYamlP should panic on invalid YAML")
		}
	}()

	type Config struct {
		Name string `mapstructure:"name"`
	}

	var result Config
	DecodeWithYamlP[Config]("name: [invalid", DynamicConfigConfig(), &result, nil)
}

func TestDecodeWithMap_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	var result Config
	input := map[string]any{
		"name": "test",
		"age":  25,
	}

	r, metadata, err := DecodeWithMap(input, DynamicConfigConfig(), &result, nil)
	a.NoError(err)
	a.NotNil(r)
	a.NotNil(metadata)
	a.Equal("test", result.Name)
	a.Equal(25, result.Age)
}

func TestDecodeWithMap_withDefaults(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}

	var result Config
	input := map[string]any{
		"name": "test",
	}
	defaults := map[string]any{
		"age": 30,
	}

	r, metadata, err := DecodeWithMap(input, DynamicConfigConfig(), &result, defaults)
	a.NoError(err)
	a.NotNil(r)
	a.NotNil(metadata)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestDecodeWithMapP_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `mapstructure:"name"`
	}

	var result Config
	input := map[string]any{
		"name": "test",
	}

	r, metadata := DecodeWithMapP[Config](input, DynamicConfigConfig(), &result, nil)
	a.NotNil(r)
	a.NotNil(metadata)
	a.Equal("test", result.Name)
}

func TestGetMapValue_withKey(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"name": "test",
	}

	result := GetMapValue(m, "name", func() string { return "default" })
	a.Equal("test", result)
}

func TestGetMapValue_withoutKey(t *testing.T) {
	a := require.New(t)

	m := map[string]any{}

	defaultValue := "default"
	result := GetMapValue(m, "name", func() string { return defaultValue })
	a.Equal(defaultValue, result)

	// Key should be added to map
	a.Equal(defaultValue, m["name"])
}

func TestSysEnvFileNames_emptyShell(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "")
	a.NotNil(result)
}

func TestSysEnvFileNames_withZsh(t *testing.T) {
	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "zsh")
	a := require.New(t)
	a.NotNil(result)
}

func TestSysEnvFileNames_withBash(t *testing.T) {
	fs := afero.NewMemMapFs()
	result := SysEnvFileNames(fs, "bash")
	a := require.New(t)
	a.NotNil(result)
}

func TestLoadEnvScripts_noFilenames(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}

	result, err := LoadEnvScripts(fs, vars)
	a.NoError(err)
	a.NotNil(result)
}

func TestLoadEnvScripts_withFilenames(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}
	// Use non-existent files - RunGoshCommand will error, but LoadEnvScripts should still return with error
	result, err := LoadEnvScripts(fs, vars, "/nonexistent/file1", "/nonexistent/file2")
	a.Error(err)  // Expected to have errors
	a.NotNil(result)
}

func TestLoadEnvScript_etcPaths(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	MkdirP(fs, "/etc")
	WriteFileLinesP(fs, "/etc/paths", "/usr/bin", "/usr/local/bin")

	vars := map[string]string{"PATH": "/existing/path"}
	result, err := LoadEnvScript(fs, vars, "/etc/paths")
	a.NoError(err)
	a.NotNil(result)
	a.Contains(result["PATH"], "/existing/path")
	a.Contains(result["PATH"], "/usr/bin")
}

func TestLoadEnvScript_nonExistent(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	vars := map[string]string{}

	result, err := LoadEnvScript(fs, vars, "/nonexistent/file")
	a.Error(err)
	a.NotNil(result)
}
