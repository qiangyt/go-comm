package qgin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

// ==================== Body 日志功能集成测试 ====================

// TestGinLogger_RequestBody_LogsTextBody 测试 request body 日志记录（文本类型）
func TestGinLogger_RequestBody_LogsTextBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestBody = BodyLogConfig{
		Strategy:     BodyTruncateHead,
		TruncateSize: 100,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"test","value":"sample"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 request_body 字段存在
	requestBody, ok := lastCall.fields["request_body"]
	if !ok {
		t.Error("日志记录缺少 request_body 字段")
	}
	if requestBodyStr, ok := requestBody.(string); ok {
		if requestBodyStr == "" {
			t.Error("request_body 不应该为空")
		}
	}
}

// TestGinLogger_RequestBody_Truncation 测试 request body 截取功能
func TestGinLogger_RequestBody_Truncation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestBody = BodyLogConfig{
		Strategy:     BodyTruncateHead,
		TruncateSize: 10,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// 创建一个长 body
	longBody := strings.Repeat("abcdefghij", 10) // 100 字符
	c.Request = httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(longBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	requestBody, ok := lastCall.fields["request_body"]
	if !ok {
		t.Fatal("日志记录缺少 request_body 字段")
	}
	if requestBodyStr, ok := requestBody.(string); ok {
		// 截取后应该包含 truncated 标记
		if !strings.Contains(requestBodyStr, "...(truncated)") {
			t.Errorf("request_body 应该被截取，实际值: %s", requestBodyStr)
		}
	}
}

// TestGinLogger_ResponseBody_LogsTextBody 测试 response body 日志记录（文本类型）
func TestGinLogger_ResponseBody_LogsTextBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.ResponseBody = BodyLogConfig{
		Strategy:     BodyTruncateFull,
		TruncateSize: 100,
	}

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/test", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.String(200, `{"status":"ok"}`)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 response_body 字段存在
	responseBody, ok := lastCall.fields["response_body"]
	if !ok {
		t.Error("日志记录缺少 response_body 字段")
	}
	if responseBodyStr, ok := responseBody.(string); ok {
		if responseBodyStr == "" {
			t.Error("response_body 不应该为空")
		}
	}
}

// TestGinLogger_Body_BinaryContentType 测试二进制 body 只记录类型和大小
func TestGinLogger_Body_BinaryContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestBody = BodyLogConfig{
		Strategy:     BodyTruncateFull,
		TruncateSize: 100,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	c.Request = httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(binaryData))
	c.Request.Header.Set("Content-Type", "image/png")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 对于二进制类型，应该记录 request_body_type 和 request_body_size
	bodyType, hasType := lastCall.fields["request_body_type"]
	bodySize, hasSize := lastCall.fields["request_body_size"]

	if hasType {
		if bodyTypeStr, ok := bodyType.(string); ok {
			if !strings.Contains(bodyTypeStr, "image/png") {
				t.Errorf("request_body_type = %s, 期望包含 image/png", bodyTypeStr)
			}
		}
	}

	if hasSize {
		if size, ok := bodySize.(int); ok {
			if size != len(binaryData) {
				t.Errorf("request_body_size = %d, 期望 %d", size, len(binaryData))
			}
		}
	}

	// 二进制内容不应该有 request_body 内容
	if _, hasBody := lastCall.fields["request_body"]; hasBody {
		body := lastCall.fields["request_body"]
		if bodyStr, ok := body.(string); ok && bodyStr != "" {
			t.Error("二进制 body 不应该记录 request_body 内容")
		}
	}
}

// TestGinLogger_Body_NoneStrategy 测试 None 策略不记录 body
func TestGinLogger_Body_NoneStrategy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestBody = BodyLogConfig{
		Strategy: BodyTruncateNone,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"test":"data"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// None 策略不应该记录 request_body
	if _, hasBody := lastCall.fields["request_body"]; hasBody {
		t.Error("BodyTruncateNone 策略不应该记录 request_body")
	}
}

// ==================== Header 日志功能集成测试 ====================

// ==================== 覆盖率补充测试 ====================

