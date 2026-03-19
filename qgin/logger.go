package qgin

import (
	"net/http"
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

		// 处理请求
		c.Next()

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
