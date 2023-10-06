package comm

import (
	"testing"
)

func TestBoolP(t *testing.T) {
	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that BoolP panics when the specified key does not exist in the map.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.not-existed-key is required"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()

	// Test that BoolP returns the value of the specified key when it exists in the map.
	r := RequiredBoolP("task", "key1", m)
	if !r {
		t.Errorf("expected true, but got false")
	}

	r = RequiredBoolP("task", "key2", m)
	if r {
		t.Errorf("expected false, but got true")
	}

	// Test that BoolP returns an error when the specified key does not exist in the map
	RequiredBoolP("task", "not-existed-key", m)
}

func TestBoolP_parsingError(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a bool, but now it is a string(not-a-bool-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()

	// Test that failed to parse bool value
	RequiredBoolP("task", "key1", m)
}

func TestBool(t *testing.T) {
	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that BoolP returns the value of the specified key when it exists in the map.
	r, err := RequiredBool("task", "key1", m)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if !r {
		t.Errorf("expected true, but got false")
	}

	r, err = RequiredBool("task", "key2", m)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if r {
		t.Errorf("expected false, but got true")
	}

	// Test that BoolP returns an error when the specified key does not exist in the map
	_, err = RequiredBool("task", "not-existed-key", m)
	if err == nil {
		t.Error("expected error message 'task.not-existed-key is required', but no")
	}
}

func TestBool_parsingError(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	// Test that failed to parse bool value
	_, err := RequiredBool("task", "key1", m)
	if err == nil {
		t.Errorf("expected error but no")
	}
	expected := "task.key1 must be a bool, but now it is a string(not-a-bool-value)"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
	}
}

func TestOptionalBoolP_panic(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-bool-value",
	}

	// Test that OptionalBoolP panics when the specified key exists in the map but cannot be parsed as a bool.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a bool, but now it is a string(not-a-bool-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalBoolP("task", "key1", m, false)
}

func TestOptionalBoolP(t *testing.T) {
	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that OptionalBoolP returns the value of the specified key when it exists in the map.
	r := OptionalBoolP("test", "key1", m, false)
	if !r {
		t.Errorf("expected true, but got false")
	}

	r = OptionalBoolP("test", "key2", m, true)
	if r {
		t.Errorf("expected false, but got true")
	}

	// Test that OptionalBoolP returns the default value when the specified key does not exist in the map.
	r = OptionalBoolP("test", "key3", m, true)
	if !r {
		t.Errorf("expected true, but got false")
	}

	// Test that OptionalBoolP panics when the specified key exists in the map but cannot be parsed as a bool.
	m["key4"] = "not-a-bool-value"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key4 must be a bool, but now it is a string(not-a-bool-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalBoolP("test", "key4", m, false)
}

func TestOptionalBool(t *testing.T) {
	m := map[string]any{
		"key1": true,
		"key2": false,
	}

	// Test that OptionalBool returns the value of the specified key when it exists in the map.
	r, err := OptionalBool("test", "key1", m, false)
	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
	if !r {
		t.Errorf("expected true, but got false")
	}

	r, err = OptionalBool("test", "key2", m, true)
	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
	if r {
		t.Errorf("expected false, but got true")
	}

	// Test that OptionalBool returns the default value when the specified key does not exist in the map.
	r, err = OptionalBool("test", "key3", m, true)
	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
	if !r {
		t.Errorf("expected true, but got false")
	}
}

func TestOptionalBool_error(t *testing.T) {
	m := map[string]any{
		"key4": "not-a-bool-value",
	}

	// Test that OptionalBool returns an error when the specified key exists in the map but cannot be parsed as a bool.
	_, err := OptionalBool("test", "key4", m, false)

	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	expected := "test.key4 must be a bool, but now it is a string(not-a-bool-value)"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
	}
}
