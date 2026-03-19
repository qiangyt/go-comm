# Gin Logger 集成 - 技术设计

## Context

go-comm 项目使用 phuslu/log 作为底层日志库，封装为 `LoggerT` 类型提供结构化日志功能。gin 框架自带 Logger 中间件，但使用独立的日志格式和输出。本设计旨在创建一个 gin 中间件，将 gin 的请求日志统一到 go-comm logger 中，并支持 request body、response body 和 header 的灵活记录。

### 现有架构

- `logger.go`: go-comm logger 实现，基于 phuslu/log
- `gin_sonic.go`: gin JSON 编解码配置
- gin vendor: `vendor/github.com/gin-gonic/gin/logger.go`

### 约束

1. 不能修改 gin vendor 代码
2. 必须与现有 logger.go 兼容
3. 必须遵循项目的 panic-based error handling 规范

## Goals / Non-Goals

**Goals:**
- 创建 gin 中间件，使用 go-comm logger 记录请求日志
- 支持与 gin 默认 Logger 相同的功能（跳过路径、自定义格式）
- 支持 traceId 自动注入和传递
- 结构化日志输出（JSON 格式）
- **支持 request/response body 日志记录（仅文本类型）**
- **支持 body 截取策略配置**
- **支持 header 日志的灵活配置**

**Non-Goals:**
- 不修改 gin 框架本身
- 不替换 gin 的 Recovery 中间件
- 不记录二进制 body 的内容（如文件上传、图片等），但记录类型和大小

## Decisions

### D1: 中间件函数签名

**决策**: 提供函数式配置模式

```go
// 基本用法
func GinLogger(logger Logger) gin.HandlerFunc

// 高级配置
func GinLoggerWithConfig(config GinLoggerConfig) gin.HandlerFunc
```

**理由**:
- 与 gin 框架风格一致（如 `LoggerWithConfig`）
- 简单场景简单使用
- 复杂场景灵活配置

### D2: 配置结构体设计

**决策**: 分离关注点，使用独立配置结构体

```go
// Body 截取策略
type BodyTruncateStrategy int

const (
    BodyTruncateNone BodyTruncateStrategy = iota   // 不记录
    BodyTruncateFull                               // 全量记录
    BodyTruncateHead                               // 只记录前 N 字符
    BodyTruncateTail                               // 只记录后 N 字符
    BodyTruncateHeadAndTail                        // 记录前后各 N 字符（默认值）
)

// Body 日志配置
type BodyLogConfig struct {
    Strategy     BodyTruncateStrategy
    TruncateSize int  // 截取字符数，默认 1024
}
```

**BodyTruncateHeadAndTail 边界情况**:
- 当 body 长度 <= TruncateSize * 2 时，返回完整 body（不做截取，避免重叠）
- 当 body 长度 > TruncateSize * 2 时，返回 `body[:n] + "...(truncated)..." + body[len(body)-n:]`

```go
// Header 日志策略
type HeaderLogStrategy int

const (
    HeaderLogNone HeaderLogStrategy = iota    // 不记录
    HeaderLogAll                              // 记录全部（默认值）
    HeaderLogWhitelist                        // 白名单模式
    HeaderLogBlacklist                        // 黑名单模式
)

// Header 日志配置
type HeaderLogConfig struct {
    Strategy        HeaderLogStrategy
    HeaderList      []string              // 白名单或黑名单
    SensitiveConfig SensitiveHeaderConfig // 敏感 header 处理配置
}

// 主配置
type GinLoggerConfig struct {
    Logger        Logger
    SkipPaths     []string
    CustomFields  func(*gin.Context) map[string]any
    TraceIdHeader string  // 默认 "X-Trace-Id"

    // Body 日志配置
    RequestBody   BodyLogConfig
    ResponseBody  BodyLogConfig

    // Header 日志配置
    RequestHeader  HeaderLogConfig
    ResponseHeader HeaderLogConfig

    // SSE 特殊配置
    SSEConfig     SSELogConfig
}
```

