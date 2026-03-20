package comm

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/qiangyt/go-event"
	"github.com/stretchr/testify/require"
)

func TestLogMapper_happy(t *testing.T) {
	// Create a real logger that returns non-nil entry
	tmpDir, err := os.MkdirTemp("", "logger-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{MaxSize: 10, MaxAge: 1, MaxBackups: 3}
	logger, err := NewLogger(nil, config, logFile)
	require.NoError(t, err)
	defer logger.Close()

	entry := logger.Info()
	mapper := &testMapper{data: map[string]any{"key": "value"}}
	result := LogMapper(entry, "test", mapper)
	require.NotNil(t, result)
}

type testMapper struct {
	data map[string]any
}

func (m *testMapper) ToMap() map[string]any {
	return m.data
}

func TestLogMap_happy(t *testing.T) {
	// Create a real logger that returns non-nil entry
	tmpDir, err := os.MkdirTemp("", "logger-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{MaxSize: 10, MaxAge: 1, MaxBackups: 3}
	logger, err := NewLogger(nil, config, logFile)
	require.NoError(t, err)
	defer logger.Close()

	entry := logger.Info()
	result := LogMap(entry, "test", map[string]any{"key": "value"})
	require.NotNil(t, result)
}

func TestNewLogContext_withTraceId(t *testing.T) {
	a := require.New(t)

	ctx := NewLogContext(true)
	a.NotNil(ctx)
}

func TestNewLogContext_withoutTraceId(t *testing.T) {
	a := require.New(t)

	ctx := NewLogContext(false)
	a.NotNil(ctx)
}

func TestNewSubLogger_happy(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	ctx := NewLogContext(false)
	subLogger := logger.NewSubLogger(ctx)

	a.NotNil(subLogger)
	a.Equal(logger, subLogger.Parent())
}

func TestNewSubLogger_nilContext(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	subLogger := logger.NewSubLogger(nil)

	a.NotNil(subLogger)
}

func TestLogger_Parent(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	a.Nil(logger.Parent())

	subLogger := logger.NewSubLogger(nil)
	a.Equal(logger, subLogger.Parent())
}

func TestLogger_Close(t *testing.T) {
	// Test that Close doesn't panic on discard logger
	logger := NewDiscardLogger()
	logger.Close()
}

func TestLogger_Error(t *testing.T) {
	// Create a real logger that returns non-nil entry
	tmpDir, err := os.MkdirTemp("", "logger-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{MaxSize: 10, MaxAge: 1, MaxBackups: 3}
	logger, err := NewLogger(nil, config, logFile)
	require.NoError(t, err)
	defer logger.Close()

	entry := logger.Error("test error")
	require.NotNil(t, entry)
}

func TestNewDiscardLogger(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	a.NotNil(logger)
	a.True(IsDiscardLogger(logger))
}

func TestIsDiscardLogger_true(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	a.True(IsDiscardLogger(logger))
}

func TestNewLogger_happy(t *testing.T) {
	a := require.New(t)

	tmpDir, err := os.MkdirTemp("", "logger-test")
	a.NoError(err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 3,
		LocalTime:  true,
		Compress:   false,
	}

	logger, err := NewLogger(nil, config, logFile)
	a.NoError(err)
	a.NotNil(logger)
	defer logger.Close()
}

func TestNewLogger_withConsole(t *testing.T) {
	a := require.New(t)

	tmpDir, err := os.MkdirTemp("", "logger-test")
	a.NoError(err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 3,
		LocalTime:  true,
		Compress:   false,
	}

	var console bytes.Buffer
	logger, err := NewLogger(&console, config, logFile)
	a.NoError(err)
	a.NotNil(logger)
	defer logger.Close()
}

func TestNewLoggerP_happy(t *testing.T) {
	a := require.New(t)

	tmpDir, err := os.MkdirTemp("", "logger-test")
	a.NoError(err)
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test.log")
	config := &LoggerConfigT{
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 3,
		LocalTime:  true,
		Compress:   false,
	}

	logger := NewLoggerP(nil, config, logFile)
	a.NotNil(logger)
	defer logger.Close()
}

func TestNewEventLogger_happy(t *testing.T) {
	a := require.New(t)

	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	a.NotNil(eventLogger)
	a.Equal(logger, eventLogger.Target())
}

func TestEventLogger_LogDebug(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogDebug(event.LogEnum(0), "hub", "topic", "listener")
}

func TestEventLogger_LogInfo(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogInfo(event.LogEnum(1), "hub", "topic", "listener")
}

func TestEventLogger_LogError(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogError(event.LogEnum(2), "hub", "topic", "listener", "test error")
}

func TestEventLogger_LogEventDebug(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogEventDebug(event.LogEnum(0), "listener", nil)
}

func TestEventLogger_LogEventInfo(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogEventInfo(event.LogEnum(1), "listener", nil)
}

func TestEventLogger_LogEventError(t *testing.T) {
	logger := NewDiscardLogger()
	eventLogger := NewEventLogger(logger)

	// Should not panic
	eventLogger.LogEventError(event.LogEnum(2), "listener", nil, "test error")
}

// Test discard logger outputs nothing
func TestDiscardLogger_noOutput(t *testing.T) {
	logger := NewDiscardLogger()

	// These should not produce any output or panic
	logger.Debug().Msg("debug message")
	logger.Info().Msg("info message")
	logger.Warn().Msg("warn message")
}

var _ io.Writer = (*testWriter)(nil)

type testWriter struct {
	buf bytes.Buffer
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}
