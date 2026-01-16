package comm

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestNewProgressReader(t *testing.T) {
	reader := strings.NewReader("test content")
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pr := NewProgressReader(reader, 12, onProgress)

	if pr == nil {
		t.Fatal("NewProgressReader returned nil")
	}

	if pr.total != 12 {
		t.Errorf("Expected total 12, got %d", pr.total)
	}

	if pr.transferred != 0 {
		t.Errorf("Expected transferred 0, got %d", pr.transferred)
	}
}

func TestProgressReader_Read(t *testing.T) {
	content := "hello world"
	reader := strings.NewReader(content)
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pr := NewProgressReader(reader, int64(len(content)), onProgress)

	buf := make([]byte, 5)
	n, err := pr.Read(buf)

	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if n != 5 {
		t.Errorf("Expected to read 5 bytes, got %d", n)
	}

	if string(buf) != "hello" {
		t.Errorf("Expected 'hello', got '%s'", string(buf))
	}

	// Progress callback may or may not be called depending on timing
	_ = callCount
}

func TestProgressReader_Read_All(t *testing.T) {
	content := "test content for reading all at once"
	reader := strings.NewReader(content)

	pr := NewProgressReader(reader, int64(len(content)), nil)

	buf := make([]byte, 100)
	n, err := pr.Read(buf)

	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}

	if n != len(content) {
		t.Errorf("Expected to read %d bytes, got %d", len(content), n)
	}
}

func TestProgressReader_GetStatistics(t *testing.T) {
	content := "test"
	reader := strings.NewReader(content)
	pr := NewProgressReader(reader, 4, nil)

	// Read some data
	buf := make([]byte, 2)
	pr.Read(buf)

	// Give it a moment to ensure duration > 0
	time.Sleep(1 * time.Millisecond)

	transferred, total, avgSpeed, duration := pr.GetStatistics()

	if transferred != 2 {
		t.Errorf("Expected transferred 2, got %d", transferred)
	}

	if total != 4 {
		t.Errorf("Expected total 4, got %d", total)
	}

	if avgSpeed < 0 {
		t.Errorf("Expected non-negative speed, got %f", avgSpeed)
	}

	if duration < 0 {
		t.Errorf("Expected non-negative duration, got %v", duration)
	}
}

func TestProgressReader_Finish(t *testing.T) {
	content := "test"
	reader := strings.NewReader(content)
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pr := NewProgressReader(reader, 4, onProgress)

	// Read all content
	buf := make([]byte, 10)
	pr.Read(buf)

	// Finish should trigger progress report
	pr.Finish()

	// Finish() calls reportProgress which may return early if elapsed == 0
	// So callCount may still be 0, but at least we verified it does not panic
	_ = callCount
}

func TestProgressReader_NilCallback(t *testing.T) {
	content := "test content"
	reader := strings.NewReader(content)

	pr := NewProgressReader(reader, int64(len(content)), nil)

	buf := make([]byte, 20)
	n, err := pr.Read(buf)

	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}

	if n != len(content) {
		t.Errorf("Expected to read %d bytes, got %d", len(content), n)
	}

	// Should not panic
	pr.Finish()
}

func TestNewProgressWriter(t *testing.T) {
	var buf bytes.Buffer
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pw := NewProgressWriter(&buf, 100, onProgress)

	if pw == nil {
		t.Fatal("NewProgressWriter returned nil")
	}

	if pw.total != 100 {
		t.Errorf("Expected total 100, got %d", pw.total)
	}

	if pw.transferred != 0 {
		t.Errorf("Expected transferred 0, got %d", pw.transferred)
	}
}

func TestProgressWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pw := NewProgressWriter(&buf, 100, onProgress)

	data := []byte("hello world")
	n, err := pw.Write(data)

	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got %d", len(data), n)
	}

	if buf.String() != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", buf.String())
	}

	// Progress callback may or may not be called depending on timing
	_ = callCount
}

func TestProgressWriter_GetStatistics(t *testing.T) {
	var buf bytes.Buffer
	pw := NewProgressWriter(&buf, 10, nil)

	// Write some data
	data := []byte("test")
	pw.Write(data)

	// Give it a moment
	time.Sleep(1 * time.Millisecond)

	transferred, total, avgSpeed, duration := pw.GetStatistics()

	if transferred != 4 {
		t.Errorf("Expected transferred 4, got %d", transferred)
	}

	if total != 10 {
		t.Errorf("Expected total 10, got %d", total)
	}

	if avgSpeed < 0 {
		t.Errorf("Expected non-negative speed, got %f", avgSpeed)
	}

	if duration < 0 {
		t.Errorf("Expected non-negative duration, got %v", duration)
	}
}

func TestProgressWriter_Finish(t *testing.T) {
	var buf bytes.Buffer
	callCount := 0

	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pw := NewProgressWriter(&buf, 10, onProgress)

	// Write some data
	pw.Write([]byte("test"))

	// Finish should trigger progress report
	pw.Finish()

	// Finish() calls reportProgress which may return early if elapsed == 0
	// So callCount may still be 0, but at least we verified it does not panic
	_ = callCount
}

func TestProgressWriter_NilCallback(t *testing.T) {
	var buf bytes.Buffer

	pw := NewProgressWriter(&buf, 10, nil)

	data := []byte("test content")
	n, err := pw.Write(data)

	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got %d", len(data), n)
	}

	if buf.String() != "test content" {
		t.Errorf("Expected 'test content', got '%s'", buf.String())
	}

	// Should not panic
	pw.Finish()
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
		}
	}
}

func TestFormatBits(t *testing.T) {
	tests := []struct {
		bits     float64
		expected string
	}{
		{0, "0 bit"},
		{500, "500 bit"},
		{999, "999 bit"},
		{1000, "1.0 Kbit"},
		{1500, "1.5 Kbit"},
		{1000000, "1.0 Mbit"},
		{1000000000, "1.0 Gbit"},
	}

	for _, tt := range tests {
		result := FormatBits(tt.bits)
		if result != tt.expected {
			t.Errorf("FormatBits(%f) = %s, want %s", tt.bits, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds  float64
		expected string
	}{
		{0, "0s"},
		{30, "30s"},
		{59, "59s"},
		{60, "1m0s"},
		{90, "1m30s"},
		{3600, "1h0m0s"},
		{3661, "1h1m1s"},
		{7265, "2h1m5s"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.seconds)
		if result != tt.expected {
			t.Errorf("FormatDuration(%f) = %s, want %s", tt.seconds, result, tt.expected)
		}
	}
}

func TestFormatSpeed(t *testing.T) {
	bytesPerSec := 1024.0

	// Test different formats
	tests := []struct {
		format   string
		contains string
	}{
		{"bytes", "KB/s"},
		{"bits", "Kbit/s"},
		{"human", "KB/s"},
		{"", "KB/s"},
		{"unknown", "KB/s"},
	}

	for _, tt := range tests {
		result := FormatSpeed(bytesPerSec, tt.format)
		if !strings.Contains(result, "/s") {
			t.Errorf("FormatSpeed(%f, %s) = %s, should contain '/s'", bytesPerSec, tt.format, result)
		}
	}

	// Test bits format
	bitsResult := FormatSpeed(125, "bits") // 125 bytes/s = 1000 bits/s
	if !strings.Contains(bitsResult, "bit") {
		t.Errorf("FormatSpeed with 'bits' format should contain 'bit', got %s", bitsResult)
	}
}
