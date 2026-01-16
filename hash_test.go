package comm

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewHashCalculator(t *testing.T) {
	calc := NewHashCalculator()
	if calc == nil {
		t.Fatal("NewHashCalculator returned nil")
	}
}

func TestHashCalculator_CalculateMD5(t *testing.T) {
	calc := NewHashCalculator()

	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("hello world")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate MD5
	md5Hash := calc.CalculateMD5(testFile)
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3" // MD5 of "hello world"

	if md5Hash != expected {
		t.Errorf("CalculateMD5() = %s, want %s", md5Hash, expected)
	}
}

func TestHashCalculator_CalculateSHA256(t *testing.T) {
	calc := NewHashCalculator()

	// Create test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("hello world")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate SHA256
	sha256Hash := calc.CalculateSHA256(testFile)
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9" // SHA256 of "hello world"

	if sha256Hash != expected {
		t.Errorf("CalculateSHA256() = %s, want %s", sha256Hash, expected)
	}
}

func TestHashCalculator_CalculateMD5FromReader(t *testing.T) {
	calc := NewHashCalculator()

	testContent := []byte("hello world")
	reader := bytes.NewReader(testContent)

	md5Hash := calc.CalculateMD5FromReader(reader)
	expected := "5eb63bbbe01eeed093cb22bb8f5acdc3"

	if md5Hash != expected {
		t.Errorf("CalculateMD5FromReader() = %s, want %s", md5Hash, expected)
	}
}

func TestHashCalculator_CalculateSHA256FromReader(t *testing.T) {
	calc := NewHashCalculator()

	testContent := []byte("hello world")
	reader := bytes.NewReader(testContent)

	sha256Hash := calc.CalculateSHA256FromReader(reader)
	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	if sha256Hash != expected {
		t.Errorf("CalculateSHA256FromReader() = %s, want %s", sha256Hash, expected)
	}
}

func TestHashCalculator_EmptyFile(t *testing.T) {
	calc := NewHashCalculator()

	// Create empty file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(testFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// MD5 of empty string
	md5Hash := calc.CalculateMD5(testFile)
	expectedMD5 := "d41d8cd98f00b204e9800998ecf8427e"
	if md5Hash != expectedMD5 {
		t.Errorf("CalculateMD5() for empty file = %s, want %s", md5Hash, expectedMD5)
	}

	// SHA256 of empty string
	sha256Hash := calc.CalculateSHA256(testFile)
	expectedSHA256 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if sha256Hash != expectedSHA256 {
		t.Errorf("CalculateSHA256() for empty file = %s, want %s", sha256Hash, expectedSHA256)
	}
}

func TestHashCalculator_NonExistentFile(t *testing.T) {
	calc := NewHashCalculator()

	// Test non-existent file should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("CalculateMD5() with non-existent file should panic")
		}
	}()

	calc.CalculateMD5("/nonexistent/file.txt")
}

func TestHashCalculator_LargeFile(t *testing.T) {
	calc := NewHashCalculator()

	// Create a larger file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	// Write 1MB of data
	largeContent := bytes.Repeat([]byte("a"), 1024*1024)
	if err := os.WriteFile(testFile, largeContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate hashes
	md5Hash := calc.CalculateMD5(testFile)
	sha256Hash := calc.CalculateSHA256(testFile)

	// Verify hashes are not empty
	if md5Hash == "" {
		t.Error("CalculateMD5() should return non-empty hash for large file")
	}
	if sha256Hash == "" {
		t.Error("CalculateSHA256() should return non-empty hash for large file")
	}

	// Verify hash lengths
	if len(md5Hash) != 32 {
		t.Errorf("MD5 hash length = %d, want 32", len(md5Hash))
	}
	if len(sha256Hash) != 64 {
		t.Errorf("SHA256 hash length = %d, want 64", len(sha256Hash))
	}
}
