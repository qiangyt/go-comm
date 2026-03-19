# Gin Logger Middleware 规格说明

## ADDED Requirements

### Requirement: 基本请求日志记录

系统 SHALL 提供 gin 中间件，使用 go-comm logger 记录 HTTP 请求信息，包括：
- 请求时间戳
- HTTP 方法和路径
- 状态码
- 响应延迟
- 客户端 IP
- 响应体大小

#### Scenario: 成功请求日志记录
- **WHEN** 一个 GET 请求被处理并返回 200 状态码
- **THEN** 系统使用 go-comm logger 以 Info 级别记录请求信息，包含方法、路径、状态码、延迟、客户端 IP

#### Scenario: 错误请求日志记录
- **WHEN** 一个请求返回 4xx 或 5xx 状态码
- **THEN** 系统使用 go-comm logger 以 Warn（4xx）或 Error（5xx）级别记录请求信息

### Requirement: 跳过特定路径

系统 SHALL 支持配置跳过特定路径的日志记录，避免健康检查等频繁请求产生过多日志。

#### Scenario: 跳过健康检查路径
- **WHEN** 配置了跳过路径 `/health`
- **THEN** 访问 `/health` 时不产生任何日志

#### Scenario: 正常路径记录日志
- **WHEN** 配置了跳过路径 `/health`，但访问 `/api/users`
- **THEN** 系统正常记录 `/api/users` 的请求日志

### Requirement: 自定义日志字段

系统 SHALL 支持自定义日志字段和格式，允许用户指定额外的日志字段。

#### Scenario: 添加自定义字段
- **WHEN** 用户配置了自定义字段函数
- **THEN** 日志输出包含用户指定的额外字段

### Requirement: TraceId 自动注入

系统 SHALL 支持自动注入 traceId 到请求上下文和日志中，便于请求追踪。

#### Scenario: 自动生成 traceId
- **WHEN** 启用 traceId 注入且请求未携带 traceId
- **THEN** 系统自动生成唯一的 traceId 并注入到日志上下文

#### Scenario: 使用请求携带的 traceId
- **WHEN** 请求头中携带了 `X-Trace-Id`
- **THEN** 系统使用该 traceId 而非生成新的

### Requirement: 与 gin 默认 Logger 兼容

系统 SHALL 提供与 gin 默认 Logger 中间件兼容的 API，便于替换使用。

#### Scenario: 替换默认 Logger
- **WHEN** 用户使用 `GinLogger(logger)` 替换 gin 默认的 `gin.Logger()`
- **THEN** 中间件正常工作，但使用 go-comm logger 输出日志

### Requirement: 日志级别配置

系统 SHALL 根据响应状态码自动选择适当的日志级别：
- 2xx/3xx: Info
- 4xx: Warn
- 5xx: Error

#### Scenario: 2xx 响应使用 Info 级别
- **WHEN** 请求返回 200 状态码
- **THEN** 日志级别为 Info

#### Scenario: 4xx 响应使用 Warn 级别
- **WHEN** 请求返回 404 状态码
- **THEN** 日志级别为 Warn

#### Scenario: 5xx 响应使用 Error 级别
- **WHEN** 请求返回 500 状态码
- **THEN** 日志级别为 Error

### Requirement: Request Body 日志记录

系统 SHALL 支持记录 request body，但仅限于文本类型（json、xml、text、form-urlencoded、sse）。

#### Scenario: 记录 JSON request body
- **WHEN** 请求 Content-Type 为 `application/json` 且配置了 request body 日志
- **THEN** 日志中包含 request body 内容

#### Scenario: 记录二进制 request body 的类型和大小
- **WHEN** 请求 Content-Type 为 `image/png`，body 大小为 10240 bytes
- **THEN** 日志中包含 `request_body_type` 为 `image/png`，`request_body_size` 为 10240，但不包含 `request_body` 内容

#### Scenario: 不记录 request body（默认）
- **WHEN** 未配置 request body 日志
- **THEN** 日志中不包含 request body

### Requirement: Response Body 日志记录

系统 SHALL 支持记录 response body，但仅限于文本类型。

#### Scenario: 记录 JSON response body
- **WHEN** 响应 Content-Type 为 `application/json` 且配置了 response body 日志
- **THEN** 日志中包含 response body 内容

