package qgin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ==================== BodyTruncateStrategy ====================

func TestBodyTruncateStrategy_Values(t *testing.T) {
	a := require.New(t)

	// 验证所有枚举值存在
	a.Equal(BodyTruncateNone, BodyTruncateStrategy(0))
	a.Equal(BodyTruncateFull, BodyTruncateStrategy(1))
	a.Equal(BodyTruncateHead, BodyTruncateStrategy(2))
	a.Equal(BodyTruncateTail, BodyTruncateStrategy(3))
	a.Equal(BodyTruncateHeadAndTail, BodyTruncateStrategy(4))
}

func TestBodyTruncateStrategy_Default(t *testing.T) {
	// 默认值应该是 HeadAndTail
	cfg := DefaultBodyLogConfig()
	require.Equal(t, BodyTruncateHeadAndTail, cfg.Strategy)
}

// ==================== BodyLogConfig ====================

func TestBodyLogConfig_Default(t *testing.T) {
	a := require.New(t)

	cfg := DefaultBodyLogConfig()

	a.Equal(BodyTruncateHeadAndTail, cfg.Strategy)
	a.Equal(1024, cfg.TruncateSize)
}

func TestBodyLogConfig_Custom(t *testing.T) {
	a := require.New(t)

	cfg := BodyLogConfig{
		Strategy:     BodyTruncateFull,
		TruncateSize: 2048,
	}

	a.Equal(BodyTruncateFull, cfg.Strategy)
	a.Equal(2048, cfg.TruncateSize)
}

// ==================== HeaderLogStrategy ====================

func TestHeaderLogStrategy_Values(t *testing.T) {
	a := require.New(t)

	// 验证所有枚举值存在
	a.Equal(HeaderLogNone, HeaderLogStrategy(0))
	a.Equal(HeaderLogAll, HeaderLogStrategy(1))
	a.Equal(HeaderLogWhitelist, HeaderLogStrategy(2))
	a.Equal(HeaderLogBlacklist, HeaderLogStrategy(3))
}

func TestHeaderLogStrategy_Default(t *testing.T) {
	// 默认值应该是 All
	cfg := DefaultHeaderLogConfig()
	require.Equal(t, HeaderLogAll, cfg.Strategy)
}

// ==================== SensitiveHeaderStrategy ====================

func TestSensitiveHeaderStrategy_Values(t *testing.T) {
	a := require.New(t)

	// 验证所有枚举值存在
	a.Equal(SensitiveHeaderFull, SensitiveHeaderStrategy(0))
	a.Equal(SensitiveHeaderExclude, SensitiveHeaderStrategy(1))
	a.Equal(SensitiveHeaderMaskAll, SensitiveHeaderStrategy(2))
	a.Equal(SensitiveHeaderMaskHead, SensitiveHeaderStrategy(3))
	a.Equal(SensitiveHeaderMaskTail, SensitiveHeaderStrategy(4))
}

// ==================== SensitiveHeaderConfig ====================

func TestSensitiveHeaderConfig_Default(t *testing.T) {
	a := require.New(t)

	cfg := DefaultSensitiveHeaderConfig()

	a.Equal(SensitiveHeaderMaskAll, cfg.Strategy)
	a.Equal(4, cfg.MaskSize)
	a.NotNil(cfg.SensitiveList)
	// 验证默认敏感 header 列表
	a.Contains(cfg.SensitiveList, "Authorization")
	a.Contains(cfg.SensitiveList, "Cookie")
	a.Contains(cfg.SensitiveList, "Set-Cookie")
	a.Contains(cfg.SensitiveList, "X-Api-Key")
	a.Contains(cfg.SensitiveList, "X-Auth-Token")
}

func TestSensitiveHeaderConfig_Custom(t *testing.T) {
	a := require.New(t)

	cfg := SensitiveHeaderConfig{
		Strategy:      SensitiveHeaderExclude,
		MaskSize:      8,
		SensitiveList: []string{"Custom-Header"},
	}

	a.Equal(SensitiveHeaderExclude, cfg.Strategy)
	a.Equal(8, cfg.MaskSize)
	a.Len(cfg.SensitiveList, 1)
}

// ==================== HeaderLogConfig ====================

func TestHeaderLogConfig_Default(t *testing.T) {
	a := require.New(t)

	cfg := DefaultHeaderLogConfig()

	a.Equal(HeaderLogAll, cfg.Strategy)
	a.Nil(cfg.HeaderList)
	a.NotNil(cfg.SensitiveConfig)
}

func TestHeaderLogConfig_Custom(t *testing.T) {
	a := require.New(t)

	cfg := HeaderLogConfig{
		Strategy:   HeaderLogWhitelist,
		HeaderList: []string{"Content-Type", "Authorization"},
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy: SensitiveHeaderExclude,
		},
	}

	a.Equal(HeaderLogWhitelist, cfg.Strategy)
	a.Len(cfg.HeaderList, 2)
	a.NotNil(cfg.SensitiveConfig)
}

// ==================== SSETruncateStrategy ====================

func TestSSETruncateStrategy_Values(t *testing.T) {
	a := require.New(t)

	// 验证所有枚举值存在
	a.Equal(SSETruncateNone, SSETruncateStrategy(0))
	a.Equal(SSETruncateFull, SSETruncateStrategy(1))
	a.Equal(SSETruncateHead, SSETruncateStrategy(2))
	a.Equal(SSETruncateTail, SSETruncateStrategy(3))
	a.Equal(SSETruncateHeadAndTail, SSETruncateStrategy(4))
}

