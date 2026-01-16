package comm

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDefaultLocalizeFunc(t *testing.T) {
	a := require.New(t)

	result := DefaultLocalizeFunc("test.message", nil)
	a.Equal("test.message", result)

	result = DefaultLocalizeFunc("test.message", map[string]any{"key": "value"})
	a.Equal("test.message", result)
}

func TestCommLocalize(t *testing.T) {
	a := require.New(t)

	InitI18n("en")
	result := CommLocalize("test.message")
	a.NotEmpty(result)
}

func TestSetCommLang(t *testing.T) {
	// Just ensure it doesn't panic
	SetCommLang("zh")
	SetCommLang("en")
}

func TestNewLocalizedFileOps(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()
	ops := NewLocalizedFileOps(fs)
	a.NotNil(ops)
}
