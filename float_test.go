package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloatP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 7.9,
	}

	// Test that FloatP panics when the specified key does not exist in the map.
	// Test that FloatP returns the value of the specified key when it exists in the map.
	r := RequiredFloatP("task", "key1", m)
	a.Equal(7.9, r)

	// Test that FloatP returns an error when the specified key does not exist in the map
	a.Panics(func() { RequiredFloatP("task", "not-existed-key", m) })
}

func TestFloatP_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-float64-value",
	}

	// Test that failed to parse float64 value
	a.Panics(func() { RequiredFloatP("task", "key1", m) })
}

func TestFloat(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 5.9,
	}

	// Test that FloatP returns the value of the specified key when it exists in the map.
	r, err := RequiredFloat("task", "key1", m)
	a.NoError(err)
	a.Equal(5.9, r)

	// Test that FloatP returns an error when the specified key does not exist in the map
	_, err = RequiredFloat("task", "not-existed-key", m)
	a.Error(err)
}

func TestFloat_parsingError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-float64-value",
	}

	// Test that failed to parse float64 value
	_, err := RequiredFloat("task", "key1", m)
	a.Error(err)
}

func TestOptionalFloatP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "not-a-float64-value",
	}

	// Test that OptionalFloatP panics when the specified key exists in the map but cannot be parsed as a float64.
	a.Panics(func() { OptionalFloatP("task", "key1", m, 9.3) })
}

func TestOptionalFloatP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 10.2,
	}

	// Test that OptionalFloatP returns the value of the specified key when it exists in the map.
	r, has := OptionalFloatP("test", "key1", m, 2.4)
	a.True(has)
	a.Equal(10.2, r)

	// Test that OptionalFloatP returns the default value when the specified key does not exist in the map.
	r, has = OptionalFloatP("test", "key3", m, 3.4)
	a.False(has)
	a.Equal(3.4, r)

	// Test that OptionalFloatP panics when the specified key exists in the map but cannot be parsed as a float64.
	m["key4"] = "not-a-float64-value"
	a.Panics(func() { OptionalFloatP("test", "key4", m, 4.5) })
}

func TestOptionalFloat(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 6.6,
	}

	// Test that OptionalFloat returns the value of the specified key when it exists in the map.
	r, has, err := OptionalFloat("test", "key1", m, -1.7)
	a.True(has)
	a.NoError(err)
	a.Equal(6.6, r)

	// Test that OptionalFloat returns the default value when the specified key does not exist in the map.
	r, has, err = OptionalFloat("test", "key3", m, -2.8)
	a.False(has)
	a.NoError(err)
	a.Equal(r, -2.8)
}

func TestOptionalFloat_error(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key4": "not-a-float64-value",
	}

	// Test that OptionalFloat returns an error when the specified key exists in the map but cannot be parsed as a float64.
	_, has, err := OptionalFloat("test", "key4", m, -3.9)
	a.True(has)
	a.Error(err)
}
