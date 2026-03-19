package qgin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== GinLogger 工厂函数测试 ====================

// TestGinLogger_ReturnsValidHandlerFunc 测试 GinLogger 返回有效的 HandlerFunc
func TestGinLogger_ReturnsValidHandlerFunc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	if handler == nil {
		t.Fatal("GinLogger() 返回 nil")
	}

	// 创建测试上下文
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// 调用中间件，不应该 panic
	handler(c)
}

// TestGinLogger_LogsRequest 测试 GinLogger 记录请求日志
func TestGinLogger_LogsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	// 验证日志被调用
	if logger.getInfoCallsCount() == 0 {
		t.Error("GinLogger() 没有记录日志")
	}
}

// ==================== GinLoggerWithConfig 配置工厂函数测试 ====================

// TestGinLoggerWithConfig_UsesCustomConfig 测试 GinLoggerWithConfig 使用自定义配置
func TestGinLoggerWithConfig_UsesCustomConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.SkipPaths = []string{"/health"}

	handler := GinLoggerWithConfig(config)
	if handler == nil {
		t.Fatal("GinLoggerWithConfig() 返回 nil")
	}

	// 测试跳过路径
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler(c)

	// /health 应该被跳过，不记录日志
	if logger.getInfoCallsCount() != 0 {
		t.Errorf("GinLoggerWithConfig() 应该跳过 /health 路径，但记录了 %d 条日志", logger.getInfoCallsCount())
	}
}

// TestGinLoggerWithConfig_NilConfig 测试 GinLoggerWithConfig 使用 nil 配置
func TestGinLoggerWithConfig_NilConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// nil 配置应该使用默认配置
	handler := GinLoggerWithConfig(nil)
	if handler == nil {
		t.Fatal("GinLoggerWithConfig(nil) 返回 nil")
	}
}

// ==================== 请求时间记录和延迟计算测试 ====================

// TestGinLogger_RecordsLatency 测试延迟计算正确
func TestGinLogger_RecordsLatency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	start := time.Now()
	handler(c)
	elapsed := time.Since(start)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 latency 字段存在
	latency, ok := lastCall.fields["latency"]
	if !ok {
		t.Error("日志记录缺少 latency 字段")
	}

	// 验证延迟值合理（应该小于测试时间）
	if latencyVal, ok := latency.(time.Duration); ok {
		if latencyVal > elapsed {
			t.Errorf("latency 值 %v 大于实际耗时 %v", latencyVal, elapsed)
		}
	} else {
		t.Errorf("latency 类型不正确: %T", latency)
	}
}

// ==================== 状态码到日志级别映射测试 ====================

// TestStatusToLevel_2xx 测试 2xx 状态码映射到 Info 级别
func TestStatusToLevel_2xx(t *testing.T) {
	statuses := []int{200, 201, 204, 299}
	for _, status := range statuses {
		level := statusToLevel(status)
		if level != "info" {
			t.Errorf("statusToLevel(%d) = %s, 期望 info", status, level)
		}
	}
}

// TestStatusToLevel_3xx 测试 3xx 状态码映射到 Info 级别
func TestStatusToLevel_3xx(t *testing.T) {
	statuses := []int{301, 302, 304, 399}
	for _, status := range statuses {
		level := statusToLevel(status)
		if level != "info" {
			t.Errorf("statusToLevel(%d) = %s, 期望 info", status, level)
		}
	}
}

// TestStatusToLevel_4xx 测试 4xx 状态码映射到 Warn 级别
func TestStatusToLevel_4xx(t *testing.T) {
	statuses := []int{400, 401, 403, 404, 499}
	for _, status := range statuses {
		level := statusToLevel(status)
		if level != "warn" {
			t.Errorf("statusToLevel(%d) = %s, 期望 warn", status, level)
		}
	}
}

// TestStatusToLevel_5xx 测试 5xx 状态码映射到 Error 级别
func TestStatusToLevel_5xx(t *testing.T) {
	statuses := []int{500, 502, 503, 504, 599}
	for _, status := range statuses {
		level := statusToLevel(status)
		if level != "error" {
			t.Errorf("statusToLevel(%d) = %s, 期望 error", status, level)
		}
	}
}

