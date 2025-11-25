package comm

import (
	"testing"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/stretchr/testify/require"
)

func TestSet2Strings_happy(t *testing.T) {
	a := require.New(t)

	set := hashset.New("apple", "banana", "cherry")
	result := Set2Strings(set)

	a.Len(result, 3)
	a.Contains(result, "apple")
	a.Contains(result, "banana")
	a.Contains(result, "cherry")
}

func TestSet2Strings_empty(t *testing.T) {
	a := require.New(t)

	set := hashset.New()
	result := Set2Strings(set)

	a.Empty(result)
}

func TestSlice2Map_happy(t *testing.T) {
	a := require.New(t)

	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	result := Slice2Map(people, func(p Person) string { return p.Name })

	a.Len(result, 2)
	a.Equal(30, result["Alice"].Age)
	a.Equal(25, result["Bob"].Age)
}

func TestDowncastMap_happy(t *testing.T) {
	a := require.New(t)

	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	result := DowncastMap(m)

	a.Len(result, 3)
	a.Equal(1, result["one"])
	a.Equal(2, result["two"])
	a.Equal(3, result["three"])
}

func TestSliceEquals_equal(t *testing.T) {
	a := require.New(t)

	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 3}

	a.True(SliceEquals(a1, a2))
}

func TestSliceEquals_notEqual(t *testing.T) {
	a := require.New(t)

	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 4}

	a.False(SliceEquals(a1, a2))
}

func TestSliceEquals_differentLength(t *testing.T) {
	a := require.New(t)

	a1 := []int{1, 2, 3}
	a2 := []int{1, 2}

	a.False(SliceEquals(a1, a2))
}

func TestSliceEquals_oneNil(t *testing.T) {
	a := require.New(t)

	var a1 []int = nil
	a2 := []int{1, 2, 3}

	a.False(SliceEquals(a1, a2))
}

func TestSliceEquals_bothNil(t *testing.T) {
	a := require.New(t)

	var a1 []int = nil
	var a2 []int = nil

	a.True(SliceEquals(a1, a2))
}

func TestMergeMap_simple(t *testing.T) {
	a := require.New(t)

	m1 := map[string]string{
		"a": "1",
		"b": "2",
	}

	m2 := map[string]string{
		"b": "3",
		"c": "4",
	}

	result := MergeMap(m1, m2)

	a.Len(result, 3)
	a.Equal("1", result["a"])
	a.Equal("3", result["b"]) // m2 overwrites m1
	a.Equal("4", result["c"])
}

func TestMergeMap_multiple(t *testing.T) {
	a := require.New(t)

	m1 := map[string]int{"x": 10}
	m2 := map[string]int{"y": 20}
	m3 := map[string]int{"z": 30}

	result := MergeMap(m1, m2, m3)

	a.Len(result, 3)
	a.Equal(10, result["x"])
	a.Equal(20, result["y"])
	a.Equal(30, result["z"])
}