// TestGinLogger_StatusZero 测试 status == 0 时的处理
func TestGinLogger_StatusZero(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger

	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/no-status", func(c *gin.Context) {
		// 不设置任何状态码，gin 默认会返回 0
		// 但实际写入响应后会变成 200
		c.Writer.Write([]byte("ok"))
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/no-status", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证状态码为 200（因为 0 会被转换为 200）
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

// TestGinLogger_ResponseBody_BinaryType 测试 response body 二进制类型记录
func TestGinLogger_ResponseBody_BinaryType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.ResponseBody = BodyLogConfig{
		Strategy:     BodyTruncateFull,
		TruncateSize: 100,
	}

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/binary", func(c *gin.Context) {
		c.Header("Content-Type", "image/png")
		c.Data(200, "image/png", []byte{0x89, 0x50, 0x4E, 0x47})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/binary", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 二进制类型应该记录 response_body_type 和 response_body_size
	bodyType, hasType := lastCall.fields["response_body_type"]
	if !hasType {
		t.Error("日志记录缺少 response_body_type 字段")
	}
	if bodyTypeStr, ok := bodyType.(string); ok {
		if !strings.Contains(bodyTypeStr, "image/png") {
			t.Errorf("response_body_type = %s, 期望包含 image/png", bodyTypeStr)
		}
	}

	_, hasSize := lastCall.fields["response_body_size"]
	if !hasSize {
		t.Error("日志记录缺少 response_body_size 字段")
	}
}

// ==================== SSE 日志功能集成测试 ====================

// TestGinLogger_SSE_EventsTruncation 测试 SSE 响应的事件截取
func TestGinLogger_SSE_EventsTruncation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.ResponseBody = BodyLogConfig{
		Strategy:     BodyTruncateFull,
		TruncateSize: 100,
	}
	config.SSEConfig = SSELogConfig{
		Strategy:     SSETruncateHeadAndTail,
		TruncateSize: 2, // 截取前后各 2 条事件
	}

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/sse", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		// 写入 10 条 SSE 事件
		for i := 0; i < 10; i++ {
			c.String(200, "data: event %d\n\n", i)
		}
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 response_body 字段存在
	responseBody, ok := lastCall.fields["response_body"]
	if !ok {
		t.Fatal("日志记录缺少 response_body 字段")
	}
	if responseBodyStr, ok := responseBody.(string); ok {
		// 应该包含截取标记
		if !strings.Contains(responseBodyStr, "...(truncated)...") {
			t.Errorf("SSE 响应应该被截取，实际值: %s", responseBodyStr)
		}
	}
}

// TestGinLogger_SSE_FullStrategy 测试 SSE Full 策略记录所有事件
func TestGinLogger_SSE_FullStrategy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.ResponseBody = BodyLogConfig{
		Strategy: BodyTruncateFull,
	}
	config.SSEConfig = SSELogConfig{
		Strategy: SSETruncateFull,
	}

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/sse", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.String(200, "data: event1\n\ndata: event2\n\n")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	responseBody, ok := lastCall.fields["response_body"]
	if !ok {
		t.Fatal("日志记录缺少 response_body 字段")
	}
	if responseBodyStr, ok := responseBody.(string); ok {
		// Full 策略不应该有截取标记
		if strings.Contains(responseBodyStr, "...(truncated)") {
			t.Errorf("SSE Full 策略不应该截取，实际值: %s", responseBodyStr)
		}
	}
}

// ==================== 错误处理测试 ====================

// TestGinLogger_Errors_ExtractGinErrors 测试 gin.Errors 中的错误信息提取
func TestGinLogger_Errors_ExtractGinErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/error", func(c *gin.Context) {
		c.Error(fmt.Errorf("test error 1"))
		c.Error(fmt.Errorf("test error 2"))
		c.String(500, "internal error")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	r.ServeHTTP(w, req)

	// 500 状态码会触发 Error 级别日志
	lastCall := logger.getLastErrorCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 error 字段存在
	errField, ok := lastCall.fields["error"]
	if !ok {
		t.Error("日志记录缺少 error 字段")
	}
	if errStr, ok := errField.(string); ok {
		if !strings.Contains(errStr, "test error 1") {
			t.Errorf("error 字段应该包含 'test error 1'，实际值: %s", errStr)
		}
	}
}

