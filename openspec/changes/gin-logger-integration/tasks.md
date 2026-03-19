# Gin Logger 集成 - 任务清单

> **TDD 流程说明**：
> 每个阶段都严格遵循 Red-Green-Refactor 循环：
> 1. **Phase 前回归**：运行所有测试，确保从干净状态开始
> 2. **Red**：编写失败的测试
> 3. **Green**：编写最小实现使测试通过
> 4. **Refactor**：在测试保护下重构代码
> 5. **Phase 后回归**：运行所有测试，确保没有破坏现有功能
> 6. **修复非本次问题**：即使不是本次修改导致的测试失败，也要修复

---

## Phase 0: 项目结构重构（前置任务）

> **目的**：将现有代码按功能模块重新组织，提高代码可维护性

### Step 0.0: Phase 前回归
- [ ] 0.0.1 运行 `go test ./...` 确保所有现有测试通过

### Step 0.1: 移动 gin 相关代码到 qgin/ 目录
**涉及文件**:
- `gin_sonic.go` → `qgin/sonic.go`
- `gin_sonic_test.go` → `qgin/sonic_test.go`

**Red:**
- [ ] 0.1.1 编写测试：验证 `github.com/qiangyt/go-comm/qgin` 包可导入且功能正常
- [ ] 0.1.2 运行测试确认失败（包不存在）
**Green:**
- [ ] 0.1.3 创建 `qgin/` 目录（已存在）
- [ ] 0.1.4 移动 `gin_sonic.go` 到 `qgin/sonic.go`，修改 package 声明
- [ ] 0.1.5 移动 `gin_sonic_test.go` 到 `qgin/sonic_test.go`
- [ ] 0.1.6 在 go-comm 根目录创建转发文件（保持向后兼容）
- [ ] 0.1.7 运行测试确认通过
**Refactor:**
- [ ] 0.1.8 检查并清理代码

### Step 0.2: 移动 config 相关代码到 qconfig/ 目录
**涉及文件**:
- `config.go` → `qconfig/config.go`
- `config_main_test.go` → `qconfig/config_main_test.go`

**Red:**
- [ ] 0.2.1 编写测试：验证 `github.com/qiangyt/go-comm/qconfig` 包可导入且功能正常
- [ ] 0.2.2 运行测试确认失败（包不存在）
**Green:**
- [ ] 0.2.3 创建 `qconfig/` 目录（已存在）
- [ ] 0.2.4 移动 `config.go` 到 `qconfig/config.go`，修改 package 声明
- [ ] 0.2.5 移动 `config_main_test.go` 到 `qconfig/config_main_test.go`
- [ ] 0.2.6 在 go-comm 根目录创建转发文件（保持向后兼容）
- [ ] 0.2.7 运行测试确认通过
**Refactor:**
- [ ] 0.2.8 检查并清理代码

### Step 0.3: 移动 io/file/afero 相关代码到 qio/ 目录
**涉及文件**:
- `afero*.go` → `qio/afero*.go`（afero.go, afero_darwin.go, afero_linux.go, afero_windows.go）
- `afero*_test.go` → `qio/afero*_test.go`（afero_test.go, afero_ext_test.go, afero_file_test.go）
- `afero_file.go` → `qio/afero_file.go`
- `file.go` → `qio/file.go`
- `file_main_test.go` → `qio/file_main_test.go`
- `file_fallback.go` → `qio/file_fallback.go`
- `file_fallback_main_test.go` → `qio/file_fallback_main_test.go`
- `file_cache.go` → `qio/file_cache.go`
- `file_cache_main_test.go` → `qio/file_cache_main_test.go`
- `fileops.go` → `qio/fileops.go`
- `fileops_main_test.go` → `qio/fileops_main_test.go`
- `devnull.go` → `qio/devnull.go`
- `devnull_test.go` → `qio/devnull_test.go`
- `remote_file.go` → `qio/remote_file.go`
- `remote_file_main_test.go` → `qio/remote_file_main_test.go`

