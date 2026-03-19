package qgin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 状态码到日志级别映射 ====================

// statusToLevel 将 HTTP 状态码映射到日志级别
// 2xx, 3xx -> info
// 4xx -> warn
// 5xx -> error
func statusToLevel(status int) string {
	if status >= 500 {
		return "error"
	}
	if status >= 400 {
		return "warn"
	}
	return "info"
}

// ==================== GinLogger 中间件 ====================

// GinLogger 创建一个使用默认配置的 gin logger 中间件
func GinLogger(logger Logger) gin.HandlerFunc {
	config := DefaultGinLoggerConfig()
	config.Logger = logger
	return GinLoggerWithConfig(config)
}

// GinLoggerWithConfig 使用自定义配置创建 gin logger 中间件
func GinLoggerWithConfig(config *GinLoggerConfig) gin.HandlerFunc {
	// 使用默认配置如果 config 为 nil
	if config == nil {
		config = DefaultGinLoggerConfig()
	}

	// 构建 skipPaths map 用于快速查找
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	// 获取 traceId header 名称，默认 "X-Trace-Id"
	traceIdHeader := config.TraceIdHeader
	if traceIdHeader == "" {
		traceIdHeader = "X-Trace-Id"
	}

	return func(c *gin.Context) {
		// 检查是否应该跳过该路径
		path := c.Request.URL.Path
		if skipPaths[path] {
			c.Next()
			return
		}

		// 记录开始时间
		startTime := time.Now()
		c.Set("_startTime", startTime)

		// 设置 logger 到 context
		if config.Logger != nil {
			c.Set("_logger", config.Logger)
		}

		// 处理 TraceId
		traceId := c.GetHeader(traceIdHeader)
		if traceId == "" {
			// 生成新的 traceId（使用时间戳纳秒）
			traceId = generateTraceId()
		}
		c.Set("trace_id", traceId)

		// === Request Body 读取 ===
		var requestBody string
		var requestBodySize int
		requestContentType := c.Request.Header.Get("Content-Type")
		if config.RequestBody.Strategy != BodyTruncateNone && c.Request.Body != nil {
			requestBody = readRequestBody(c)
			requestBodySize = len(requestBody)
		}

		// === Response Body 捕获 ===
		var captureWriter *bodyCaptureWriter
		if config.ResponseBody.Strategy != BodyTruncateNone {
			captureWriter = newBodyCaptureWriter(c.Writer)
			c.Writer = captureWriter
		}

		// 处理请求
		c.Next()

		// 在 c.Next() 之后读取 response Content-Type（因为 handler 会设置它）
		responseContentType := c.Writer.Header().Get("Content-Type")

		// 计算延迟
		latency := time.Since(startTime)

		// 获取状态码
		status := c.Writer.Status()
		if status == 0 {
			status = http.StatusOK
		}

		// 获取客户端 IP
		clientIP := c.ClientIP()

		// 获取响应体大小
		bodySize := c.Writer.Size()
		if bodySize < 0 {
			bodySize = 0
		}

		// 构建日志字段
		fields := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", latency,
			"client_ip", clientIP,
			"body_size", bodySize,
			"trace_id", traceId,
		}

		// === 添加 Request Body 字段 ===
		if config.RequestBody.Strategy != BodyTruncateNone && requestBody != "" {
			if isTextContentType(requestContentType) {
				truncatedBody := applyTruncateStrategy(requestBody, config.RequestBody)
				if truncatedBody != "" {
					fields = append(fields, "request_body", truncatedBody)
				}
			} else {
				// 二进制类型只记录类型和大小
				fields = append(fields, "request_body_type", requestContentType)
				fields = append(fields, "request_body_size", requestBodySize)
			}
		}

		// === 添加 Response Body 字段 ===
		if config.ResponseBody.Strategy != BodyTruncateNone && captureWriter != nil {
			responseBody := captureWriter.CapturedBody()
			if responseBody != "" {
				// 检查是否为 SSE 响应
				if strings.Contains(responseContentType, "text/event-stream") {
					// SSE 处理
					events := parseSSEEvents(responseBody)
					if len(events) > 0 {
						truncatedSSE := truncateSSEEvents(events, config.SSEConfig)
						if truncatedSSE != "" {
							fields = append(fields, "response_body", truncatedSSE)
						}
					}
				} else if isTextContentType(responseContentType) {
					truncatedBody := applyTruncateStrategy(responseBody, config.ResponseBody)
					if truncatedBody != "" {
						fields = append(fields, "response_body", truncatedBody)
					}
				} else {
					// 二进制类型只记录类型和大小
					fields = append(fields, "response_body_type", responseContentType)
					fields = append(fields, "response_body_size", len(responseBody))
				}
			}
		}

		// === 添加 Request Headers 字段 ===
		if config.RequestHeader.Strategy != HeaderLogNone {
			requestHeaders := filterHeaders(c.Request.Header, config.RequestHeader)
			if len(requestHeaders) > 0 {
				fields = append(fields, "request_headers", requestHeaders)
			}
		}

		// === 添加 Response Headers 字段 ===
		if config.ResponseHeader.Strategy != HeaderLogNone {
			responseHeaders := filterHeaders(c.Writer.Header(), config.ResponseHeader)
			if len(responseHeaders) > 0 {
				fields = append(fields, "response_headers", responseHeaders)
			}
		}

		// === 添加 gin.Errors 字段 ===
		if len(c.Errors) > 0 {
			var errorMessages []string
			for _, err := range c.Errors {
				errorMessages = append(errorMessages, err.Error())
			}
			fields = append(fields, "error", strings.Join(errorMessages, "; "))
		}

		// 添加自定义字段
		if config.CustomFields != nil {
			customFields := config.CustomFields(c)
			for key, value := range customFields {
				fields = append(fields, key, value)
			}
		}

		// 根据状态码选择日志级别并记录日志
		if config.Logger != nil {
			level := statusToLevel(status)
			switch level {
			case "error":
				config.Logger.Error("HTTP Request", fields...)
			case "warn":
				config.Logger.Warn("HTTP Request", fields...)
			default:
				config.Logger.Info("HTTP Request", fields...)
			}
		}
	}
}

// generateTraceId 生成一个唯一的 traceId
func generateTraceId() string {
	return time.Now().Format("20060102150405.999999999")
}
