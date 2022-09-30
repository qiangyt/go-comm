package test

import (
	"testing"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fastgh/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_Slice2Set_happy(t *testing.T) {
	a := require.New(t)

	set := comm.Slice2Set("A", "A", "B")

	a.Equal(hashset.New("A", "B"), set)

	a.True(set.Contains("A"))
	a.True(set.Contains("B"))
	a.False(set.Contains("C"))
}

func Test_Slice2Map_happy(t *testing.T) {
	a := require.New(t)

	a.Equal(
		map[string]string{
			"1k": "1",
			"2k": "2",
		},
		comm.Slice2Map([]string{"1", "2"}, func(i string) string { return i + "k" }))
}
