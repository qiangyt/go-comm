package test

import (
	"testing"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/fastgh/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_Set2Strings_happy(t *testing.T) {
	a := require.New(t)

	strs := comm.Set2Strings(hashset.New("A", "A", "B"))
	a.Len(strs, 2)

	a.Equal("A", strs[0])
	a.Equal("B", strs[1])
}

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

func Test_MergeMap_happy(t *testing.T) {
	a := require.New(t)

	mapA := map[string]any{
		"key-A-1": "value-A-1",
		"key-A-2": "value-A-2",
		"key-A-3": map[string]any{
			"key-A-3-1": "value-A-3-1",
			"key-A-3-2": "value-A-3-2",
		},
	}
	mapB := map[string]any{
		"key-B-1": "value-B-1",
		"key-A-2": "value-B-2",
		"key-A-3": map[string]any{
			"key-B-3-1": "value-B-3-1",
			"key-A-3-2": "value-B-3-2",
		},
	}

	r := comm.MergeMap(mapA, mapB)
	a.Len(r, 4)

	a.Equal(r["key-A-1"], "value-A-1")
	a.Equal(r["key-A-2"], "value-B-2")
	a.Equal(r["key-B-1"], "value-B-1")

	A3 := r["key-A-3"].(map[string]any)
	a.Len(A3, 3)
	a.Equal(A3["key-A-3-1"], "value-A-3-1")
	a.Equal(A3["key-A-3-2"], "value-B-3-2")
	a.Equal(A3["key-B-3-1"], "value-B-3-1")
}
