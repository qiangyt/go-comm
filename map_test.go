package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequiredMapP(t *testing.T) {
	a := require.New(t)

	// Test that RequiredMapP returns the expected result when the specified key exists in the map and its value is a map[string]any.
	m := map[string]any{
		"key1": map[string]any{"key2": "value2"},
	}
	r := RequiredMapP("test", "key1", m)
	expected := map[string]any{"key2": "value2"}
	a.Equal(expected, r)

	// Test that RequiredMapP panics when the specified key does not exist in the map.
	f := func() { RequiredMapP("test", "key2", m) }
	a.Panics(f)
}

func TestRequiredMap(t *testing.T) {
	a := require.New(t)

	// Test that RequiredMap returns the expected result when the specified key exists in the map and its value is a map[string]any.
	m := map[string]any{
		"key1": map[string]any{"key2": "value2"},
	}
	r, err := RequiredMap("test", "key1", m)
	a.NoError(err)
	expected := map[string]any{"key2": "value2"}
	a.Equal(expected, r)

	// Test that RequiredMap returns an error when the specified key does not exist in the map.
	_, err = RequiredMap("test", "key2", m)
	expectedError := "test.key2 is required"
	a.EqualError(err, expectedError)

	// Test that RequiredMap returns an error when the value of the specified key is not a map[string]any.
	m["key3"] = "not-a-map"
	_, err = RequiredMap("test", "key3", m)
	expectedError = "test.key3 must be a map[string]any, but now it is a string(not-a-map)"
	a.EqualError(err, expectedError)
}

func TestOptionalMapP(t *testing.T) {
	a := require.New(t)

	// Test that OptionalMap returns the expected result when the specified key exists in the map and its value is a map[string]any.
	m := map[string]any{
		"key1": map[string]any{"key2": "value2"},
	}
	r := OptionalMapP("test", "key1", m, nil)
	expected := map[string]any{"key2": "value2"}
	a.Equal(expected, r)

	defaultValue := map[string]any{"default": "value"}
	// Test that OptionalMap returns default value when the specified key does not exist in the map.
	r = OptionalMapP("test", "key2", m, defaultValue)
	a.Equal(defaultValue, r)

	// Test that OptionalMapP panices when the value of the specified key is not a map[string]any.
	m["key3"] = "not-a-map"
	a.Panics(func() { OptionalMapP("test", "key3", m, nil) })
}

func TestOptionalMap(t *testing.T) {
	a := require.New(t)

	// Test that OptionalMap returns the expected result when the specified key exists in the map and its value is a map[string]any.
	m := map[string]any{
		"key1": map[string]any{"key2": "value2"},
	}
	r, err := OptionalMap("test", "key1", m, nil)
	a.NoError(err)
	expected := map[string]any{"key2": "value2"}
	a.Equal(expected, r)

	defaultValue := map[string]any{"default": "value"}
	// Test that OptionalMap returns default value when the specified key does not exist in the map.
	r, err = OptionalMap("test", "key2", m, defaultValue)
	a.NoError(err)
	a.Equal(defaultValue, r)

	// Test that OptionalMap returns error when the value of the specified key is not a map[string]any.
	m["key3"] = "not-a-map"
	r, err = OptionalMap("test", "key3", m, defaultValue)
	a.Error(err)
	a.Nil(r)
}

func TestMapP(t *testing.T) {
	a := require.New(t)

	// Test that Map returns the expected result when the input value is a map[string]any.
	input := map[string]any{"key1": "value1", "key2": "value2"}
	r, err := Map("test", input)
	a.NoError(err)
	expected := map[string]interface{}{"key1": "value1", "key2": "value2"}
	a.Equal(expected, r)

	// Test that Map panicesr when the input value is not a map[string]any.
	a.Panics(func() { MapP("test", 123) }, "test must be a map[string]any, but now it is a int(123)")
}

func TestMap(t *testing.T) {
	a := require.New(t)

	// Test that Map returns the expected result when the input value is a map[string]any.
	input := map[string]any{"key1": "value1", "key2": "value2"}
	r, err := Map("test", input)
	a.NoError(err)
	expected := map[string]interface{}{"key1": "value1", "key2": "value2"}
	a.Equal(expected, r)

	// Test that Map returns an error when the input value is not a map[string]any.
	_, err = Map("test", 123)
	a.EqualError(err, "test must be a map[string]any, but now it is a int(123)")
}