// TestGinLogger_LogLevelByStatus 测试根据状态码选择日志级别
func TestGinLogger_LogLevelByStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		statusCode  int
		expectInfo  bool
		expectWarn  bool
		expectError bool
	}{
		{"200 OK", 200, true, false, false},
		{"404 Not Found", 404, false, true, false},
		{"500 Internal Server Error", 500, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := newMockLogger()
			config := DefaultGinLoggerConfig()
			config.Logger = logger
			handler := GinLoggerWithConfig(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// 设置状态码
			c.Status(tt.statusCode)

			handler(c)

			// 手动触发响应写入（因为 gin 需要写入才能真正设置状态）
			w.WriteHeader(tt.statusCode)

			// 验证日志级别
			infoCount := logger.getInfoCallsCount()
			warnCount := logger.getWarnCallsCount()
			errorCount := logger.getErrorCallsCount()

			if tt.expectInfo && infoCount == 0 {
				t.Errorf("期望 Info 日志，但没有调用")
			}
			if tt.expectWarn && warnCount == 0 {
				t.Errorf("期望 Warn 日志，但没有调用")
			}
			if tt.expectError && errorCount == 0 {
				t.Errorf("期望 Error 日志，但没有调用")
			}
		})
	}
}

// ==================== 基本日志字段输出测试 ====================

// TestGinLogger_BasicFields 测试基本日志字段（method, path, status, latency, client_ip, body_size）
func TestGinLogger_BasicFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/test", nil)
	c.Request.Header.Set("X-Real-IP", "192.168.1.100")

	handler(c)

	// 手动设置状态码和 body size
	w.WriteHeader(200)
	w.Write([]byte("test response"))

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证基本字段
	expectedFields := []string{"method", "path", "status", "latency", "client_ip"}
	for _, field := range expectedFields {
		if _, ok := lastCall.fields[field]; !ok {
			t.Errorf("日志记录缺少 %s 字段", field)
		}
	}

	// 验证 method
	if method, ok := lastCall.fields["method"].(string); ok {
		if method != http.MethodPost {
			t.Errorf("method = %s, 期望 %s", method, http.MethodPost)
		}
	}

	// 验证 path
	if path, ok := lastCall.fields["path"].(string); ok {
		if path != "/api/test" {
			t.Errorf("path = %s, 期望 /api/test", path)
		}
	}
}

// TestGinLogger_BodySize 测试 body_size 字段
func TestGinLogger_BodySize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// 模拟响应
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte("hello world"))

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 body_size 字段存在
	if _, ok := lastCall.fields["body_size"]; !ok {
		t.Error("日志记录缺少 body_size 字段")
	}
}

// ==================== SkipPaths 功能测试 ====================

// TestGinLogger_SkipPaths 测试跳过指定路径
func TestGinLogger_SkipPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 测试应该跳过的路径
	for _, path := range []string{"/health", "/metrics"} {
		t.Run("skip_"+path, func(t *testing.T) {
			logger := newMockLogger()
			config := DefaultGinLoggerConfig()
			config.Logger = logger
			config.SkipPaths = []string{"/health", "/metrics"}
			handler := GinLoggerWithConfig(config)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, path, nil)

			handler(c)

			if logger.getInfoCallsCount() != 0 {
				t.Errorf("路径 %s 应该被跳过，但记录了 %d 条日志", path, logger.getInfoCallsCount())
			}
		})
	}

	// 测试不应该跳过的路径
	t.Run("log_/api/test", func(t *testing.T) {
		logger := newMockLogger()
		config := DefaultGinLoggerConfig()
		config.Logger = logger
		config.SkipPaths = []string{"/health", "/metrics"}
		handler := GinLoggerWithConfig(config)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)

		handler(c)

		if logger.getInfoCallsCount() == 0 {
			t.Error("路径 /api/test 不应该被跳过")
		}
	})
}

// ==================== TraceId 功能测试 ====================

// TestGinLogger_TraceId_Generate 测试自动生成 traceId
func TestGinLogger_TraceId_Generate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 trace_id 字段存在且非空
	traceId, ok := lastCall.fields["trace_id"]
	if !ok {
		t.Error("日志记录缺少 trace_id 字段")
	}
	if traceIdStr, ok := traceId.(string); ok {
		if traceIdStr == "" {
			t.Error("trace_id 不应该为空")
		}
	}
}

// TestGinLogger_TraceId_FromHeader 测试从请求头读取 traceId
func TestGinLogger_TraceId_FromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.TraceIdHeader = "X-Request-Id"
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Request-Id", "test-trace-123")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 trace_id 使用请求头中的值
	traceId, ok := lastCall.fields["trace_id"]
	if !ok {
		t.Fatal("日志记录缺少 trace_id 字段")
	}
	if traceIdStr, ok := traceId.(string); ok {
		if traceIdStr != "test-trace-123" {
			t.Errorf("trace_id = %s, 期望 test-trace-123", traceIdStr)
		}
	}
}

