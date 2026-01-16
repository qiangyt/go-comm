package comm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOutput(t *testing.T) {
	result := DefaultOutput()
	a := require.New(t)
	a.NotNil(result)
}

func TestReadBytes_happy(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("test content")
	result, err := ReadBytes(reader)
	a.NoError(err)
	a.Equal([]byte("test content"), result)
}

func TestReadBytesP_happy(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("test content")
	result := ReadBytesP(reader)
	a.Equal([]byte("test content"), result)
}

func TestReadText_happy(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("test content")
	result, err := ReadText(reader)
	a.NoError(err)
	a.Equal("test content", result)
}

func TestReadTextP_happy(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("test content")
	result := ReadTextP(reader)
	a.Equal("test content", result)
}

func TestReadLines_happy(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("line1\nline2\nline3")
	result := ReadLines(reader)
	a.Equal(3, len(result))
	a.Equal("line1", result[0])
	a.Equal("line2", result[1])
	a.Equal("line3", result[2])
}

func TestReadLines_empty(t *testing.T) {
	a := require.New(t)

	reader := strings.NewReader("")
	result := ReadLines(reader)
	a.NotNil(result)
	a.Empty(result)
}
