package qgin

// ==================== Body 截取策略 ====================

// BodyTruncateStrategy 定义 body 截取策略
type BodyTruncateStrategy int

const (
	// BodyTruncateNone 不记录 body
	BodyTruncateNone BodyTruncateStrategy = iota
	// BodyTruncateFull 全量记录
	BodyTruncateFull
	// BodyTruncateHead 只记录前 N 字符
	BodyTruncateHead
	// BodyTruncateTail 只记录后 N 字符
	BodyTruncateTail
	// BodyTruncateHeadAndTail 记录前后各 N 字符（默认值）
	BodyTruncateHeadAndTail
)

// BodyLogConfig body 日志配置
type BodyLogConfig struct {
	// Strategy 截取策略
	Strategy BodyTruncateStrategy
	// TruncateSize 截取字符数，默认 1024
	TruncateSize int
}

// DefaultBodyLogConfig 返回默认的 body 日志配置
func DefaultBodyLogConfig() BodyLogConfig {
	return BodyLogConfig{
		Strategy:     BodyTruncateHeadAndTail,
		TruncateSize: 1024,
	}
}

// ==================== Header 日志策略 ====================

// HeaderLogStrategy 定义 header 日志策略
type HeaderLogStrategy int

const (
	// HeaderLogNone 不记录 header
	HeaderLogNone HeaderLogStrategy = iota
	// HeaderLogAll 记录全部（默认值）
	HeaderLogAll
	// HeaderLogWhitelist 白名单模式
	HeaderLogWhitelist
	// HeaderLogBlacklist 黑名单模式
	HeaderLogBlacklist
)

// ==================== 敏感 Header 处理策略 ====================

// SensitiveHeaderStrategy 敏感 header 处理策略
type SensitiveHeaderStrategy int

const (
	// SensitiveHeaderFull 完全记录
	SensitiveHeaderFull SensitiveHeaderStrategy = iota
	// SensitiveHeaderExclude 不记录
	SensitiveHeaderExclude
	// SensitiveHeaderMaskAll mask 全部值（替换为 ****）
	SensitiveHeaderMaskAll
	// SensitiveHeaderMaskHead mask 前 N 字符
	SensitiveHeaderMaskHead
	// SensitiveHeaderMaskTail mask 后 N 字符
	SensitiveHeaderMaskTail
)

// SensitiveHeaderConfig 敏感 header 配置
type SensitiveHeaderConfig struct {
	// Strategy 处理策略
	Strategy SensitiveHeaderStrategy
	// MaskSize mask 字符数，默认 4
	MaskSize int
	// SensitiveList 敏感 header 列表
	SensitiveList []string
}

// DefaultSensitiveHeaderConfig 返回默认的敏感 header 配置
func DefaultSensitiveHeaderConfig() *SensitiveHeaderConfig {
	return &SensitiveHeaderConfig{
		Strategy: SensitiveHeaderMaskAll,
		MaskSize: 4,
		SensitiveList: []string{
			"Authorization",
			"Cookie",
			"Set-Cookie",
			"X-Api-Key",
			"X-Auth-Token",
		},
	}
}

// HeaderLogConfig header 日志配置
type HeaderLogConfig struct {
	// Strategy 日志策略
	Strategy HeaderLogStrategy
	// HeaderList 白名单或黑名单
	HeaderList []string
	// SensitiveConfig 敏感 header 处理配置
	SensitiveConfig *SensitiveHeaderConfig
}

// DefaultHeaderLogConfig 返回默认的 header 日志配置
func DefaultHeaderLogConfig() HeaderLogConfig {
	return HeaderLogConfig{
		Strategy:       HeaderLogAll,
		HeaderList:     nil,
		SensitiveConfig: DefaultSensitiveHeaderConfig(),
	}
}

// ==================== SSE 事件截取策略 ====================

// SSETruncateStrategy SSE 事件截取策略
type SSETruncateStrategy int

const (
	// SSETruncateNone 不记录
	SSETruncateNone SSETruncateStrategy = iota
	// SSETruncateFull 全记录
	SSETruncateFull
	// SSETruncateHead 记录前 N 条事件
	SSETruncateHead
	// SSETruncateTail 记录后 N 条事件
	SSETruncateTail
	// SSETruncateHeadAndTail 记录前后各 N 条事件（默认值）
	SSETruncateHeadAndTail
)

// SSELogConfig SSE 日志配置
type SSELogConfig struct {
	// Strategy 截取策略
	Strategy SSETruncateStrategy
	// TruncateSize 截取事件数，默认 10
	TruncateSize int
}

// DefaultSSELogConfig 返回默认的 SSE 日志配置
func DefaultSSELogConfig() SSELogConfig {
	return SSELogConfig{
		Strategy:     SSETruncateHeadAndTail,
		TruncateSize: 10,
	}
}

// ==================== 主配置 ====================

// Logger 定义日志接口
type Logger interface {
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Trace(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Panic(msg string, fields ...any)
	WithField(key string, value any) any
	WithFields(fields map[string]any) any
}

// GinLoggerConfig gin logger 主配置
type GinLoggerConfig struct {
	// Logger 日志器实例
	Logger Logger
	// SkipPaths 跳过日志记录的路径
	SkipPaths []string
	// CustomFields 自定义字段回调
	CustomFields func(ctx any) map[string]any
	// TraceIdHeader traceId 请求头名称，默认 "X-Trace-Id"
	TraceIdHeader string

	// RequestBody 请求 body 日志配置
	RequestBody BodyLogConfig
	// ResponseBody 响应 body 日志配置
	ResponseBody BodyLogConfig

	// RequestHeader 请求 header 日志配置
	RequestHeader HeaderLogConfig
	// ResponseHeader 响应 header 日志配置
	ResponseHeader HeaderLogConfig

	// SSEConfig SSE 特殊配置
	SSEConfig SSELogConfig
}

// DefaultGinLoggerConfig 返回默认的 gin logger 配置
func DefaultGinLoggerConfig() *GinLoggerConfig {
	return &GinLoggerConfig{
		Logger:         nil,
		SkipPaths:      nil,
		TraceIdHeader:  "X-Trace-Id",
		CustomFields:   nil,
		RequestBody:    DefaultBodyLogConfig(),
		ResponseBody:   DefaultBodyLogConfig(),
		RequestHeader:  DefaultHeaderLogConfig(),
		ResponseHeader: DefaultHeaderLogConfig(),
		SSEConfig:      DefaultSSELogConfig(),
	}
}
