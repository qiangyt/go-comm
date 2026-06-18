package qio

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevNull_Read(t *testing.T) {
	a := require.New(t)

	var dn DevNull
	buf := make([]byte, 10)

	n, err := dn.Read(buf)
	a.Equal(0, n)
	a.Equal(io.EOF, err)
}

func TestDevNull_Write(t *testing.T) {
	a := require.New(t)

	var dn DevNull
	data := []byte("hello world")

	n, err := dn.Write(data)
	a.Equal(len(data), n)
	a.NoError(err)
}

func TestDevNull_Close(t *testing.T) {
	a := require.New(t)

	var dn DevNull
	err := dn.Close()
	a.NoError(err)
}

func TestDevNull_Interface(t *testing.T) {
	a := require.New(t)

	// Test that DevNull implements io.ReadWriteCloser
	var _ io.ReadWriteCloser = DevNull{}
	a.True(true) // If we got here, the interface is implemented
}