**理由**:
- 配置项复杂，需要结构化设计
- 分离 body 和 header 配置，职责清晰
- 枚举类型保证类型安全

### D3: Body 日志字段设计

**日志字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| method | string | HTTP 方法 |
| path | string | 请求路径 |
| status | int | 状态码 |
| latency | time.Duration | 响应延迟 |
| client_ip | string | 客户端 IP |
| body_size | int | 响应体大小 |
| trace_id | string | 追踪 ID（可选） |
| error | string | 错误信息（可选） |
| **request_body** | string | 请求体（可选，可截取） |
| **response_body** | string | 响应体（可选，可截取） |
| **request_headers** | map[string]string | 请求头（可选） |
| **response_headers** | map[string]string | 响应头（可选） |

### D4: Body 类型判断

**决策**: 文本类型记录内容，二进制类型只记录类型和大小

判断逻辑：
1. 检查 Content-Type 头
2. 文本类型（记录内容）：
   - `application/json`
   - `application/xml` / `text/xml`
   - `text/*`（text/plain, text/html, text/event-stream 等）
   - `application/x-www-form-urlencoded`
3. 二进制类型（只记录类型和大小）：
   - `image/*`
   - `video/*`
   - `audio/*`
   - `application/octet-stream`
   - `multipart/form-data`（文件上传）
   - 其他非文本类型

**二进制 body 日志字段**:
```go
// 二进制 body 输出示例
{
    "request_body_type": "image/png",
    "request_body_size": 102400,  // bytes
    // 不包含 request_body 内容
}
```

**理由**: 避免日志膨胀和安全问题，同时保留类型信息便于调试

### D5: Body 截取实现

**决策**: 使用字符串截取，添加省略标记

```go
// 截取前 N 字符
func truncateHead(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n] + "...(truncated)"
}

// 截取后 N 字符
func truncateTail(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return "...(truncated)" + s[len(s)-n:]
}

// 截取前后各 N 字符（默认策略）
func truncateHeadAndTail(s string, n int) string {
    if len(s) <= n*2 {
        // 避免 overlap，返回完整内容
        return s
    }
    return s[:n] + "...(truncated)..." + s[len(s)-n:]
}
```

### D6: Response Body 捕获

**决策**: 使用自定义 ResponseWriter 包装器

```go
type bodyCaptureWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w *bodyCaptureWriter) Write(b []byte) (int, error) {
    w.body.Write(b)  // 捕获
    return w.ResponseWriter.Write(b)  // 实际写入
}
```

**理由**:
- 需要捕获 response body 进行日志记录
- 不能影响原始响应

### D7: Header 过滤实现

**决策**: 根据策略过滤 header

```go
func filterHeaders(headers http.Header, config HeaderLogConfig) map[string]string {
    switch config.Strategy {
    case HeaderLogNone:
        return nil
    case HeaderLogAll:
        // 返回所有 header
    case HeaderLogWhitelist:
        // 只返回白名单中的 header
    case HeaderLogBlacklist:
        // 返回所有 header，但对黑名单中的敏感 header 进行 mask
    }
}
```

### D7.1: 敏感 Header 处理策略

**决策**: 敏感 header 的记录方式可配置

```go
// 敏感 Header 处理策略
type SensitiveHeaderStrategy int

const (
    SensitiveHeaderFull      SensitiveHeaderStrategy = iota  // 完全记录
    SensitiveHeaderExclude                                   // 不记录
    SensitiveHeaderMaskAll                                   // mask 全部值（替换为 ****）
    SensitiveHeaderMaskHead                                  // mask 前 N 字符
    SensitiveHeaderMaskTail                                  // mask 后 N 字符
)

// 敏感 Header 配置
type SensitiveHeaderConfig struct {
    Strategy      SensitiveHeaderStrategy
    MaskSize      int      // mask 字符数，默认 4
    SensitiveList []string // 敏感 header 列表
}
```

