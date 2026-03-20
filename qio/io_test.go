package qio

import (
	"strings"
	"testing"

	"github.com/qiangyt/go-comm/v2/qlang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReadText_happy(t *testing.T) {
	a := require.New(t)

	a.Equal("", ReadTextP(strings.NewReader("")))
	a.Equal("xyz", ReadTextP(strings.NewReader("xyz")))
}

func Test_ReadLines_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]string{}, ReadLines(strings.NewReader("")))
	a.Equal([]string{""}, ReadLines(strings.NewReader("\n")))
	a.Equal([]string{"", ""}, ReadLines(strings.NewReader("\n\r")))
	a.Equal([]string{"", ""}, ReadLines(strings.NewReader("\n\r\n")))

	a.Equal([]string{"1", "2", "3"}, ReadLines(strings.NewReader("1\n2\n3")))
	a.Equal([]string{"", "1", "2", "3"}, ReadLines(strings.NewReader("\n1\n2\n3\n")))
}

func Test_Text2Lines_happy(t *testing.T) {
	a := require.New(t)

	a.Equal([]string{"A", "B", "C"}, qlang.Text2Lines("A\nB\nC"))
}

// CloseQuietly 测试

type mockCloser struct {
	closed bool
	err    error
}

func (m *mockCloser) Close() error {
	m.closed = true
	if m.err != nil {
		return m.err
	}
	return nil
}

func Test_CloseQuietly_happy(t *testing.T) {
	a := require.New(t)

	// 正常关闭
	mc := &mockCloser{}
	CloseQuietly(mc)
	a.True(mc.closed)

	// nil closer 不应该 panic
	CloseQuietly(nil)
}

func Test_CloseQuietly_withError(t *testing.T) {
	a := require.New(t)

	// Close 返回错误时应该忽略错误
	mc := &mockCloser{err: assert.AnError}
	CloseQuietly(mc)
	a.True(mc.closed)
}
