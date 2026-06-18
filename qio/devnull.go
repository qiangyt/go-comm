// copy of github.com/go-task/task/v3/internal/execext/devnull.go
package qio

import (
	"io"
)

var _ io.ReadWriteCloser = DevNull{}

type DevNull struct{}

func (DevNull) Read(p []byte) (int, error)  { return 0, io.EOF }
func (DevNull) Write(p []byte) (int, error) { return len(p), nil }
func (DevNull) Close() error                { return nil }
