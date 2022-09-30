package test

import (
	"strings"
	"testing"

	"github.com/fastgh/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_ReadText_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("", comm.ReadTextP(strings.NewReader("")))
	a.Equal("xyz", comm.ReadTextP(strings.NewReader("xyz")))
}

func Test_ReadLines_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]string{}, comm.ReadLines(strings.NewReader("")))
	a.Equal([]string{""}, comm.ReadLines(strings.NewReader("\n")))
	a.Equal([]string{"", ""}, comm.ReadLines(strings.NewReader("\n\r")))
	a.Equal([]string{"", ""}, comm.ReadLines(strings.NewReader("\n\r\n")))

	a.Equal([]string{"1", "2", "3"}, comm.ReadLines(strings.NewReader("1\n2\n3")))
	a.Equal([]string{"", "1", "2", "3"}, comm.ReadLines(strings.NewReader("\n1\n2\n3\n")))
}

func Test_Text2Lines_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]string{"A", "B", "C"}, comm.Text2Lines("A\nB\nC"))
}