func TestRequiredMapArrayP(t *testing.T) {
	a := require.New(t)

	// Test that RequiredMapArrayP returns the expected result when the specified key exists in the map and its value is a []map[string]any.
	m := map[string]any{
		"key1": []map[string]any{{"key2": "value2"}},
	}
	r := RequiredMapArrayP("test", "key1", m)
	expected := []map[string]any{{"key2": "value2"}}
	a.Equal(expected, r)

	// Test that RequiredMapArrayP panics when the specified key does not exist in the map.
	a.Panics(func() { RequiredMapArrayP("test", "key2", m) })

	// Test that RequiredMapArrayP panics when the value of the specified key is not a []map[string]any.
	m["key3"] = "not-an-array"
	a.Panics(func() { RequiredMapArrayP("test", "key3", m) })
}

func TestRequiredMapArray(t *testing.T) {
	a := require.New(t)

	// Test that RequiredMapArray returns the expected result when the specified key exists in the map and its value is a []map[string]any.
	m := map[string]any{
		"key1": []map[string]any{{"key2": "value2"}},
	}
	r, err := RequiredMapArray("test", "key1", m)
	a.NoError(err)
	expected := []map[string]any{{"key2": "value2"}}
	a.Equal(expected, r)

	// Test that RequiredMapArray returns error when the specified key does not exist in the map.
	_, err = RequiredMapArray("test", "key2", m)
	a.Error(err)

	// Test that OptionalMapArray returns error when the value of the specified key is not a []map[string]any.
	m["key3"] = "not-an-array"
	_, err = RequiredMapArray("test", "key3", m)
	a.Error(err)
}

func TestOptionalMapArrayP(t *testing.T) {
	a := require.New(t)

	// Test that OptionalMapArrayP returns the expected result when the specified key exists in the map and its value is a []map[string]any.
	m := map[string]any{
		"key1": []map[string]any{{"key2": "value2"}},
	}
	r := OptionalMapArrayP("test", "key1", m, nil)
	expected := []map[string]any{{"key2": "value2"}}
	a.Equal(expected, r)

	// Test that OptionalMapArrayP returns the default result when the specified key does not exist in the map.
	defaultResult := []map[string]any{{"default": "value"}}
	r = OptionalMapArrayP("test", "key2", m, defaultResult)
	a.Equal(defaultResult, r)

	// Test that OptionalMapArrayP panics when the value of the specified key is not a []map[string]any.
	m["key3"] = "not-an-array"
	a.Panics(func() { OptionalMapArrayP("test", "key3", m, defaultResult) })
}

func TestOptionalMapArray(t *testing.T) {
	a := require.New(t)

	// Test that OptionalMapArray returns the expected result when the specified key exists in the map and its value is a []map[string]any.
	m := map[string]any{
		"key1": []map[string]any{{"key2": "value2"}},
	}
	r, err := OptionalMapArray("test", "key1", m, nil)
	a.NoError(err)
	expected := []map[string]any{{"key2": "value2"}}
	a.Equal(expected, r)

	// Test that OptionalMapArray returns the default result when the specified key does not exist in the map.
	defaultResult := []map[string]any{{"default": "value"}}
	r, err = OptionalMapArray("test", "key2", m, defaultResult)
	a.NoError(err)
	a.Equal(defaultResult, r)

	// Test that OptionalMapArray returns error when the value of the specified key is not a []map[string]any.
	m["key3"] = "not-an-array"
	_, err = OptionalMapArray("test", "key3", m, defaultResult)
	a.Error(err)
}

func TestMapArrayP(t *testing.T) {
	a := require.New(t)

	// Test that MapArrayP returns the expected result when the input value is a []map[string]any.
	input := []map[string]any{{"key1": "value1"}, {"key2": "value2"}}
	r := MapArrayP("test", input)
	a.Equal(input, r)

	// Test that MapArrayP panics when the input value is not a []map[string]any.
	a.Panics(func() { MapArrayP("test", 123) })
}

func TestMapArray(t *testing.T) {
	a := require.New(t)

	// Test that MapArray returns the expected result when the input value is a []map[string]any.
	input := []map[string]any{{"key1": "value1"}, {"key2": "value2"}}
	r, err := MapArray("test", input)
	a.NoError(err)
	a.Equal(input, r)

	// Test that MapArray returns an error when the input value is not a []map[string]any.
	_, err = MapArray("test", 123)
	require.Error(t, err)
	expected := "test must be a map[string]any array, but now it is a int(123)"
	a.EqualError(err, expected)
}
