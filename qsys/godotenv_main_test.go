package qsys

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestReadEnv_noFile(t *testing.T) {
	a := require.New(t)

	fs := afero.NewMemMapFs()

	_, err := ReadEnv(fs)
	a.Error(err)
}

func TestMarshalEnv_happy(t *testing.T) {
	a := require.New(t)

	envMap := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	result, err := MarshalEnv(envMap)
	a.NoError(err)
	a.NotEmpty(result)
	// Should be sorted
	lines := strings.Split(result, "\n")
	a.True(lines[0] < lines[1])
}

func TestMarshalEnv_numeric(t *testing.T) {
	a := require.New(t)

	envMap := map[string]string{
		"PORT": "8080",
	}

	result, err := MarshalEnv(envMap)
	a.NoError(err)
	a.Contains(result, "PORT=8080")
}

func TestEnvFilenamesOrDefault_empty(t *testing.T) {
	result := envFilenamesOrDefault([]string{})
	a := require.New(t)
	a.Equal([]string{".env"}, result)
}

func TestEnvFilenamesOrDefault_provided(t *testing.T) {
	a := require.New(t)

	input := []string{".env1", ".env2"}
	result := envFilenamesOrDefault(input)
	a.Equal(input, result)
}

func TestEnvDoubleQuoteEscape_happy(t *testing.T) {
	a := require.New(t)

	// Test escaping of special characters
	result := envDoubleQuoteEscape("test\nvalue")
	a.Equal(`test\nvalue`, result)
}
