package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequiredStringP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 123,
	}

	// Test that StringP panics when the specified key exists in the map but cannot be parsed as a string.
	a.Panics(func() { RequiredStringP("task", "key1", m) })
}

func TestRequiredStringP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	// Test that StringP returns the value of the specified key when it exists in the map.
	r := RequiredStringP("test", "key1", m)
	a.Equal("value1", r)

	// Test that StringP panics when the specified key does not exist in the map.
	a.Panics(func() { RequiredStringP("test", "key3", m) })
}

func TestOptionalStringP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 123,
	}

	// Test that OptionalStringP panics when the specified key exists in the map but cannot be parsed as a string.
	a.Panics(func() { OptionalStringP("task", "key1", m, "default") })
}

func TestOptionalStringP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	// Test that OptionalStringP returns the value of the specified key when it exists in the map.
	r, has := OptionalStringP("test", "key1", m, "default")
	a.True(has)
	a.Equal("value1", r)

	// Test that OptionalStringP returns the default value when the specified key does not exist in the map.
	r, has = OptionalStringP("test", "key3", m, "default")
	a.False(has)
	a.Equal("default", r)

	// Test that OptionalStringP panics when the specified key exists in the map but cannot be parsed as a string.
	m["key4"] = 123
	a.Panics(func() { OptionalStringP("test", "key4", m, "default") })
}

func TestStringArrayValueP_panic(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": 123, // Use an invalid type (int cannot be converted to string array)
	}

	// Test that StringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	a.Panics(func() { StringArrayValueP("task", "key1", m) })
}

func TestStringArrayValueP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3", "value4"},
	}

	// Test that StringArrayValueP returns the value of the specified key when it exists in the map.
	r := StringArrayValueP("test", "key1", m)
	a.Equal([]string{"value1", "value2"}, r)

	// Test that StringArrayValueP panics when the specified key does not exist in the map.
	a.Panics(func() { StringArrayValueP("test", "key3", m) })
}

func TestOptionalStringArrayValueP_panic(t *testing.T) {
	a := require.New(t)
	m := map[string]any{
		"key1": 123, // Use an invalid type (int cannot be converted to string array)
	}

	// Test that OptionalStringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	a.Panics(func() { OptionalStringArrayValueP("task", "key1", m, []string{"default"}) })
}

func TestOptionalStringArrayValueP(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3", "value4"},
	}

	// Test that OptionalStringArrayValueP returns the value of the specified key when it exists in the map.
	r, has := OptionalStringArrayValueP("test", "key1", m, []string{"default"})
	a.True(has)
	a.Equal([]string{"value1", "value2"}, r)
	//if !reflect.DeepEqual(r, []string{"value1", "value2"}) {
	//	t.Errorf("expected ['value1', 'value2'], but got '%v'", r)
	//}

	// Test that OptionalStringArrayValueP returns the default value when the specified key does not exist in the map.
	r, has = OptionalStringArrayValueP("test", "key3", m, []string{"default"})
	a.False(has)
	a.Equal([]string{"default"}, r)

	// Test that OptionalStringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	m["key4"] = 456 // Use an invalid type
	a.Panics(func() { OptionalStringArrayValueP("test", "key4", m, []string{"default"}) })
}

// TestString_I18nError tests that error messages are localized
func TestString_I18nError(t *testing.T) {
	a := require.New(t)

	// Test with English
	InitI18n("en")
	_, err := String("testField", 123)
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "string")

	// Test with Chinese
	InitI18n("zh")
	_, err = String("testField", 123)
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "字符串")
}

// TestStringArray tests StringArray function with various inputs
func TestStringArray(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with []string
	r, err := StringArray("test", []string{"a", "b", "c"})
	a.NoError(err)
	a.Equal([]string{"a", "b", "c"}, r)

	// Test with []any containing strings
	r, err = StringArray("test", []any{"x", "y", "z"})
	a.NoError(err)
	a.Equal([]string{"x", "y", "z"}, r)

	// Test with single string value (converts to array)
	r, err = StringArray("test", "single")
	a.NoError(err)
	a.Equal([]string{"single"}, r)

	// Test with invalid element type in array
	_, err = StringArray("test", []any{"valid", 123})
	a.Error(err)

	// Test with completely invalid type
	_, err = StringArray("test", 123)
	a.Error(err)
	a.Contains(err.Error(), "array")
}

// TestStringMap tests StringMap function with various inputs
func TestStringMap(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with map[string]string
	r, err := StringMap("test", map[string]string{"k1": "v1", "k2": "v2"})
	a.NoError(err)
	a.Equal("v1", r["k1"])
	a.Equal("v2", r["k2"])

	// Test with map[string]any
	r, err = StringMap("test", map[string]any{"key": "value"})
	a.NoError(err)
	a.Equal("value", r["key"])

	// Test with string "key:value" format
	r, err = StringMap("test", "name:John")
	a.NoError(err)
	a.Equal("John", r["name"])

	// Test with string "key: value" format (with spaces)
	r, err = StringMap("test", "city: Beijing")
	a.NoError(err)
	a.Equal("Beijing", r["city"])

	// Test with invalid type
	_, err = StringMap("test", 123)
	a.Error(err)
	a.Contains(err.Error(), "map")
}

// TestStringP tests StringP function
func TestStringP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := StringP("test", "hello")
	a.Equal("hello", r)

	// Test panic on invalid type
	a.Panics(func() { StringP("test", 123) })
}

// TestStringArrayP tests StringArrayP function
func TestStringArrayP(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test successful conversion
	r := StringArrayP("test", []string{"a", "b"})
	a.Equal([]string{"a", "b"}, r)

	// Test panic on invalid type
	a.Panics(func() { StringArrayP("test", 123) })
}

// TestStringMapP tests StringMapP function
func TestStringMapP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := StringMapP("test", map[string]string{"key": "value"})
	a.Equal("value", r["key"])

	// Test panic on invalid type
	a.Panics(func() { StringMapP("test", []string{"a", "b"}) })
}

// TestRequiredString_I18nError tests that error messages are localized
func TestRequiredString_I18nError(t *testing.T) {
	a := require.New(t)

	m := map[string]any{}

	// Test with English
	InitI18n("en")
	_, err := RequiredString("config", "name", m)
	a.Error(err)
	a.Contains(err.Error(), "config.name")
	a.Contains(err.Error(), "required")

	// Test with Chinese
	InitI18n("zh")
	_, err = RequiredString("config", "name", m)
	a.Error(err)
	a.Contains(err.Error(), "config.name")
	a.Contains(err.Error(), "必需")
}