#### Scenario: 记录二进制 response body 的类型和大小
- **WHEN** 响应 Content-Type 为 `application/octet-stream`，body 大小为 5120 bytes
- **THEN** 日志中包含 `response_body_type` 为 `application/octet-stream`，`response_body_size` 为 5120，但不包含 `response_body` 内容

### Requirement: Body 截取策略 - 不记录

系统 SHALL 支持 `BodyTruncateNone` 策略，即不记录 body。

#### Scenario: 配置为不记录
- **WHEN** body 截取策略配置为 `BodyTruncateNone`
- **THEN** 日志中不包含该 body

### Requirement: Body 截取策略 - 全量记录

系统 SHALL 支持 `BodyTruncateFull` 策略，即全量记录 body。

#### Scenario: 全量记录短 body
- **WHEN** body 截取策略为 `BodyTruncateFull` 且 body 长度为 100 字符
- **THEN** 日志中包含完整的 100 字符 body

#### Scenario: 全量记录长 body
- **WHEN** body 截取策略为 `BodyTruncateFull` 且 body 长度为 10000 字符
- **THEN** 日志中包含完整的 10000 字符 body

### Requirement: Body 截取策略 - 截取前 N 字符

系统 SHALL 支持 `BodyTruncateHead` 策略，只记录 body 的前 N 个字符。

#### Scenario: body 长度小于截取长度
- **WHEN** body 截取策略为 `BodyTruncateHead`，截取长度 100，body 长度 50
- **THEN** 日志中包含完整的 50 字符 body，无截取标记

#### Scenario: body 长度大于截取长度
- **WHEN** body 截取策略为 `BodyTruncateHead`，截取长度 100，body 长度 200
- **THEN** 日志中包含前 100 字符，并添加 `...(truncated)` 标记

### Requirement: Body 截取策略 - 截取后 N 字符

系统 SHALL 支持 `BodyTruncateTail` 策略，只记录 body 的后 N 个字符。

#### Scenario: body 长度大于截取长度
- **WHEN** body 截取策略为 `BodyTruncateTail`，截取长度 100，body 长度 200
- **THEN** 日志中包含后 100 字符，前面添加 `...(truncated)` 标记

### Requirement: Body 截取策略 - 截取前后各 N 字符

系统 SHALL 支持 `BodyTruncateHeadAndTail` 策略，记录 body 的前 N 和后 N 个字符。

#### Scenario: body 长度大于两倍截取长度
- **WHEN** body 截取策略为 `BodyTruncateHeadAndTail`，截取长度 100，body 长度 500
- **THEN** 日志中包含前 100 字符，中间 `...(truncated)...` 标记，后 100 字符

#### Scenario: body 长度小于两倍截取长度
- **WHEN** body 截取策略为 `BodyTruncateHeadAndTail`，截取长度 100，body 长度 150
- **THEN** 日志中包含完整的 150 字符 body，无截取标记

### Requirement: Request Header 日志记录

系统 SHALL 支持记录 request headers。

#### Scenario: 不记录 request headers
- **WHEN** header 日志策略为 `HeaderLogNone`
- **THEN** 日志中不包含 request headers

#### Scenario: 记录全部 request headers
- **WHEN** header 日志策略为 `HeaderLogAll`
- **THEN** 日志中包含所有 request headers

### Requirement: Header 白名单过滤

系统 SHALL 支持 `HeaderLogWhitelist` 策略，只记录白名单中的 headers。

#### Scenario: 白名单过滤
- **WHEN** header 日志策略为 `HeaderLogWhitelist`，白名单为 `["Content-Type", "Authorization"]`
- **THEN** 日志中只包含 `Content-Type` 和 `Authorization` 两个 header

#### Scenario: 白名单中不存在的 header
- **WHEN** header 日志策略为 `HeaderLogWhitelist`，白名单为 `["X-Custom"]`，但请求中没有该 header
- **THEN** 日志中的 headers 字段为空或不包含该 key

### Requirement: 敏感 Header 处理策略配置

系统 SHALL 支持配置敏感 header 的处理策略，包括：完全记录、不记录、mask 全部值、mask 前 N 字符、mask 后 N 字符。