**默认敏感 header 列表**:
```go
var defaultSensitiveHeaders = []string{
    "Authorization",
    "Cookie",
    "Set-Cookie",
    "X-Api-Key",
    "X-Auth-Token",
}
```

**Mask 实现示例**:
```go
// Mask 全部值
func maskAll(value string) string {
    return "****"
}

// Mask 前 N 字符
func maskHead(value string, n int) string {
    if len(value) <= n {
        return "****"
    }
    return "****" + value[n:]
}

// Mask 后 N 字符
func maskTail(value string, n int) string {
    if len(value) <= n {
        return "****"
    }
    return value[:len(value)-n] + "****"
}
```

**示例输出**:
| 策略 | 原值 | 输出 |
|------|------|------|
| Full | `Bearer eyJhbGci...` | `Bearer eyJhbGci...` |
| Exclude | `Bearer eyJhbGci...` | (不记录该 header) |
| MaskAll | `Bearer eyJhbGci...` | `****` |
| MaskHead(4) | `Bearer eyJhbGci...` | `****r eyJhbGci...` |
| MaskTail(4) | `Bearer eyJhbGci...` | `Bearer eyJhbG****` |

**理由**:
- 不同安全级别场景需要不同策略
- 开发环境可能需要完全记录便于调试
- 生产环境可能需要完全排除或严格 mask

### D8: TraceId 处理

**决策**: 支持从请求头读取或自动生成

1. 检查请求头 `X-Trace-Id`（可配置）
2. 如存在，使用该值
3. 如不存在，从 `comm.TraceId` 生成新 ID
4. 将 traceId 存入 gin.Context 供后续使用

### D9: 日志级别映射

**决策**: 按状态码自动选择级别

| 状态码范围 | 日志级别 |
|-----------|---------|
| 200-399 | Info |
| 400-499 | Warn |
| 500+ | Error |

### D10: SSE 流式响应处理

**决策**: SSE 响应支持事件级别的截取策略

```go
// SSE 事件截取策略
type SSETruncateStrategy int

const (
    SSETruncateNone        SSETruncateStrategy = iota  // 不记录
    SSETruncateFull                                    // 全记录
    SSETruncateHead                                    // 记录前 N 条事件
    SSETruncateTail                                    // 记录后 N 条事件
    SSETruncateHeadAndTail                             // 记录前后各 N 条事件（默认值）
)

// SSE 日志配置
type SSELogConfig struct {
    Strategy     SSETruncateStrategy
    TruncateSize int  // 截取事件数，默认 10
}
```

**SSETruncateHeadAndTail 边界情况**:
- 当事件总数 <= TruncateSize * 2 时，返回所有事件（不做截取，避免重叠）
- 当事件总数 > TruncateSize * 2 时，返回前 N 条 + `...(N events truncated)...` + 后 N 条

**SSE 事件解析**:
- SSE 事件以 `\n\n` 或 `\r\n\r\n` 分隔
- 每个事件可能包含 `event:`, `data:`, `id:`, `retry:` 等字段

## Risks / Trade-offs

### R1: 性能影响
- **风险**: body 捕获和日志记录增加处理开销
- **缓解**: 默认关闭 body 日志，按需开启

### R2: 内存使用
- **风险**: 大 body 会占用更多内存
- **缓解**: 提供截取策略，限制捕获大小

### R3: 敏感信息泄露
- **风险**: body 或 header 可能包含敏感信息
- **缓解**:
  - Header 默认使用黑名单过滤敏感 header
  - 提供白名单模式严格控制
  - Body 默认不记录

### R4: SSE 流式响应
- **风险**: SSE 是长连接流式响应，无法在请求结束时记录完整 body
- **缓解**: 支持配置 SSE 事件截取策略，默认记录前后各 10 条事件
