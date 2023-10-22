package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that BoolP returns the value of the specified key when it exists in the map.
	r := RequiredBoolP("task", "key1", m)
	a.True(r)

	r = RequiredBoolP("task", "key2", m)
	a.False(r)

	// Test that BoolP returns an error when the specified key does not exist in the map
	a.Panics(func() { RequiredBoolP("task", "not-existed-key", m) })
}

func TestBoolP_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	// Test that failed to parse bool value
	a.Panics(func() { RequiredBoolP("task", "key1", m) })
}

func TestBool(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that BoolP returns the value of the specified key when it exists in the map.
	r, err := RequiredBool("task", "key1", m)
	a.NoError(err)
	a.True(r)

	r, err = RequiredBool("task", "key2", m)
	a.NoError(err)
	a.False(r)

	// Test that BoolP returns an error when the specified key does not exist in the map
	_, err = RequiredBool("task", "not-existed-key", m)
	a.Error(err)
}

func TestBool_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	// Test that failed to parse bool value
	_, err := RequiredBool("task", "key1", m)
	a.Error(err)
}

func TestOptionalBoolP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	// Test that OptionalBoolP panics when the specified key exists in the map but cannot be parsed as a bool.
	a.Panics(func() { OptionalBoolP("task", "key1", m, false) })
}

func TestOptionalBoolP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that OptionalBoolP returns the value of the specified key when it exists in the map.
	r, has := OptionalBoolP("test", "key1", m, false)
	a.True(has)
	a.True(r)

	r, has = OptionalBoolP("test", "key2", m, true)
	a.True(has)
	a.False(r)

	// Test that OptionalBoolP returns the default value when the specified key does not exist in the map.
	r, has = OptionalBoolP("test", "key3", m, true)
	a.False(has)
	a.True(r)

	// Test that OptionalBoolP panics when the specified key exists in the map but cannot be parsed as a bool.
	m["key4"] = "not-a-bool-value"
	a.Panics(func() { OptionalBoolP("test", "key4", m, false) })
}

func TestOptionalBool(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that OptionalBool returns the value of the specified key when it exists in the map.
	r, has, err := OptionalBool("test", "key1", m, false)
	a.True(has)
	a.NoError(err)
	a.True(r)

	r, has, err = OptionalBool("test", "key2", m, true)
	a.True(has)
	a.NoError(err)
	a.False(r)

	// Test that OptionalBool returns the default value when the specified key does not exist in the map.
	r, has, err = OptionalBool("test", "key3", m, true)
	a.False(has)
	a.NoError(err)
	a.True(r)
}

func TestOptionalBool_error(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key4": "not-a-bool-value",
	}

	// Test that OptionalBool returns an error when the specified key exists in the map but cannot be parsed as a bool.
	_, has, err := OptionalBool("test", "key4", m, false)
	a.True(has)
	a.Error(err)
}
