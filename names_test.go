package comm

import "testing"

func TestNameOfKey(t *testing.T) {
	// Test that NameOfKey returns the correct name when the key contains a period.
	r := NameOfKey("util.NameOfKey")
	if r != "NameOfKey" {
		t.Errorf("expected 'NameOfKey', but got '%s'", r)
	}

	// Test that NameOfKey returns the correct name when the key does not contain a period.
	r = NameOfKey("key")
	if r != "key" {
		t.Errorf("expected 'key', but got '%s'", r)
	}
}
