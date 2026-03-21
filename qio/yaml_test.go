package qio

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func Test_MapFromYamlFileP_happy(t *testing.T) {
	a := require.New(t)
	fs := afero.NewMemMapFs()

	WriteFileTextP(fs, "test.yaml", `k: v`)

	configMap := MapFromYamlFileP(fs, "test.yaml", false)

	a.Len(configMap, 1)
	a.Equal("v", configMap["k"])
}
