package comm

import (
	"fmt"
	"reflect"
)

// RequiredFloatP returns the float64 value of the key in the map. If either parsing error or the key is not
// found, raise a panic.
func RequiredFloatP(hint string, key string, m map[string]any) float64 {
	r, err := RequiredFloat(hint, key, m)
	if err != nil {
		panic(err)
	}
	return r
}

// RequiredFloat returns the float64 value of the key in the map. If either parsing error or the key is not
// found, an error is returned.
func RequiredFloat(hint string, key string, m map[string]any) (float64, error) {
	v, has := m[key]
	if !has {
		return 0, fmt.Errorf("%s.%s is required", hint, key)
	}

	return Float(hint+"."+key, v)
}

// OptionalFloatP returns the float64 value of the key in the map. If parsing error occurred,
// raise a panic. If the key is not found, return the default value.
func OptionalFloatP(hint string, key string, m map[string]any, devault float64) (result float64, has bool) {
	var err error
	result, has, err = OptionalFloat(hint, key, m, devault)
	if err != nil {
		panic(err)
	}
	return
}

// OptionalFloat returns the float64 value of the key in the map. If parsing error occrred,
// returns the error. If the key is not found, return the default value.
func OptionalFloat(hint string, key string, m map[string]any, devault float64) (result float64, has bool, err error) {
	var v any

	v, has = m[key]
	if !has {
		result = devault
		return
	}

	result, err = Float(hint+"."+key, v)
	return
}

// Cast the value to float64. If parsing error occurred, raise a panic.
func FloatP(hint string, v any) float64 {
	r, err := Float(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

// Cast the value to float64. If parsing error occurred, returns the error.
func Float(hint string, v any) (float64, error) {
	r, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("%s must be a float64, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}

func FloatArrayP(hint string, v any) []float64 {
	r, err := FloatArray(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

func FloatArray(hint string, v any) ([]float64, error) {
	r, ok := v.([]float64)
	if !ok {
		return nil, fmt.Errorf("%s must be a float64 array, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}
