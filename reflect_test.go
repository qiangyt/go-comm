package comm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsEmptyReflectValue_string(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf("")))
	a.False(IsEmptyReflectValue(reflect.ValueOf("hello")))
}

func TestIsEmptyReflectValue_int(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf(0)))
	a.False(IsEmptyReflectValue(reflect.ValueOf(42)))
	a.False(IsEmptyReflectValue(reflect.ValueOf(-1)))
}

func TestIsEmptyReflectValue_float(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf(0.0)))
	a.False(IsEmptyReflectValue(reflect.ValueOf(3.14)))
}

func TestIsEmptyReflectValue_bool(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf(false)))
	a.False(IsEmptyReflectValue(reflect.ValueOf(true)))
}

func TestIsEmptyReflectValue_slice(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf([]int{})))
	a.False(IsEmptyReflectValue(reflect.ValueOf([]int{1, 2, 3})))
}

func TestIsEmptyReflectValue_map(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf(map[string]int{})))
	a.False(IsEmptyReflectValue(reflect.ValueOf(map[string]int{"key": 1})))
}

func TestIsEmptyReflectValue_array(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf([0]int{})))
	a.False(IsEmptyReflectValue(reflect.ValueOf([3]int{1, 2, 3})))
}

func TestIsEmptyReflectValue_pointer(t *testing.T) {
	a := require.New(t)

	var nilPtr *int = nil
	a.True(IsEmptyReflectValue(reflect.ValueOf(nilPtr)))

	value := 42
	a.False(IsEmptyReflectValue(reflect.ValueOf(&value)))
}

func TestIsEmptyReflectValue_uint(t *testing.T) {
	a := require.New(t)

	a.True(IsEmptyReflectValue(reflect.ValueOf(uint(0))))
	a.False(IsEmptyReflectValue(reflect.ValueOf(uint(100))))
}

func TestIsEmptyValue(t *testing.T) {
	a := require.New(t)

	// Test various empty values
	a.True(IsEmptyValue(""))
	a.True(IsEmptyValue(0))
	a.True(IsEmptyValue(false))
	a.True(IsEmptyValue([]int{}))
	a.True(IsEmptyValue(map[string]int{}))

	// Test non-empty values
	a.False(IsEmptyValue("hello"))
	a.False(IsEmptyValue(42))
	a.False(IsEmptyValue(true))
	a.False(IsEmptyValue([]int{1}))
	a.False(IsEmptyValue(map[string]int{"key": 1}))
}

func TestIsPrimitiveReflectValue_primitives(t *testing.T) {
	a := require.New(t)

	// Test primitive types
	a.True(IsPrimitiveReflectValue(reflect.ValueOf("string")))
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(42)))
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(3.14)))
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(true)))
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(uint(10))))
}

func TestIsPrimitiveReflectValue_nonPrimitives(t *testing.T) {
	a := require.New(t)

	// Test non-primitive types
	a.False(IsPrimitiveReflectValue(reflect.ValueOf([]int{1, 2, 3})))
	a.False(IsPrimitiveReflectValue(reflect.ValueOf(map[string]int{"key": 1})))

	type MyStruct struct{ Field int }
	a.False(IsPrimitiveReflectValue(reflect.ValueOf(MyStruct{Field: 1})))
}

func TestIsPrimitiveReflectValue_pointer(t *testing.T) {
	a := require.New(t)

	// Pointer to primitive is primitive
	value := 42
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(&value)))

	// Nil pointer is primitive
	var nilPtr *int = nil
	a.True(IsPrimitiveReflectValue(reflect.ValueOf(nilPtr)))

	// Pointer to non-primitive is non-primitive
	slice := []int{1, 2, 3}
	a.False(IsPrimitiveReflectValue(reflect.ValueOf(&slice)))
}

func TestIsPrimitiveValue(t *testing.T) {
	a := require.New(t)

	// Test primitive values
	a.True(IsPrimitiveValue("test"))
	a.True(IsPrimitiveValue(123))
	a.True(IsPrimitiveValue(true))
	a.True(IsPrimitiveValue(3.14))

	// Test non-primitive values
	a.False(IsPrimitiveValue([]int{1, 2, 3}))
	a.False(IsPrimitiveValue(map[string]int{"key": 1}))

	type TestStruct struct{ Value int }
	a.False(IsPrimitiveValue(TestStruct{Value: 1}))
}
