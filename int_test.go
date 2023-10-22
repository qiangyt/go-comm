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
