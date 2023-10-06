package comm

import (
	"fmt"
	"reflect"
)

func RequiredStringP(hint string, key string, m map[string]any) string {
	r, err := RequiredString(hint, key, m)
	if err != nil {
		panic(err)
	}
	return r
}

func RequiredString(hint string, key string, m map[string]any) (string, error) {
	v, has := m[key]
	if !has {
		return "", fmt.Errorf("%s.%s is required", hint, key)
	}

	return String(hint+"."+key, v)
}

func OptionalStringP(hint string, key string, m map[string]any, devault string) string {
	r, err := OptionalString(hint, key, m, devault)
	if err != nil {
		panic(err)
	}
	return r
}

func OptionalString(hint string, key string, m map[string]any, devault string) (string, error) {
	v, has := m[key]
	if !has {
		return devault, nil
	}

	return String(hint+"."+key, v)
}

func StringP(hint string, v any) string {
	r, err := String(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

func String(hint string, v any) (string, error) {
	r, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}

func StringArrayValueP(hint string, key string, m map[string]any) []string {
	r, err := StringArrayValue(hint, key, m)
	if err != nil {
		panic(err)
	}
	return r
}

func StringArrayValue(hint string, key string, m map[string]any) ([]string, error) {
	v, has := m[key]
	if !has {
		return nil, fmt.Errorf("%s.%s is required", hint, key)
	}

	return StringArray(hint+"."+key, v)
}

func OptionalStringArrayValueP(hint string, key string, m map[string]any, devault []string) []string {
	r, err := OptionalStringArrayValue(hint, key, m, devault)
	if err != nil {
		panic(err)
	}
	return r
}

func OptionalStringArrayValue(hint string, key string, m map[string]any, devault []string) ([]string, error) {
	v, has := m[key]
	if !has {
		return devault, nil
	}

	return StringArray(hint+"."+key, v)
}

func StringArrayP(hint string, v any) []string {
	r, err := StringArray(hint, v)
	if err != nil {
		panic(err)
	}
	return r
}

func StringArray(hint string, v any) ([]string, error) {
	r, ok := v.([]string)
	if !ok {
		return nil, fmt.Errorf("%s must be a string array, but now it is a %v(%v)", hint, reflect.TypeOf(v), v)
	}
	return r, nil
}
