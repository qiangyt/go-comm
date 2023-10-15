package comm

import (
	"testing"
)

func TestIntP(t *testing.T) {
	m := map[string]any{
		"key1": 7,
	}

	// Test that IntP panics when the specified key does not exist in the map.
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

	// Test that IntP returns the value of the specified key when it exists in the map.
	r := RequiredIntP("task", "key1", m)
	if r != 7 {
		t.Errorf("expected 7, but got %d", r)
	}

	// Test that IntP returns an error when the specified key does not exist in the map
	RequiredIntP("task", "not-existed-key", m)
}

func TestIntP_parsingError(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-int-value",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a int, but now it is a string(not-a-int-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()

	// Test that failed to parse int value
	RequiredIntP("task", "key1", m)
}

func TestInt(t *testing.T) {
	m := map[string]any{
		"key1": 5,
	}

	// Test that IntP returns the value of the specified key when it exists in the map.
	r, err := RequiredInt("task", "key1", m)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}
	if r != 5 {
		t.Errorf("expected 5, but got %d", r)
	}

	// Test that IntP returns an error when the specified key does not exist in the map
	_, err = RequiredInt("task", "not-existed-key", m)
	if err == nil {
		t.Error("expected error message 'task.not-existed-key is required', but no")
	}
}

func TestInt_parsingError(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-int-value",
	}

	// Test that failed to parse int value
	_, err := RequiredInt("task", "key1", m)
	if err == nil {
		t.Errorf("expected error but no")
	}
	expected := "task.key1 must be a int, but now it is a string(not-a-int-value)"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
	}
}

func TestOptionalIntP_panic(t *testing.T) {
	m := map[string]any{
		"key1": "not-a-int-value",
	}

	// Test that OptionalIntP panics when the specified key exists in the map but cannot be parsed as a int.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "task.key1 must be a int, but now it is a string(not-a-int-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalIntP("task", "key1", m, 9)
}

func TestOptionalIntP(t *testing.T) {
	m := map[string]any{
		"key1": 10,
	}

	// Test that OptionalIntP returns the value of the specified key when it exists in the map.
	r := OptionalIntP("test", "key1", m, 2)
	if r != 10 {
		t.Errorf("expected 10, but got %d", r)
	}

	// Test that OptionalIntP returns the default value when the specified key does not exist in the map.
	r = OptionalIntP("test", "key3", m, 3)
	if r != 3 {
		t.Errorf("expected 3, but got %d", r)
	}

	// Test that OptionalIntP panics when the specified key exists in the map but cannot be parsed as a int.
	m["key4"] = "not-a-int-value"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, but got nil")
		} else {
			err, ok := r.(error)
			if !ok {
				t.Errorf("expected error, but got %v", r)
			}
			expected := "test.key4 must be a int, but now it is a string(not-a-int-value)"
			if err.Error() != expected {
				t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
			}
		}
	}()
	OptionalIntP("test", "key4", m, 4)
}

func TestOptionalInt(t *testing.T) {
	m := map[string]any{
		"key1": 6,
	}

	// Test that OptionalInt returns the value of the specified key when it exists in the map.
	r, err := OptionalInt("test", "key1", m, -1)
	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
	if r != 6 {
		t.Errorf("expected 6, but got %d", r)
	}

	// Test that OptionalInt returns the default value when the specified key does not exist in the map.
	r, err = OptionalInt("test", "key3", m, -2)
	if err != nil {
		t.Errorf("expected nil error, but got %v", err)
	}
	if r != -2 {
		t.Errorf("expected -2, but got %d", r)
	}
}

func TestOptionalInt_error(t *testing.T) {
	m := map[string]any{
		"key4": "not-a-int-value",
	}

	// Test that OptionalInt returns an error when the specified key exists in the map but cannot be parsed as a int.
	_, err := OptionalInt("test", "key4", m, -3)

	if err == nil {
		t.Errorf("expected error, but got nil")
	}

	expected := "test.key4 must be a int, but now it is a string(not-a-int-value)"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', but got '%s'", expected, err.Error())
	}
}
