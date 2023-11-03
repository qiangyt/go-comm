package comm

import (
	"fmt"
	"reflect"
)

// RequiredBoolP returns the bool value of the key in the map. If either parsing error or the key is not
// found, raise a panic.
func RequiredBoolP(hint string, key string, m map[string]any) bool {
	r, err := RequiredBool(hint, key, m)
	if err != nil {
		panic(err)
	}
	return r
}

// RequiredBool returns the bool value of the key in the map. If either parsing error or the key is not
// found, an error is returned.
func RequiredBool(hint string, key string, m map[string]any) (bool, error) {
	v, has := m[key]
	if !has {
		return false, fmt.Errorf("%s.%s is required", hint, key)
	}

	return Bool(hint+"."+key, v)
}

// OptionalBoolP returns the bool value of the key in the map. If parsing error occurred,
// raise a panic. If the key is not found, return the default value.
func OptionalBoolP(hint string, key string, m map[string]any, devault bool) (result bool, has bool) {
	var err error
	result, has, err = OptionalBool(hint, key, m, devault)
	if err != nil {
		panic(err)
	}
	return
}

// OptionalBool returns the bool value of the key in the map. If parsing error occrred,
// returns the error. If the key is not found, return the default value.
func OptionalBool(hint string, key string, m map[string]any, devault bool) (result bool, has bool, err error) {
	var v any
	v, has = m[key]
	if !has {
		result = devault
		return
	}

	result, err = Bool(hint+"."+key, v)
	return
}

// Cast the value to bool. If parsing error occurred, raise a panic.
func BoolP(hint string, v any) bool {
	r, err := Bool(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

// Cast the value to bool. If parsing error occurred, returns the error.
func Bool(hint string, v any) (bool, error) {
	r, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%s must be a bool, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}

func BoolArrayP(hint string, v any) []bool {
	r, err := BoolArray(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

func BoolArray(hint string, v any) ([]bool, error) {
	r, ok := v.([]bool)
	if !ok {
		return nil, fmt.Errorf("%s must be a bool array, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}
