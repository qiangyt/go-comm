package test

import (
	"testing"

	"github.com/qiangyt/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_EnvironMap_happy(t *testing.T) {
	a := require.New(t)

	actual := comm.EnvironMapP(nil)
	a.True(len(actual) > 0)
	a.True(len(actual["PATH"]) > 0)

	overrides := map[string]string{
		"k1":   "v1",
		"k2":   "v2",
		"PATH": "overrided_path",
	}

	actual = comm.EnvironMapP(overrides)
	a.Equal("v1", actual["k1"])
	a.Equal("v2", actual["k2"])
	a.Equal("overrided_path", actual["PATH"])
}

func Test_AbsPath_happy(t *testing.T) {
	a := require.New(t)

	if comm.IsWindows() {
		a.Equal("\\home", comm.AbsPath("/home", "."))
		a.Equal("\\home", comm.AbsPath("/home", "./"))
		a.Equal("\\", comm.AbsPath("/home", ".."))
		a.Equal("\\", comm.AbsPath("/home", "../"))

		a.Equal("\\home\\1\\2", comm.AbsPath("/home/1", "2"))
		a.Equal("\\home\\1\\2\\3", comm.AbsPath("/home/1", "2/3"))
		a.Equal("\\home\\1\\2\\3", comm.AbsPath("/home/1", "./2/3"))
		a.Equal("\\home\\2\\3", comm.AbsPath("/home/1", "../2/3"))

		a.Equal("\\home\\1\\2\\3", comm.AbsPath("/home/1", "2/./3"))
		a.Equal("\\home\\1\\3", comm.AbsPath("/home/1", "2/../3"))
	} else {
		a.Equal("/home", comm.AbsPath("/home", "."))
		a.Equal("/home", comm.AbsPath("/home", "./"))
		a.Equal("/", comm.AbsPath("/home", ".."))
		a.Equal("/", comm.AbsPath("/home", "../"))

		a.Equal("/home/1/2", comm.AbsPath("/home/1", "2"))
		a.Equal("/home/1/2/3", comm.AbsPath("/home/1", "2/3"))
		a.Equal("/home/1/2/3", comm.AbsPath("/home/1", "./2/3"))
		a.Equal("/home/2/3", comm.AbsPath("/home/1", "../2/3"))

		a.Equal("/home/1/2/3", comm.AbsPath("/home/1", "2/./3"))
		a.Equal("/home/1/3", comm.AbsPath("/home/1", "2/../3"))
	}
}