**Red:**
- [ ] 0.3.1 编写测试：验证 `github.com/qiangyt/go-comm/qio` 包可导入且功能正常
- [ ] 0.3.2 运行测试确认失败（包不存在）
**Green:**
- [ ] 0.3.3 创建 `qio/` 目录（已存在）
- [ ] 0.3.4 移动所有 afero 相关文件到 `qio/`，修改 package 声明
- [ ] 0.3.5 移动所有 file 相关文件到 `qio/`，修改 package 声明
- [ ] 0.3.6 移动 devnull 相关文件到 `qio/`，修改 package 声明
- [ ] 0.3.7 移动 remote_file 相关文件到 `qio/`，修改 package 声明
- [ ] 0.3.8 在 go-comm 根目录创建转发文件（保持向后兼容）
- [ ] 0.3.9 运行测试确认通过
**Refactor:**
- [ ] 0.3.10 检查并清理代码

### Step 0.4: 更新 go-comm 内部引用
**Red:**
- [ ] 0.4.1 编写测试：验证所有原有功能仍然正常工作
- [ ] 0.4.2 运行测试确认失败（如有引用问题）
**Green:**
- [ ] 0.4.3 更新 go-comm 内部所有对移动文件的 import 引用
- [ ] 0.4.4 运行测试确认通过
**Refactor:**
- [ ] 0.4.5 检查并清理 import 语句

### Step 0.5: Phase 后回归
- [ ] 0.5.1 运行 `go test ./...` 确保所有测试通过
- [ ] 0.5.2 确保行覆盖率和分支覆盖率均为 100%
- [ ] 0.5.3 运行 `go vet ./...` 检查代码质量

---

## Phase 1: 基础结构和类型定义

> **代码位置**：所有 gin logger 相关代码放在 `qgin/` 目录下
> - 类型定义：`qgin/logger_config.go`
> - 工具函数：`qgin/logger_truncate.go`, `qgin/logger_mask.go`
> - 中间件：`qgin/logger.go`
> - 测试文件：`qgin/logger_*_test.go`

### Step 1.0: Phase 前回归
- [ ] 1.0.1 运行 `go test ./...` 确保所有现有测试通过

### Step 1.1: BodyTruncateStrategy 枚举
**Red:**
- [ ] 1.1.1 编写测试：验证 BodyTruncateStrategy 所有枚举值
- [ ] 1.1.2 运行测试确认失败
**Green:**
- [ ] 1.1.3 实现 BodyTruncateStrategy 枚举（None, Full, Head, Tail, HeadAndTail）
- [ ] 1.1.4 运行测试确认通过
**Refactor:**
- [ ] 1.1.5 检查并优化代码结构

### Step 1.2: BodyLogConfig 结构体
**Red:**
- [ ] 1.2.1 编写测试：验证 BodyLogConfig 字段和默认值
- [ ] 1.2.2 运行测试确认失败
**Green:**
- [ ] 1.2.3 实现 BodyLogConfig 结构体（Strategy, TruncateSize 默认 1024）
- [ ] 1.2.4 运行测试确认通过
**Refactor:**
- [ ] 1.2.5 检查并优化代码结构

### Step 1.3: HeaderLogStrategy 枚举
**Red:**
- [ ] 1.3.1 编写测试：验证 HeaderLogStrategy 所有枚举值
- [ ] 1.3.2 运行测试确认失败
**Green:**
- [ ] 1.3.3 实现 HeaderLogStrategy 枚举（None, All, Whitelist, Blacklist），All 为默认值
- [ ] 1.3.4 运行测试确认通过
**Refactor:**
- [ ] 1.3.5 检查并优化代码结构

### Step 1.4: SensitiveHeaderStrategy 枚举和配置
**Red:**
- [ ] 1.4.1 编写测试：验证 SensitiveHeaderStrategy 所有枚举值
- [ ] 1.4.2 编写测试：验证 SensitiveHeaderConfig 字段和默认值
- [ ] 1.4.3 运行测试确认失败
**Green:**
- [ ] 1.4.4 实现 SensitiveHeaderStrategy 枚举（Full, Exclude, MaskAll, MaskHead, MaskTail）
- [ ] 1.4.5 实现 SensitiveHeaderConfig 结构体（Strategy, MaskSize 默认 4, SensitiveList）
- [ ] 1.4.6 运行测试确认通过
**Refactor:**
- [ ] 1.4.7 检查并优化代码结构