#### Scenario: 完全记录敏感 header
- **WHEN** 敏感 header 策略为 `SensitiveHeaderFull`，Authorization 值为 `Bearer eyJhbGci...`
- **THEN** 日志中包含完整的 `Authorization` 值

#### Scenario: 不记录敏感 header
- **WHEN** 敏感 header 策略为 `SensitiveHeaderExclude`，Authorization 值为 `Bearer eyJhbGci...`
- **THEN** 日志中不包含 `Authorization` header

#### Scenario: Mask 全部值
- **WHEN** 敏感 header 策略为 `SensitiveHeaderMaskAll`，Authorization 值为 `Bearer eyJhbGci...`
- **THEN** 日志中 `Authorization` 值为 `****`

#### Scenario: Mask 前 N 字符
- **WHEN** 敏感 header 策略为 `SensitiveHeaderMaskHead`，mask 长度 4，Authorization 值为 `Bearer eyJhbGci...`
- **THEN** 日志中 `Authorization` 值为 `****r eyJhbGci...`

#### Scenario: Mask 后 N 字符
- **WHEN** 敏感 header 策略为 `SensitiveHeaderMaskTail`，mask 长度 4，Authorization 值为 `Bearer eyJhbGci...`
- **THEN** 日志中 `Authorization` 值为 `Bearer eyJhbG****`

#### Scenario: 敏感 header 值长度小于 mask 长度
- **WHEN** 敏感 header 策略为 `SensitiveHeaderMaskHead`，mask 长度 10，值长度 5
- **THEN** 日志中该 header 值为 `****`

#### Scenario: 默认敏感 header 列表
- **WHEN** 用户未指定敏感 header 列表
- **THEN** 系统默认对 `Authorization`、`Cookie`、`Set-Cookie`、`X-Api-Key`、`X-Auth-Token` 进行敏感处理

### Requirement: Header 黑名单策略

系统 SHALL 支持 `HeaderLogBlacklist` 策略，记录所有 headers，但根据敏感 header 配置处理指定的 headers。

#### Scenario: 黑名单模式下非敏感 header 保持原样
- **WHEN** header 日志策略为 `HeaderLogBlacklist`，敏感列表为 `["Authorization"]`，Content-Type 值为 `application/json`
- **THEN** 日志中包含 `Content-Type` header，值为 `application/json`（按敏感策略处理）

### Requirement: Response Header 日志记录

系统 SHALL 支持记录 response headers，与 request headers 使用相同的过滤策略。

#### Scenario: 记录 response headers
- **WHEN** 配置了 response header 日志
- **THEN** 日志中包含 response headers，应用配置的过滤策略

### Requirement: SSE 流式响应事件截取

系统 SHALL 对 SSE（Server-Sent Events）类型的 response body 支持事件级别的截取策略。

#### Scenario: SSE 不记录
- **WHEN** SSE 策略为 `SSETruncateNone`
- **THEN** 日志中不包含 SSE 事件

#### Scenario: SSE 全记录
- **WHEN** SSE 策略为 `SSETruncateFull`，有 20 条事件
- **THEN** 日志中包含全部 20 条事件

#### Scenario: SSE 记录前 N 条
- **WHEN** SSE 策略为 `SSETruncateHead`，截取 5 条，有 20 条事件
- **THEN** 日志中包含前 5 条事件 + `...(15 events truncated)`

#### Scenario: SSE 记录后 N 条
- **WHEN** SSE 策略为 `SSETruncateTail`，截取 5 条，有 20 条事件
- **THEN** 日志中包含 `...(15 events truncated)...` + 后 5 条事件

#### Scenario: SSE 记录前后各 N 条（默认）
- **WHEN** SSE 策略为 `SSETruncateHeadAndTail`，截取 5 条，有 20 条事件
- **THEN** 日志中包含前 5 条 + `...(10 events truncated)...` + 后 5 条

#### Scenario: SSE 事件数小于截取数（避免 overlap）
- **WHEN** SSE 策略为 `SSETruncateHeadAndTail`，截取 10 条，只有 15 条事件
- **THEN** 日志中包含全部 15 条事件（不截取，避免重叠）

#### Scenario: SSE 默认截取配置
- **WHEN** 用户未指定 SSE 配置
- **THEN** 系统默认使用 `SSETruncateHeadAndTail` 策略，截取 10 条事件
