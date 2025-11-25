package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 7,
	}

	// Test that IntP panics when the specified key does not exist in the map.
	// Test that IntP returns the value of the specified key when it exists in the map.
	r := RequiredIntP("task", "key1", m)
	a.Equal(7, r)

	// Test that IntP returns an error when the specified key does not exist in the map
	a.Panics(func() { RequiredIntP("task", "not-existed-key", m) })
}

func TestIntP_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-int-value",
	}

	// Test that failed to parse int value
	a.Panics(func() { RequiredIntP("task", "key1", m) })
}

func TestInt(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 5,
	}

	// Test that IntP returns the value of the specified key when it exists in the map.
	r, err := RequiredInt("task", "key1", m)
	a.NoError(err)
	a.Equal(5, r)

	// Test that IntP returns an error when the specified key does not exist in the map
	_, err = RequiredInt("task", "not-existed-key", m)
	a.Error(err)
}

func TestIntP_direct(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	result := IntP("test", 42)
	a.Equal(42, result)

	// Test panic on invalid type
	a.Panics(func() { IntP("test", "not an int") })
}

func TestInt_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-int-value",
	}

	// Test that failed to parse int value
	_, err := RequiredInt("task", "key1", m)
	a.Error(err)
}

func TestOptionalIntP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-int-value",
	}

	// Test that OptionalIntP panics when the specified key exists in the map but cannot be parsed as a int.
	a.Panics(func() { OptionalIntP("task", "key1", m, 9) })
}

func TestOptionalIntP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 10,
	}

	// Test that OptionalIntP returns the value of the specified key when it exists in the map.
	r, has := OptionalIntP("test", "key1", m, 2)
	a.True(has)
	a.Equal(10, r)

	// Test that OptionalIntP returns the default value when the specified key does not exist in the map.
	r, has = OptionalIntP("test", "key3", m, 3)
	a.False(has)
	a.Equal(3, r)

	// Test that OptionalIntP panics when the specified key exists in the map but cannot be parsed as a int.
	m["key4"] = "not-a-int-value"
	a.Panics(func() { OptionalIntP("test", "key4", m, 4) })
}

func TestOptionalInt(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 6,
	}

	// Test that OptionalInt returns the value of the specified key when it exists in the map.
	r, has, err := OptionalInt("test", "key1", m, -1)
	a.True(has)
	a.NoError(err)
	a.Equal(6, r)

	// Test that OptionalInt returns the default value when the specified key does not exist in the map.
	r, has, err = OptionalInt("test", "key3", m, -2)
	a.False(has)
	a.NoError(err)
	a.Equal(r, -2)
}

func TestOptionalInt_error(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key4": "not-a-int-value",
	}

	// Test that OptionalInt returns an error when the specified key exists in the map but cannot be parsed as a int.
	_, has, err := OptionalInt("test", "key4", m, -3)
	a.True(has)
	a.Error(err)
}

// TestInt_FloatConversion tests converting float32/float64 to int
func TestInt_FloatConversion(t *testing.T) {
	a := require.New(t)

	// Test float32 to int conversion
	r, err := Int("test", float32(42.7))
	a.NoError(err)
	a.Equal(42, r)

	// Test float64 to int conversion
	r, err = Int("test", float64(99.9))
	a.NoError(err)
	a.Equal(99, r)

	// Test negative float
	r, err = Int("test", float64(-15.3))
	a.NoError(err)
	a.Equal(-15, r)
}

// TestInt_I18nError tests that error messages are localized
func TestInt_I18nError(t *testing.T) {
	a := require.New(t)

	// Test with English
	InitI18n("en")
	_, err := Int("testField", "not-an-int")
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "int")

	// Test with Chinese
	InitI18n("zh")
	_, err = Int("testField", "not-an-int")
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "整数")
}

// TestIntArray tests IntArray function
func TestIntArray(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with []int
	r, err := IntArray("test", []int{1, 2, 3})
	a.NoError(err)
	a.Equal([]int{1, 2, 3}, r)

	// Test with []any containing ints
	r, err = IntArray("test", []any{4, 5, 6})
	a.NoError(err)
	a.Equal([]int{4, 5, 6}, r)

	// Test with []any containing float64
	r, err = IntArray("test", []any{float64(7), float64(8)})
	a.NoError(err)
	a.Equal([]int{7, 8}, r)

	// Test with single int value (converts to array)
	r, err = IntArray("test", 42)
	a.NoError(err)
	a.Equal([]int{42}, r)

	// Test with invalid type
	_, err = IntArray("test", "not-an-array")
	a.Error(err)
	a.Contains(err.Error(), "array")
}

// TestIntArrayP tests IntArrayP panic behavior
func TestIntArrayP(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test successful conversion
	r := IntArrayP("test", []int{10, 20, 30})
	a.Equal([]int{10, 20, 30}, r)

	// Test panic on invalid type
	a.Panics(func() { IntArrayP("test", "invalid") })
}

// TestIntMap tests IntMap function
func TestIntMap(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with map[string]int
	r, err := IntMap("test", map[string]int{"a": 1, "b": 2})
	a.NoError(err)
	a.Equal(1, r["a"])
	a.Equal(2, r["b"])

	// Test with map[string]any
	r, err = IntMap("test", map[string]any{"x": 10, "y": 20})
	a.NoError(err)
	a.Equal(10, r["x"])
	a.Equal(20, r["y"])

	// Test with string "key:value" format
	r, err = IntMap("test", "count:42")
	a.NoError(err)
	a.Equal(42, r["count"])

	// Test with invalid type
	_, err = IntMap("test", 123)
	a.Error(err)
	a.Contains(err.Error(), "map")
}

// TestIntMapP tests IntMapP panic behavior
func TestIntMapP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := IntMapP("test", map[string]int{"a": 100})
	a.Equal(100, r["a"])

	// Test panic on invalid type
	a.Panics(func() { IntMapP("test", []int{1, 2, 3}) })
}