### Step 1.5: HeaderLogConfig 结构体
**Red:**
- [ ] 1.5.1 编写测试：验证 HeaderLogConfig 字段（Strategy, HeaderList, SensitiveConfig）
- [ ] 1.5.2 运行测试确认失败
**Green:**
- [ ] 1.5.3 实现 HeaderLogConfig 结构体
- [ ] 1.5.4 运行测试确认通过
**Refactor:**
- [ ] 1.5.5 检查并优化代码结构

### Step 1.6: SSELogConfig 结构体
**Red:**
- [ ] 1.6.1 编写测试：验证 SSELogConfig 字段（Strategy, TruncateSize 默认 10）
- [ ] 1.6.2 运行测试确认失败
**Green:**
- [ ] 1.6.3 实现 SSETruncateStrategy 枚举（None, Full, Head, Tail, HeadAndTail）
- [ ] 1.6.4 实现 SSELogConfig 结构体
- [ ] 1.6.5 运行测试确认通过
**Refactor:**
- [ ] 1.6.6 检查并优化代码结构

### Step 1.7: GinLoggerConfig 主配置
**Red:**
- [ ] 1.7.1 编写测试：验证 GinLoggerConfig 所有字段和默认值
- [ ] 1.7.2 运行测试确认失败
**Green:**
- [ ] 1.7.3 实现 GinLoggerConfig 主配置结构体
- [ ] 1.7.4 运行测试确认通过
**Refactor:**
- [ ] 1.7.5 检查并优化代码结构

### Step 1.8: Phase 后回归
- [ ] 1.8.1 运行 `go test ./...` 确保所有测试通过
- [ ] 1.8.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 2: Body 截取工具函数

### Step 2.0: Phase 前回归
- [ ] 2.0.1 运行 `go test ./...` 确保所有测试通过

### Step 2.1: truncateHead 函数
**Red:**
- [ ] 2.1.1 编写测试：截取前 N 字符的各种场景（短字符串、长字符串、边界情况）
- [ ] 2.1.2 运行测试确认失败
**Green:**
- [ ] 2.1.3 实现 truncateHead(s string, n int) string
- [ ] 2.1.4 运行测试确认通过
**Refactor:**
- [ ] 2.1.5 检查并优化代码结构

### Step 2.2: truncateTail 函数
**Red:**
- [ ] 2.2.1 编写测试：截取后 N 字符的各种场景
- [ ] 2.2.2 运行测试确认失败
**Green:**
- [ ] 2.2.3 实现 truncateTail(s string, n int) string
- [ ] 2.2.4 运行测试确认通过
**Refactor:**
- [ ] 2.2.5 检查并优化代码结构

### Step 2.3: truncateHeadAndTail 函数
**Red:**
- [ ] 2.3.1 编写测试：截取前后各 N 字符的各种场景，**特别注意 overlap 边界情况**（len <= n*2 时不截取）
- [ ] 2.3.2 运行测试确认失败
**Green:**
- [ ] 2.3.3 实现 truncateHeadAndTail(s string, n int) string，正确处理 overlap
- [ ] 2.3.4 运行测试确认通过
**Refactor:**
- [ ] 2.3.5 检查并优化代码结构

### Step 2.4: applyTruncateStrategy 统一入口
**Red:**
- [ ] 2.4.1 编写测试：验证所有截取策略的统一入口
- [ ] 2.4.2 运行测试确认失败
**Green:**
- [ ] 2.4.3 实现 applyTruncateStrategy(s string, config BodyLogConfig) string
- [ ] 2.4.4 运行测试确认通过
**Refactor:**
- [ ] 2.4.5 检查并优化代码结构

### Step 2.5: Phase 后回归
- [ ] 2.5.1 运行 `go test ./...` 确保所有测试通过
- [ ] 2.5.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 3: Body 类型判断

