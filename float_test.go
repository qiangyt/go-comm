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

// TestFloat_IntConversion tests converting int/int32/int64 to float64
func TestFloat_IntConversion(t *testing.T) {
	a := require.New(t)

	// Test int to float64 conversion
	r, err := Float("test", int(42))
	a.NoError(err)
	a.Equal(42.0, r)

	// Test int32 to float64 conversion
	r, err = Float("test", int32(99))
	a.NoError(err)
	a.Equal(99.0, r)

	// Test int64 to float64 conversion
	r, err = Float("test", int64(-15))
	a.NoError(err)
	a.Equal(-15.0, r)
}

// TestFloat_I18nError tests that error messages are localized
func TestFloat_I18nError(t *testing.T) {
	a := require.New(t)

	// Test with English
	InitI18n("en")
	_, err := Float("testField", "not-a-float")
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "float")

	// Test with Chinese
	InitI18n("zh")
	_, err = Float("testField", "not-a-float")
	a.Error(err)
	a.Contains(err.Error(), "testField")
	a.Contains(err.Error(), "浮点数")
}

// TestFloatArray tests FloatArray function
func TestFloatArray(t *testing.T) {
	a := require.New(t)

	// Test with []float64
	r, err := FloatArray("test", []float64{1.1, 2.2, 3.3})
	a.NoError(err)
	a.Equal([]float64{1.1, 2.2, 3.3}, r)

	// Test with []any containing float64
	r, err = FloatArray("test", []any{4.4, 5.5, 6.6})
	a.NoError(err)
	a.Equal([]float64{4.4, 5.5, 6.6}, r)

	// Test with []any containing ints
	r, err = FloatArray("test", []any{int(7), int32(8), int64(9)})
	a.NoError(err)
	a.Equal([]float64{7.0, 8.0, 9.0}, r)

	// Test with single float64 value (converts to array)
	r, err = FloatArray("test", 42.5)
	a.NoError(err)
	a.Equal([]float64{42.5}, r)

	// Test with invalid type
	_, err = FloatArray("test", "not-an-array")
	a.Error(err)
	a.Contains(err.Error(), "array")
}

// TestFloatArrayP tests FloatArrayP panic behavior
func TestFloatArrayP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := FloatArrayP("test", []float64{10.5, 20.5, 30.5})
	a.Equal([]float64{10.5, 20.5, 30.5}, r)

	// Test panic on invalid type
	a.Panics(func() { FloatArrayP("test", "invalid") })
}

// TestFloatMap tests FloatMap function
func TestFloatMap(t *testing.T) {
	a := require.New(t)
	InitI18n("en")

	// Test with map[string]float64
	r, err := FloatMap("test", map[string]float64{"a": 1.1, "b": 2.2})
	a.NoError(err)
	a.Equal(1.1, r["a"])
	a.Equal(2.2, r["b"])

	// Test with map[string]any
	r, err = FloatMap("test", map[string]any{"x": 10.5, "y": 20.5})
	a.NoError(err)
	a.Equal(10.5, r["x"])
	a.Equal(20.5, r["y"])

	// Test with string "key:value" format
	r, err = FloatMap("test", "price:99.99")
	a.NoError(err)
	a.Equal(99.99, r["price"])

	// Test with invalid type
	_, err = FloatMap("test", 123)
	a.Error(err)
	a.Contains(err.Error(), "map")
}

// TestFloatMapP tests FloatMapP panic behavior
func TestFloatMapP(t *testing.T) {
	a := require.New(t)

	// Test successful conversion
	r := FloatMapP("test", map[string]float64{"a": 100.5})
	a.Equal(100.5, r["a"])

	// Test panic on invalid type
	a.Panics(func() { FloatMapP("test", []float64{1.1, 2.2, 3.3}) })
}

