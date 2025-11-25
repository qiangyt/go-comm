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

func TestBoolP_direct(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	result := BoolP("test", true)
	a.True(result)

	result = BoolP("test", false)
	a.False(result)

	// Test panic on invalid type
	a.Panics(func() { BoolP("test", "not a bool") })
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

// TestBool_StringConversion tests converting string to bool
func TestBool_StringConversion(t *testing.T) {
	a := require.New(t)

	// Test "true" string
	r, err := Bool("test", "true")
	a.NoError(err)
	a.True(r)

	// Test "TRUE" string (case insensitive)
	r, err = Bool("test", "TRUE")
	a.NoError(err)
	a.True(r)

	// Test " True " string (with spaces)
	r, err = Bool("test", " True ")
	a.NoError(err)
	a.True(r)

	// Test "false" string
	r, err = Bool("test", "false")
	a.NoError(err)
	a.False(r)

	// Test "FALSE" string (case insensitive)
	r, err = Bool("test", "FALSE")
	a.NoError(err)
	a.False(r)

	// Test invalid string
	_, err = Bool("test", "yes")
	a.Error(err)
}

// TestBool_I18nError tests that error messages are localized
func TestBool_I18nError(t *testing.T) {
	a := require.New(t)

	// Test with English
	InitI18n("en")
	_, err := Bool("testField", 123)
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "bool")

	// Test with Chinese
	InitI18n("zh")
	_, err = Bool("testField", 123)
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "布尔")
}

// TestBoolArray tests BoolArray function
func TestBoolArray(t *testing.T) {
	a := require.New(t)

	// Test with []bool
	r, err := BoolArray("test", []bool{true, false, true})
	a.NoError(err)
	a.Equal([]bool{true, false, true}, r)

	// Test with []any containing bools
	r, err = BoolArray("test", []any{true, false})
	a.NoError(err)
	a.Equal([]bool{true, false}, r)

	// Test with []any containing string bools
	r, err = BoolArray("test", []any{"true", "false"})
	a.NoError(err)
	a.Equal([]bool{true, false}, r)

	// Test with single bool value (converts to array)
	r, err = BoolArray("test", true)
	a.NoError(err)
	a.Equal([]bool{true}, r)

	// Test with invalid type
	_, err = BoolArray("test", "not-an-array")
	a.Error(err)
	a.Contains(err.Error(), "array")
}

// TestBoolArrayP tests BoolArrayP panic behavior
func TestBoolArrayP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := BoolArrayP("test", []bool{true, false, true})
	a.Equal([]bool{true, false, true}, r)

	// Test panic on invalid type
	a.Panics(func() { BoolArrayP("test", "invalid") })
}

// TestBoolMap tests BoolMap function
func TestBoolMap(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with map[string]bool
	r, err := BoolMap("test", map[string]bool{"enabled": true, "disabled": false})
	a.NoError(err)
	a.True(r["enabled"])
	a.False(r["disabled"])

	// Test with map[string]any
	r, err = BoolMap("test", map[string]any{"active": true, "inactive": false})
	a.NoError(err)
	a.True(r["active"])
	a.False(r["inactive"])

	// Test with string "key:value" format
	r, err = BoolMap("test", "visible:true")
	a.NoError(err)
	a.True(r["visible"])

	// Test with string "key:false" format
	r, err = BoolMap("test", "hidden:false")
	a.NoError(err)
	a.False(r["hidden"])

	// Test with invalid type
	_, err = BoolMap("test", 123)
	a.Error(err)
	a.Contains(err.Error(), "map")
}

// TestBoolMapP tests BoolMapP panic behavior
func TestBoolMapP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := BoolMapP("test", map[string]bool{"ok": true})
	a.True(r["ok"])

	// Test panic on invalid type
	a.Panics(func() { BoolMapP("test", []bool{true, false}) })
}