// ==================== SSELogConfig ====================

func TestSSELogConfig_Default(t *testing.T) {
	a := require.New(t)

	cfg := DefaultSSELogConfig()

	a.Equal(SSETruncateHeadAndTail, cfg.Strategy)
	a.Equal(10, cfg.TruncateSize)
}

func TestSSELogConfig_Custom(t *testing.T) {
	a := require.New(t)

	cfg := SSELogConfig{
		Strategy:     SSETruncateFull,
		TruncateSize: 20,
	}

	a.Equal(SSETruncateFull, cfg.Strategy)
	a.Equal(20, cfg.TruncateSize)
}

// ==================== GinLoggerConfig ====================

func TestGinLoggerConfig_Default(t *testing.T) {
	a := require.New(t)

	cfg := DefaultGinLoggerConfig()

	// Logger 默认为 nil，需要用户设置
	a.Nil(cfg.Logger)
	a.Nil(cfg.SkipPaths)
	a.Equal("X-Trace-Id", cfg.TraceIdHeader)

	// Body 配置默认值
	a.Equal(BodyTruncateHeadAndTail, cfg.RequestBody.Strategy)
	a.Equal(1024, cfg.RequestBody.TruncateSize)
	a.Equal(BodyTruncateHeadAndTail, cfg.ResponseBody.Strategy)
	a.Equal(1024, cfg.ResponseBody.TruncateSize)

	// Header 配置默认值
	a.Equal(HeaderLogAll, cfg.RequestHeader.Strategy)
	a.Equal(HeaderLogAll, cfg.ResponseHeader.Strategy)

	// SSE 配置默认值
	a.Equal(SSETruncateHeadAndTail, cfg.SSEConfig.Strategy)
	a.Equal(10, cfg.SSEConfig.TruncateSize)
}

func TestGinLoggerConfig_Custom(t *testing.T) {
	a := require.New(t)

	logger := &mockLogger{}
	cfg := GinLoggerConfig{
		Logger:        logger,
		SkipPaths:     []string{"/health", "/metrics"},
		TraceIdHeader: "X-Request-Id",
		CustomFields: func(ctx any) map[string]any {
			return map[string]any{"custom": "field"}
		},
		RequestBody: BodyLogConfig{
			Strategy:     BodyTruncateFull,
			TruncateSize: 2048,
		},
		ResponseBody: BodyLogConfig{
			Strategy:     BodyTruncateNone,
			TruncateSize: 0,
		},
	}

	a.Equal(logger, cfg.Logger)
	a.Len(cfg.SkipPaths, 2)
	a.Equal("X-Request-Id", cfg.TraceIdHeader)
	a.NotNil(cfg.CustomFields)
	a.Equal(BodyTruncateFull, cfg.RequestBody.Strategy)
	a.Equal(BodyTruncateNone, cfg.ResponseBody.Strategy)
}

// mockLogger 用于测试的 mock logger
type mockLogger struct {
	infoCalls  []logCall
	warnCalls  []logCall
	errorCalls []logCall
	debugCalls []logCall
	traceCalls []logCall
}

type logCall struct {
	msg    string
	fields map[string]any
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		infoCalls:  make([]logCall, 0),
		warnCalls:  make([]logCall, 0),
		errorCalls: make([]logCall, 0),
		debugCalls: make([]logCall, 0),
		traceCalls: make([]logCall, 0),
	}
}

func fieldsToMap(fields []any) map[string]any {
	result := make(map[string]any)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}
	return result
}

func (m *mockLogger) Info(msg string, fields ...any) {
	m.infoCalls = append(m.infoCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Warn(msg string, fields ...any) {
	m.warnCalls = append(m.warnCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Error(msg string, fields ...any) {
	m.errorCalls = append(m.errorCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Debug(msg string, fields ...any) {
	m.debugCalls = append(m.debugCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Trace(msg string, fields ...any) {
	m.traceCalls = append(m.traceCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Fatal(msg string, fields ...any) {
	m.errorCalls = append(m.errorCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) Panic(msg string, fields ...any) {
	m.errorCalls = append(m.errorCalls, logCall{msg: msg, fields: fieldsToMap(fields)})
}

func (m *mockLogger) WithField(key string, value any) any { return m }
func (m *mockLogger) WithFields(fields map[string]any) any { return m }

func (m *mockLogger) getLastInfoCall() *logCall {
	if len(m.infoCalls) == 0 {
		return nil
	}
	return &m.infoCalls[len(m.infoCalls)-1]
}

func (m *mockLogger) getLastWarnCall() *logCall {
	if len(m.warnCalls) == 0 {
		return nil
	}
	return &m.warnCalls[len(m.warnCalls)-1]
}

func (m *mockLogger) getLastErrorCall() *logCall {
	if len(m.errorCalls) == 0 {
		return nil
	}
	return &m.errorCalls[len(m.errorCalls)-1]
}

func (m *mockLogger) getInfoCallsCount() int { return len(m.infoCalls) }
func (m *mockLogger) getWarnCallsCount() int { return len(m.warnCalls) }
func (m *mockLogger) getErrorCallsCount() int { return len(m.errorCalls) }