// TestGinLogger_TraceId_StoredInContext 测试 traceId 存入 gin.Context
func TestGinLogger_TraceId_StoredInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Trace-Id", "context-trace-456")

	handler(c)

	// 验证 trace_id 存入 context
	traceId := c.GetString("trace_id")
	if traceId != "context-trace-456" {
		t.Errorf("context trace_id = %s, 期望 context-trace-456", traceId)
	}
}

// ==================== CustomFields 功能测试 ====================

// TestGinLogger_CustomFields 测试 CustomFields 回调
func TestGinLogger_CustomFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.CustomFields = func(ctx any) map[string]any {
		return map[string]any{
			"user_id":    "user-123",
			"tenant_id":  "tenant-456",
			"custom_key": "custom_value",
		}
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证自定义字段
	if lastCall.fields["user_id"] != "user-123" {
		t.Errorf("user_id = %v, 期望 user-123", lastCall.fields["user_id"])
	}
	if lastCall.fields["tenant_id"] != "tenant-456" {
		t.Errorf("tenant_id = %v, 期望 tenant-456", lastCall.fields["tenant_id"])
	}
	if lastCall.fields["custom_key"] != "custom_value" {
		t.Errorf("custom_key = %v, 期望 custom_value", lastCall.fields["custom_key"])
	}
}

// ==================== gin.Context 设置测试 ====================

// TestGinLogger_SetStartTime 测试设置请求开始时间
func TestGinLogger_SetStartTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	// 验证开始时间存入 context
	startTime, exists := c.Get("_startTime")
	if !exists {
		t.Error("开始时间没有存入 context")
	}
	if _, ok := startTime.(time.Time); !ok {
		t.Errorf("开始时间类型不正确: %T", startTime)
	}
}

// TestGinLogger_SetContextLogger 测试设置 context logger
func TestGinLogger_SetContextLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	// 验证 logger 存入 context
	contextLogger, exists := c.Get("_logger")
	if !exists {
		t.Error("logger 没有存入 context")
	}
	if contextLogger != logger {
		t.Error("context 中的 logger 与配置的 logger 不一致")
	}
}

// TestGinLogger_Context 测试 context 传递
func TestGinLogger_Context(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// 设置一个有值的 context
	ctx := context.WithValue(context.Background(), "test_key", "test_value")
	c.Request = c.Request.WithContext(ctx)

	handler(c)

	// 验证日志记录成功（不应该因为 context 问题而失败）
	if logger.getInfoCallsCount() == 0 {
		t.Error("GinLogger() 没有记录日志")
	}
}

// TestGinLogger_EmptyTraceIdHeader 测试空 TraceIdHeader 使用默认值
func TestGinLogger_EmptyTraceIdHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := &GinLoggerConfig{
		Logger:         logger,
		TraceIdHeader:  "", // 空值应该使用默认 "X-Trace-Id"
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Trace-Id", "default-header-trace")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证使用默认 header 名读取 traceId
	traceId, ok := lastCall.fields["trace_id"]
	if !ok {
		t.Fatal("日志记录缺少 trace_id 字段")
	}
	if traceIdStr, ok := traceId.(string); ok {
		if traceIdStr != "default-header-trace" {
			t.Errorf("trace_id = %s, 期望 default-header-trace", traceIdStr)
		}
	}
}

// TestGinLogger_ZeroStatus 测试状态码为 0 时的处理
// 注意：gin 的测试上下文默认状态码是 200，所以这个测试主要验证正常流程
// status == 0 的分支是为了代码健壮性保留的，在实际使用中几乎不可能触发
func TestGinLogger_ZeroStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// gin 测试上下文默认状态码是 200
	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证状态码
	status, ok := lastCall.fields["status"]
	if !ok {
		t.Fatal("日志记录缺少 status 字段")
	}
	if statusInt, ok := status.(int); ok {
		if statusInt != 200 {
			t.Errorf("status = %d, 期望 200", statusInt)
		}
	}
}

// TestGinLoggerWithConfig_NilConfigUsesDefault 测试 nil config 使用默认值
func TestGinLoggerWithConfig_NilConfigUsesDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// nil 配置应该创建一个默认配置并正常工作
	handler := GinLoggerWithConfig(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// 不应该 panic
	handler(c)
}

// TestGinLogger_NegativeBodySize 测试 bodySize 为负数时的处理
func TestGinLogger_NegativeBodySize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	handler := GinLogger(logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 body_size 字段存在且非负
	bodySize, ok := lastCall.fields["body_size"]
	if !ok {
		t.Fatal("日志记录缺少 body_size 字段")
	}
	if bodySizeInt, ok := bodySize.(int); ok {
		if bodySizeInt < 0 {
			t.Errorf("body_size = %d, 应该非负", bodySizeInt)
		}
	}
}
