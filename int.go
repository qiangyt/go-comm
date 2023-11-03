package comm

import (
	"fmt"
	"reflect"
)

// RequiredIntP returns the int value of the key in the map. If either parsing error or the key is not
// found, raise a panic.
func RequiredIntP(hint string, key string, m map[string]any) int {
	r, err := RequiredInt(hint, key, m)
	if err != nil {
		panic(err)
	}
	return r
}

// RequiredInt returns the int value of the key in the map. If either parsing error or the key is not
// found, an error is returned.
func RequiredInt(hint string, key string, m map[string]any) (int, error) {
	v, has := m[key]
	if !has {
		return 0, fmt.Errorf("%s.%s is required", hint, key)
	}

	return Int(hint+"."+key, v)
}

// OptionalIntP returns the int value of the key in the map. If parsing error occurred,
// raise a panic. If the key is not found, return the default value.
func OptionalIntP(hint string, key string, m map[string]any, devault int) (result int, has bool) {
	var err error
	result, has, err = OptionalInt(hint, key, m, devault)
	if err != nil {
		panic(err)
	}
	return
}

// OptionalInt returns the int value of the key in the map. If parsing error occrred,
// returns the error. If the key is not found, return the default value.
func OptionalInt(hint string, key string, m map[string]any, devault int) (result int, has bool, err error) {
	var v any

	v, has = m[key]
	if !has {
		result = devault
		return
	}

	result, err = Int(hint+"."+key, v)
	return
}

// Cast the value to int. If parsing error occurred, raise a panic.
func IntP(hint string, v any) int {
	r, err := Int(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

// Cast the value to int. If parsing error occurred, returns the error.
func Int(hint string, v any) (int, error) {
	r, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("%s must be a int, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}

func IntArrayP(hint string, v any) []int {
	r, err := IntArray(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

func IntArray(hint string, v any) ([]int, error) {
	r, ok := v.([]int)
	if !ok {
		if r0, ok0 := v.(int); ok0 {
			return []int{r0}, nil
		}
		return nil, fmt.Errorf("%s must be a int array, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}