### Step 3.0: Phase 前回归
- [ ] 3.0.1 运行 `go test ./...` 确保所有测试通过

### Step 3.1: isTextContentType 函数
**Red:**
- [ ] 3.1.1 编写测试：验证各种 Content-Type 的判断（json, xml, text/*, form-urlencoded, sse, 二进制类型）
- [ ] 3.1.2 运行测试确认失败
**Green:**
- [ ] 3.1.3 实现 isTextContentType(contentType string) bool
- [ ] 3.1.4 运行测试确认通过
**Refactor:**
- [ ] 3.1.5 检查并优化代码结构

### Step 3.2: Phase 后回归
- [ ] 3.2.1 运行 `go test ./...` 确保所有测试通过
- [ ] 3.2.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 4: Header 过滤和敏感 Header 处理

### Step 4.0: Phase 前回归
- [ ] 4.0.1 运行 `go test ./...` 确保所有测试通过

### Step 4.1: maskAll 函数
**Red:**
- [ ] 4.1.1 编写测试：验证 mask 全部值为 ****
- [ ] 4.1.2 运行测试确认失败
**Green:**
- [ ] 4.1.3 实现 maskAll(value string) string
- [ ] 4.1.4 运行测试确认通过
**Refactor:**
- [ ] 4.1.5 检查并优化代码结构

### Step 4.2: maskHead 函数
**Red:**
- [ ] 4.2.1 编写测试：验证 mask 前 N 字符（含短值边界）
- [ ] 4.2.2 运行测试确认失败
**Green:**
- [ ] 4.2.3 实现 maskHead(value string, n int) string
- [ ] 4.2.4 运行测试确认通过
**Refactor:**
- [ ] 4.2.5 检查并优化代码结构

### Step 4.3: maskTail 函数
**Red:**
- [ ] 4.3.1 编写测试：验证 mask 后 N 字符（含短值边界）
- [ ] 4.3.2 运行测试确认失败
**Green:**
- [ ] 4.3.3 实现 maskTail(value string, n int) string
- [ ] 4.3.4 运行测试确认通过
**Refactor:**
- [ ] 4.3.5 检查并优化代码结构

### Step 4.4: applySensitiveStrategy 统一入口
**Red:**
- [ ] 4.4.1 编写测试：验证所有敏感策略的统一入口
- [ ] 4.4.2 运行测试确认失败
**Green:**
- [ ] 4.4.3 实现 applySensitiveStrategy(value string, config SensitiveHeaderConfig) string
- [ ] 4.4.4 运行测试确认通过
**Refactor:**
- [ ] 4.4.5 检查并优化代码结构

### Step 4.5: filterHeaders 系列函数
**Red:**
- [ ] 4.5.1 编写测试：验证 filterHeadersNone
- [ ] 4.5.2 编写测试：验证 filterHeadersAll（含敏感处理）
- [ ] 4.5.3 编写测试：验证 filterHeadersWhitelist
- [ ] 4.5.4 编写测试：验证 filterHeadersBlacklist（敏感 header mask）
- [ ] 4.5.5 运行测试确认失败
**Green:**
- [ ] 4.5.6 实现 filterHeadersNone
- [ ] 4.5.7 实现 filterHeadersAll
- [ ] 4.5.8 实现 filterHeadersWhitelist
- [ ] 4.5.9 实现 filterHeadersBlacklist
- [ ] 4.5.10 实现 filterHeaders 统一入口
- [ ] 4.5.11 运行测试确认通过
**Refactor:**
- [ ] 4.5.12 检查并优化代码结构

### Step 4.6: 默认敏感 header 列表
**Red:**
- [ ] 4.6.1 编写测试：验证默认敏感 header 列表
- [ ] 4.6.2 运行测试确认失败
**Green:**
- [ ] 4.6.3 定义 defaultSensitiveHeaders（Authorization, Cookie, Set-Cookie, X-Api-Key, X-Auth-Token）
- [ ] 4.6.4 运行测试确认通过
**Refactor:**
- [ ] 4.6.5 检查并优化代码结构

### Step 4.7: Phase 后回归
- [ ] 4.7.1 运行 `go test ./...` 确保所有测试通过
- [ ] 4.7.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 5: Response Body 捕获

### Step 5.0: Phase 前回归
- [ ] 5.0.1 运行 `go test ./...` 确保所有测试通过

### Step 5.1: bodyCaptureWriter 结构体
**Red:**
- [ ] 5.1.1 编写测试：验证 bodyCaptureWriter 捕获写入的数据
- [ ] 5.1.2 编写测试：验证 bodyCaptureWriter 同时写入原始 ResponseWriter
- [ ] 5.1.3 运行测试确认失败
**Green:**
- [ ] 5.1.4 实现 bodyCaptureWriter 结构体
- [ ] 5.1.5 实现 Write(b []byte) 方法
- [ ] 5.1.6 实现 WriteString(s string) 方法
- [ ] 5.1.7 运行测试确认通过
**Refactor:**
- [ ] 5.1.8 检查并优化代码结构

### Step 5.2: Phase 后回归
- [ ] 5.2.1 运行 `go test ./...` 确保所有测试通过
- [ ] 5.2.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 6: Request Body 读取

### Step 6.0: Phase 前回归
- [ ] 6.0.1 运行 `go test ./...` 确保所有测试通过

### Step 6.1: readRequestBody 函数
**Red:**
- [ ] 6.1.1 编写测试：验证读取 request body 并保留供后续处理
- [ ] 6.1.2 运行测试确认失败
**Green:**
- [ ] 6.1.3 实现 readRequestBody 函数
- [ ] 6.1.4 运行测试确认通过
**Refactor:**
- [ ] 6.1.5 检查并优化代码结构

### Step 6.2: Phase 后回归
- [ ] 6.2.1 运行 `go test ./...` 确保所有测试通过
- [ ] 6.2.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 7: SSE 事件处理

### Step 7.0: Phase 前回归
- [ ] 7.0.1 运行 `go test ./...` 确保所有测试通过

### Step 7.1: parseSSEEvents 函数
**Red:**
- [ ] 7.1.1 编写测试：解析 SSE 事件（以 \n\n 或 \r\n\r\n 分隔）
- [ ] 7.1.2 运行测试确认失败
**Green:**
- [ ] 7.1.3 实现 parseSSEEvents(body string) []string
- [ ] 7.1.4 运行测试确认通过
**Refactor:**
- [ ] 7.1.5 检查并优化代码结构

### Step 7.2: truncateSSEEvents 函数
**Red:**
- [ ] 7.2.1 编写测试：验证所有 SSE 截取策略（含 overlap 边界）
- [ ] 7.2.2 运行测试确认失败
**Green:**
- [ ] 7.2.3 实现 truncateSSEEvents(events []string, config SSELogConfig) string
- [ ] 7.2.4 运行测试确认通过
**Refactor:**
- [ ] 7.2.5 检查并优化代码结构

### Step 7.3: Phase 后回归
- [ ] 7.3.1 运行 `go test ./...` 确保所有测试通过
- [ ] 7.3.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 8: 核心中间件实现

### Step 8.0: Phase 前回归
- [ ] 8.0.1 运行 `go test ./...` 确保所有测试通过

### Step 8.1: GinLogger 简单工厂函数
**Red:**
- [ ] 8.1.1 编写测试：验证 GinLogger(logger) 返回有效的 HandlerFunc
- [ ] 8.1.2 运行测试确认失败
**Green:**
- [ ] 8.1.3 实现 GinLogger(logger Logger) gin.HandlerFunc
- [ ] 8.1.4 运行测试确认通过
**Refactor:**
- [ ] 8.1.5 检查并优化代码结构

### Step 8.2: GinLoggerWithConfig 配置工厂函数
**Red:**
- [ ] 8.2.1 编写测试：验证 GinLoggerWithConfig 使用自定义配置
- [ ] 8.2.2 运行测试确认失败
**Green:**
- [ ] 8.2.3 实现 GinLoggerWithConfig(config GinLoggerConfig) gin.HandlerFunc
- [ ] 8.2.4 运行测试确认通过
**Refactor:**
- [ ] 8.2.5 检查并优化代码结构

### Step 8.3: 请求时间记录和延迟计算
**Red:**
- [ ] 8.3.1 编写测试：验证延迟计算正确
- [ ] 8.3.2 运行测试确认失败
**Green:**
- [ ] 8.3.3 实现请求时间记录和延迟计算
- [ ] 8.3.4 运行测试确认通过
**Refactor:**
- [ ] 8.3.5 检查并优化代码结构

### Step 8.4: 状态码到日志级别映射
**Red:**
- [ ] 8.4.1 编写测试：验证 2xx→Info, 4xx→Warn, 5xx→Error
- [ ] 8.4.2 运行测试确认失败
**Green:**
- [ ] 8.4.3 实现 statusToLevel(status int) 函数
- [ ] 8.4.4 运行测试确认通过
**Refactor:**
- [ ] 8.4.5 检查并优化代码结构

### Step 8.5: 基本日志字段输出
**Red:**
- [ ] 8.5.1 编写测试：验证 method, path, status, latency, client_ip, body_size 字段
- [ ] 8.5.2 运行测试确认失败
**Green:**
- [ ] 8.5.3 实现基本日志字段输出
- [ ] 8.5.4 运行测试确认通过
**Refactor:**
- [ ] 8.5.5 检查并优化代码结构

### Step 8.6: Phase 后回归
- [ ] 8.6.1 运行 `go test ./...` 确保所有测试通过
- [ ] 8.6.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 9: 功能集成

### Step 9.0: Phase 前回归
- [ ] 9.0.1 运行 `go test ./...` 确保所有测试通过

### Step 9.1: 跳过路径功能
**Red:**
- [ ] 9.1.1 编写测试：验证跳过指定路径
- [ ] 9.1.2 运行测试确认失败
**Green:**
- [ ] 9.1.3 实现 SkipPaths 配置解析和跳过逻辑
- [ ] 9.1.4 运行测试确认通过
**Refactor:**
- [ ] 9.1.5 检查并优化代码结构

### Step 9.2: TraceId 功能
**Red:**
- [ ] 9.2.1 编写测试：验证自动生成 traceId
- [ ] 9.2.2 编写测试：验证从请求头读取 traceId
- [ ] 9.2.3 运行测试确认失败
**Green:**
- [ ] 9.2.4 实现 traceId 读取、生成、注入
- [ ] 9.2.5 运行测试确认通过
**Refactor:**
- [ ] 9.2.6 检查并优化代码结构

### Step 9.3: 自定义字段功能
**Red:**
- [ ] 9.3.1 编写测试：验证 CustomFields 回调
- [ ] 9.3.2 运行测试确认失败
**Green:**
- [ ] 9.3.3 实现 CustomFields 回调支持
- [ ] 9.3.4 运行测试确认通过
**Refactor:**
- [ ] 9.3.5 检查并优化代码结构

### Step 9.4: Body 日志功能集成
**Red:**
- [ ] 9.4.1 编写测试：验证 request body 日志记录（含截取）
- [ ] 9.4.2 编写测试：验证 response body 日志记录（含截取）
- [ ] 9.4.3 编写测试：验证二进制 body 只记录类型和大小
- [ ] 9.4.4 运行测试确认失败
**Green:**
- [ ] 9.4.5 集成 request body 读取和日志记录
- [ ] 9.4.6 集成 response body 捕获和日志记录
- [ ] 9.4.7 运行测试确认通过
**Refactor:**
- [ ] 9.4.8 检查并优化代码结构

### Step 9.5: Header 日志功能集成
**Red:**
- [ ] 9.5.1 编写测试：验证 request header 过滤和敏感处理
- [ ] 9.5.2 编写测试：验证 response header 过滤和敏感处理
- [ ] 9.5.3 运行测试确认失败
**Green:**
- [ ] 9.5.4 集成 request header 过滤和日志记录
- [ ] 9.5.5 集成 response header 过滤和日志记录
- [ ] 9.5.6 运行测试确认通过
**Refactor:**
- [ ] 9.5.7 检查并优化代码结构

### Step 9.6: SSE 日志功能集成
**Red:**
- [ ] 9.6.1 编写测试：验证 SSE 响应的事件截取
- [ ] 9.6.2 运行测试确认失败
**Green:**
- [ ] 9.6.3 集成 SSE 日志处理
- [ ] 9.6.4 运行测试确认通过
**Refactor:**
- [ ] 9.6.5 检查并优化代码结构

### Step 9.7: 错误处理
**Red:**
- [ ] 9.7.1 编写测试：验证 gin.Errors 中的错误信息提取
- [ ] 9.7.2 运行测试确认失败
**Green:**
- [ ] 9.7.3 实现错误信息提取和日志记录
- [ ] 9.7.4 运行测试确认通过
**Refactor:**
- [ ] 9.7.5 检查并优化代码结构

### Step 9.8: Phase 后回归
- [ ] 9.8.1 运行 `go test ./...` 确保所有测试通过
- [ ] 9.8.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 10: 代码审查

### Step 10.0: Phase 前回归
- [ ] 10.0.1 运行 `go test ./...` 确保所有测试通过

### Step 10.1: 全量 Review
- [ ] 10.1.1 运行 `superpowers:requesting-code-review` skill
- [ ] 10.1.2 处理审查反馈，修复问题
- [ ] 10.1.3 确保审查通过

### Step 10.2: Phase 后回归
- [ ] 10.2.1 运行 `go test ./...` 确保所有测试通过
- [ ] 10.2.2 确保行覆盖率和分支覆盖率均为 100%

---

## Phase 11: 发布 go-comm

- [ ] 11.1 提交 go-comm 的 gin logger 变更
- [ ] 11.2 打 tag v2.7.2
- [ ] 11.3 push commit 和 tag 到远程仓库

---

## Phase 12: 集成到 mobile-claude

> **重要规则**：在本 Phase 的任何阶段发现测试失败，即使不是本次修改导致的，也要**无条件立刻修复**，不可以以任何理由回避

### Step 12.0: Phase 前准备
- [ ] 12.0.1 切换到 mobile-claude 项目目录
- [ ] 12.0.2 确认当前 mobile-claude 代码已提交或暂存

### Step 12.1: Phase 前全量回归（mise test-all）
- [ ] 12.1.1 运行 `mise test-all` 进行全量回归测试
- [ ] 12.1.2 如有任何测试失败，**无条件立刻修复**（即使不是本次修改导致）
- [ ] 12.1.3 确认 `mise test-all` 完全通过后再继续

### Step 12.2: 更新依赖和 import 引用
**Red:**
- [ ] 12.2.1 编写测试：验证新 gin logger 中间件的使用（预期失败，依赖未更新）
- [ ] 12.2.2 运行测试确认失败
**Green:**
- [ ] 12.2.3 更新 go-comm 依赖到 v2.7.2
- [ ] 12.2.4 更新所有 import 引用：
  - `github.com/qiangyt/go-comm` 中的 gin 相关 → `github.com/qiangyt/go-comm/qgin`
  - `github.com/qiangyt/go-comm` 中的 config 相关 → `github.com/qiangyt/go-comm/qconfig`
  - `github.com/qiangyt/go-comm` 中的 io/file 相关 → `github.com/qiangyt/go-comm/qio`
- [ ] 12.2.5 替换原有的 gin logger，使用 go-comm/qgin 的新 gin logger 中间件
- [ ] 12.2.6 运行测试确认通过
**Refactor:**
- [ ] 12.2.7 检查并优化代码结构
- [ ] 12.2.8 如有任何问题，**无条件立刻修复**

### Step 12.3: Phase 后全量回归（mise test-all）
- [ ] 12.3.1 运行 `mise test-all` 进行全量回归测试
- [ ] 12.3.2 如有任何测试失败，**无条件立刻修复**（即使不是本次修改导致）
- [ ] 12.3.3 确保行覆盖率和分支覆盖率均为 100%
- [ ] 12.3.4 确认 `mise test-all` 完全通过

### Step 12.4: 提交变更
- [ ] 12.4.1 提交 mobile-claude 的变更

