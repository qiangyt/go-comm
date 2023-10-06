package comm

import (
	"reflect"
	"testing"
)

func TestRequiredStringP_panic(t *testing.T) {
	m := map[string]any{
		"key1": 123,
	}

	// Test that StringP panics when the specified key exists in the map but cannot be parsed as a string.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a string, but now it is a int(123)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	RequiredStringP("task", "key1", m)
}

func TestRequiredStringP(t *testing.T) {
	m := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	// Test that StringP returns the value of the specified key when it exists in the map.
	r := RequiredStringP("test", "key1", m)

	if r != "value1" {
		t.Errorf("expected 'value1', but got '%s'", r)
	}

	r = RequiredStringP("test", "key2", m)

	if r != "value2" {
		t.Errorf("expected 'value2', but got '%s'", r)
	}

	// Test that StringP panics when the specified key does not exist in the map.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key3 is required"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	RequiredStringP("test", "key3", m)
}

func TestOptionalStringP_panic(t *testing.T) {
	m := map[string]any{
		"key1": 123,
	}

	// Test that OptionalStringP panics when the specified key exists in the map but cannot be parsed as a string.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a string, but now it is a int(123)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalStringP("task", "key1", m, "default")
}

func TestOptionalStringP(t *testing.T) {
	m := map[string]any{
		"key1": "value1",
		"key2": "value2",
	}

	// Test that OptionalStringP returns the value of the specified key when it exists in the map.
	r := OptionalStringP("test", "key1", m, "default")

	if r != "value1" {
		t.Errorf("expected 'value1', but got '%s'", r)
	}

	r = OptionalStringP("test", "key2", m, "default")

	if r != "value2" {
		t.Errorf("expected 'value2', but got '%s'", r)
	}

	// Test that OptionalStringP returns the default value when the specified key does not exist in the map.
	r = OptionalStringP("test", "key3", m, "default")

	if r != "default" {
		t.Errorf("expected 'default', but got '%s'", r)
	}

	// Test that OptionalStringP panics when the specified key exists in the map but cannot be parsed as a string.
	m["key4"] = 123
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key4 must be a string, but now it is a int(123)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalStringP("test", "key4", m, "default")
}

func TestStringArrayValueP_panic(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-slice-value",
	}

	// Test that StringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a string array, but now it is a string(not-a-slice-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	StringArrayValueP("task", "key1", m)
}

func TestStringArrayValueP(t *testing.T) {
	m := map[string]any{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3", "value4"},
	}

	// Test that StringArrayValueP returns the value of the specified key when it exists in the map.
	r := StringArrayValueP("test", "key1", m)

	if !reflect.DeepEqual(r, []string{"value1", "value2"}) {
		t.Errorf("expected ['value1', 'value2'], but got '%v'", r)
	}

	r = StringArrayValueP("test", "key2", m)

	if !reflect.DeepEqual(r, []string{"value3", "value4"}) {
		t.Errorf("expected ['value3', 'value4'], but got '%v'", r)
	}

	// Test that StringArrayValueP panics when the specified key does not exist in the map.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key3 is required"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	StringArrayValueP("test", "key3", m)
}

func TestOptionalStringArrayValueP_panic(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-slice-value",
	}

	// Test that OptionalStringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a string array, but now it is a string(not-a-slice-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalStringArrayValueP("task", "key1", m, []string{"default"})
}

func TestOptionalStringArrayValueP(t *testing.T) {
	m := map[string]any{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3", "value4"},
	}

	// Test that OptionalStringArrayValueP returns the value of the specified key when it exists in the map.
	r := OptionalStringArrayValueP("test", "key1", m, []string{"default"})

	if !reflect.DeepEqual(r, []string{"value1", "value2"}) {
		t.Errorf("expected ['value1', 'value2'], but got '%v'", r)
	}

	r = OptionalStringArrayValueP("test", "key2", m, []string{"default"})

	if !reflect.DeepEqual(r, []string{"value3", "value4"}) {
		t.Errorf("expected ['value3', 'value4'], but got '%v'", r)
	}

	// Test that OptionalStringArrayValueP returns the default value when the specified key does not exist in the map.
	r = OptionalStringArrayValueP("test", "key3", m, []string{"default"})

	if !reflect.DeepEqual(r, []string{"default"}) {
		t.Errorf("expected ['default'], but got '%v'", r)
	}

	// Test that OptionalStringArrayValueP panics when the specified key exists in the map but cannot be parsed as a string slice.
	m["key4"] = "not-a-slice-value"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key4 must be a string array, but now it is a string(not-a-slice-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalStringArrayValueP("test", "key4", m, []string{"default"})
}