// TestGinLogger_Errors_NoErrors 测试没有错误时不记录 error 字段
func TestGinLogger_Errors_NoErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/ok", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 没有 error 字段
	if _, hasError := lastCall.fields["error"]; hasError {
		t.Error("没有错误时不应该记录 error 字段")
	}
}
func TestGinLogger_RequestHeader_LogsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestHeader = HeaderLogConfig{
		Strategy: HeaderLogAll,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Custom-Header", "custom-value")
	c.Request.Header.Set("Accept", "application/json")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 request_headers 字段存在
	requestHeaders, ok := lastCall.fields["request_headers"]
	if !ok {
		t.Error("日志记录缺少 request_headers 字段")
	}
	if headers, ok := requestHeaders.(map[string]string); ok {
		if headers["X-Custom-Header"] != "custom-value" {
			t.Errorf("X-Custom-Header = %s, 期望 custom-value", headers["X-Custom-Header"])
		}
	}
}

// TestGinLogger_RequestHeader_Whitelist 测试 request header 白名单模式
func TestGinLogger_RequestHeader_Whitelist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestHeader = HeaderLogConfig{
		Strategy:   HeaderLogWhitelist,
		HeaderList: []string{"X-Custom-Header"},
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Custom-Header", "custom-value")
	c.Request.Header.Set("Accept", "application/json")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	requestHeaders, ok := lastCall.fields["request_headers"]
	if !ok {
		t.Fatal("日志记录缺少 request_headers 字段")
	}
	headers, ok := requestHeaders.(map[string]string)
	if !ok {
		t.Fatal("request_headers 类型不正确")
	}

	// 白名单中的 header 应该存在
	if headers["X-Custom-Header"] != "custom-value" {
		t.Errorf("X-Custom-Header = %s, 期望 custom-value", headers["X-Custom-Header"])
	}

	// 不在白名单中的 header 不应该存在
	if _, exists := headers["Accept"]; exists {
		t.Error("Accept 不在白名单中，不应该被记录")
	}
}

// TestGinLogger_RequestHeader_Sensitive 测试敏感 header 处理
func TestGinLogger_RequestHeader_Sensitive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestHeader = HeaderLogConfig{
		Strategy: HeaderLogAll,
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderMaskAll,
			SensitiveList: []string{"Authorization"},
		},
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("Authorization", "Bearer secret-token")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	requestHeaders, ok := lastCall.fields["request_headers"]
	if !ok {
		t.Fatal("日志记录缺少 request_headers 字段")
	}
	headers, ok := requestHeaders.(map[string]string)
	if !ok {
		t.Fatal("request_headers 类型不正确")
	}

	// 敏感 header 应该被 mask
	if headers["Authorization"] != "****" {
		t.Errorf("Authorization = %s, 期望 **** (被 mask)", headers["Authorization"])
	}
}

// TestGinLogger_ResponseHeader_LogsHeaders 测试 response header 日志记录
func TestGinLogger_ResponseHeader_LogsHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.ResponseHeader = HeaderLogConfig{
		Strategy: HeaderLogAll,
	}

	// 创建一个完整的 gin 引擎来测试
	r := gin.New()
	r.Use(GinLoggerWithConfig(config))
	r.GET("/test", func(c *gin.Context) {
		c.Header("X-Response-Header", "response-value")
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// 验证 response_headers 字段存在
	responseHeaders, ok := lastCall.fields["response_headers"]
	if !ok {
		t.Error("日志记录缺少 response_headers 字段")
	}
	if headers, ok := responseHeaders.(map[string]string); ok {
		if headers["X-Response-Header"] != "response-value" {
			t.Errorf("X-Response-Header = %s, 期望 response-value", headers["X-Response-Header"])
		}
	}
}

// TestGinLogger_Header_NoneStrategy 测试 None 策略不记录 headers
func TestGinLogger_Header_NoneStrategy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := newMockLogger()
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	config.RequestHeader = HeaderLogConfig{
		Strategy: HeaderLogNone,
	}
	handler := GinLoggerWithConfig(config)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set("X-Custom-Header", "custom-value")

	handler(c)

	lastCall := logger.getLastInfoCall()
	if lastCall == nil {
		t.Fatal("GinLogger() 没有记录日志")
	}

	// None 策略不应该记录 request_headers
	if _, hasHeaders := lastCall.fields["request_headers"]; hasHeaders {
		t.Error("HeaderLogNone 策略不应该记录 request_headers")
	}
}
