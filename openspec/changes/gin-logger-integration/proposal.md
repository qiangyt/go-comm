# Gin Logger 集成

## Why

当前 go-comm 项目已有完善的 logger 实现（基于 phuslu/log），且已有 gin 框架的集成（gin_sonic.go 配置 JSON 编解码）。但 gin 框架默认的日志中间件使用自己的日志格式和输出方式，无法与 go-comm 的 logger 统一，导致：

1. 日志格式不统一：gin 请求日志与业务日志格式不同
2. 日志输出分离：gin 日志与业务日志可能输出到不同位置
3. 缺乏统一的日志上下文：无法在请求日志中注入 traceId 等上下文信息
4. 无法复用 logger 的高级特性：如子日志器、结构化日志等
5. **无法记录请求/响应体**：调试和问题排查时缺少关键信息
6. **无法灵活控制 header 记录**：安全和隐私场景需要精细化控制
7. **项目结构不清晰**：相关功能代码分散在根目录，缺乏模块化组织

## What Changes

### 项目结构重构（前置任务）
- 将 gin 相关代码移动到 `qgin/` 目录（gin_sonic.go → qgin/sonic.go）
- 将 config 相关代码移动到 `qconfig/` 目录（config.go → qconfig/config.go）
- 将 io/file/afero 相关代码移动到 `qio/` 目录（afero*.go, file*.go, devnull.go, remote_file.go 等）
- 在根目录创建转发文件保持向后兼容

### 新增 gin logger 中间件
- 新增 gin 日志中间件，使用 go-comm logger 记录 HTTP 请求
- 提供与 gin 默认 Logger 中间件兼容的 API
- 支持自定义日志格式和字段
- 支持跳过特定路径的日志记录
- 支持 traceId 自动注入
- **支持 request body 日志记录**（仅文本类型：json、xml、text、form-urlencoded、sse）
- **支持 response body 日志记录**（仅文本类型，二进制类型只记录类型和大小）
- **支持 body 截取策略**：不记录、全量、截取前N字符、截取后N字符、截取前后各N字符（默认）
- **支持 header 日志记录**：全部记录（默认）、白名单、黑名单、不记录
- **支持敏感 header 处理**：完全记录、不记录、mask全部值、mask前N字符、mask后N字符
- **支持 SSE 事件截取**：不记录、全记录、记录前N条、记录后N条、记录前后各N条（默认）

## Capabilities

### New Capabilities

- `gin-logger-middleware`: 提供 gin 中间件，将 go-comm logger 集成到 gin 框架中，支持请求日志记录、traceId 注入、自定义格式、body 日志、header 日志等功能

### Modified Capabilities

无现有 capabilities 需要修改

## Impact

### go-comm 项目

**结构重构**:
- `gin_sonic.go` → `qgin/sonic.go`
- `config.go` → `qconfig/config.go`
- `afero*.go`, `file*.go`, `devnull.go`, `remote_file.go` → `qio/`

**新增文件**:
- `qgin/logger.go` - gin logger 中间件实现
- `qgin/logger_config.go` - 配置类型定义
- `qgin/logger_truncate.go` - body 截取工具
- `qgin/logger_mask.go` - header mask 工具
- `qgin/logger_test.go` - 中间件测试

**其他**:
- 依赖现有 `logger.go` 和 gin 框架
- 发布 tag v2.7.2

### mobile-claude 项目

- **依赖更新**: go-comm 升级到 v2.7.2
- **import 更新**: 更新 import 路径（`go-comm` → `go-comm/qgin`, `go-comm/qconfig`, `go-comm/qio`）
- **变更**: 替换原有的 gin logger，使用 go-comm/qgin 的新 gin logger 中间件
