package test

import (
	comm "github.com/qiangyt/go-comm/v2"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestProgressReader(t *testing.T) {
	content := []byte("hello world, this is a test content for progress reader")
	reader := bytes.NewReader(content)

	var progressCalled bool
	var lastTransferred int64

	onProgress := func(transferred, total int64, speed float64) {
		progressCalled = true
		lastTransferred = transferred

		if transferred > total && total > 0 {
			t.Errorf("Transferred (%d) should not exceed total (%d)", transferred, total)
		}
	}

	pr := comm.NewProgressReader(reader, int64(len(content)), onProgress)

	// 读取所有内容
	buf := make([]byte, 10)
	var totalRead int64
	for {
		n, err := pr.Read(buf)
		totalRead += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read error: %v", err)
		}
	}

	// 触发最后一次进度报告
	pr.Finish()

	if totalRead != int64(len(content)) {
		t.Errorf("Total read = %d, want %d", totalRead, len(content))
	}

	// 验证进度回调被调用
	if progressCalled {
		if lastTransferred != totalRead {
			t.Errorf("Last transferred = %d, want %d", lastTransferred, totalRead)
		}
	}

	// 获取统计信息
	transferred, total, avgSpeed, duration := pr.GetStatistics()
	if transferred != totalRead {
		t.Errorf("Statistics transferred = %d, want %d", transferred, totalRead)
	}
	if total != int64(len(content)) {
		t.Errorf("Statistics total = %d, want %d", total, len(content))
	}
	if avgSpeed < 0 {
		t.Error("Average speed should be non-negative")
	}
	if duration <= 0 {
		t.Error("Duration should be positive")
	}
}

func TestProgressWriter(t *testing.T) {
	var buf bytes.Buffer
	content := []byte("hello world, this is a test content for progress writer")

	var progressCalled bool
	var lastTransferred int64

	onProgress := func(transferred, total int64, speed float64) {
		progressCalled = true
		lastTransferred = transferred

		if transferred > total && total > 0 {
			t.Errorf("Transferred (%d) should not exceed total (%d)", transferred, total)
		}
	}

	pw := comm.NewProgressWriter(&buf, int64(len(content)), onProgress)

	// 写入所有内容
	chunkSize := 10
	var totalWritten int64
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		n, err := pw.Write(content[i:end])
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
		totalWritten += int64(n)
	}

	// 触发最后一次进度报告
	pw.Finish()

	if totalWritten != int64(len(content)) {
		t.Errorf("Total written = %d, want %d", totalWritten, len(content))
	}

	// 验证写入的内容正确
	if buf.String() != string(content) {
		t.Errorf("Written content = %s, want %s", buf.String(), content)
	}

	// 验证进度回调被调用
	if progressCalled {
		if lastTransferred != totalWritten {
			t.Errorf("Last transferred = %d, want %d", lastTransferred, totalWritten)
		}
	}

	// 获取统计信息
	transferred, total, avgSpeed, duration := pw.GetStatistics()
	if transferred != totalWritten {
		t.Errorf("Statistics transferred = %d, want %d", transferred, totalWritten)
	}
	if total != int64(len(content)) {
		t.Errorf("Statistics total = %d, want %d", total, len(content))
	}
	if avgSpeed < 0 {
		t.Error("Average speed should be non-negative")
	}
	if duration <= 0 {
		t.Error("Duration should be positive")
	}
}

func TestProgressReaderNilCallback(t *testing.T) {
	content := []byte("test content")
	reader := bytes.NewReader(content)

	// 使用nil回调
	pr := comm.NewProgressReader(reader, int64(len(content)), nil)

	// 应该能正常读取
	buf := make([]byte, len(content))
	n, err := pr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != len(content) {
		t.Errorf("Read = %d bytes, want %d", n, len(content))
	}
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
		result := comm.FormatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("comm.FormatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
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
		result := comm.FormatBits(tt.bits)
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
		result := comm.FormatDuration(tt.seconds)
		if result != tt.expected {
			t.Errorf("FormatDuration(%f) = %s, want %s", tt.seconds, result, tt.expected)
		}
	}
}

func TestFormatSpeed(t *testing.T) {
	bytesPerSec := 1024.0

	// 测试不同格式
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
		result := comm.FormatSpeed(bytesPerSec, tt.format)
		if !strings.Contains(result, "/s") {
			t.Errorf("FormatSpeed(%f, %s) = %s, should contain '/s'", bytesPerSec, tt.format, result)
		}
	}

	// 测试bits格式
	bitsResult := comm.FormatSpeed(125, "bits") // 125 bytes/s = 1000 bits/s
	if !strings.Contains(bitsResult, "bit") {
		t.Errorf("FormatSpeed with 'bits' format should contain 'bit', got %s", bitsResult)
	}
}

func TestProgressReaderSlowUpdate(t *testing.T) {
	// 测试进度更新间隔
	content := bytes.Repeat([]byte("a"), 1000)
	reader := bytes.NewReader(content)

	callCount := 0
	onProgress := func(transferred, total int64, speed float64) {
		callCount++
	}

	pr := comm.NewProgressReader(reader, int64(len(content)), onProgress)

	// 快速读取所有内容（应该不会触发多次进度更新）
	io.Copy(io.Discard, pr)

	// 由于读取速度快，可能不会触发间隔更新
	// 但调用Finish应该触发一次
	pr.Finish()

	if callCount > 100 {
		t.Errorf("Progress callback called too many times: %d", callCount)
	}
}

func TestProgressReaderZeroTotal(t *testing.T) {
	content := []byte("test content")
	reader := bytes.NewReader(content)

	// total为0的情况（未知大小）
	pr := comm.NewProgressReader(reader, 0, func(transferred, total int64, speed float64) {
		if total != 0 {
			t.Errorf("Total should be 0, got %d", total)
		}
	})

	buf := make([]byte, len(content))
	pr.Read(buf)
	pr.Finish()
}
