## 为什么

go-comm 库目前使用 Go 标准库 `encoding/json` 进行 JSON 操作。将其替换为 `github.com/bytedance/sonic` v1.15.0 可以提供 2-5 倍更快的 JSON 编码/解码性能，且代码改动极小（API 兼容的即插即用替换）。此外，升级 gin 到 v1.12.0 可以为 gin 用户提供内置的 sonic 支持，使他们只需调用一个简单的辅助函数即可获得性能提升。

## 变更内容

- 在 go-comm 内部代码中将 `encoding/json` 替换为 `github.com/bytedance/sonic` v1.15.0
- 升级 `github.com/gin-gonic/gin` 到 v1.12.0 及相关 gin 依赖
- 添加 `gin_sonic.go` 辅助函数，为 gin 框架用户启用 gin 内置的 sonic 编解码器
- 更新 5 个使用 JSON 操作的文件：`file.go`、`command.go`、`lock_file.go`、`lock_file_posix.go`、`lock_file_windows.go`

## 能力

### 新增能力

- `gin-sonic-helper`：提供 `ConfigureGinWithSonic()` 函数，启用 gin 内置的 sonic 编解码器用于 JSON 序列化/反序列化

### 修改的能力

无 - 这是内部性能改进，没有 API 行为变更。

## 影响

**修改的文件:**
- `go.mod` - 添加 sonic v1.15.0，升级 gin 到 v1.12.0
- `file.go` - 将 encoding/json 导入替换为 sonic
- `command.go` - 将 encoding/json 导入替换为 sonic
- `lock_file.go` - 将 encoding/json 导入替换为 sonic
- `lock_file_posix.go` - 将 encoding/json 导入替换为 sonic
- `lock_file_windows.go` - 将 encoding/json 导入替换为 sonic

**新增文件:**
- `gin_sonic.go` - Gin-sonic 集成辅助函数
- `gin_sonic_test.go` - Gin 集成测试

**依赖变更:**
- `github.com/bytedance/sonic` v1.15.0（新增）
- `github.com/gin-gonic/gin` v1.12.0（升级）

**受影响的消费者:**
- 使用 go-comm 的项目将自动获得更快的 JSON 操作性能
- 使用 gin 和 go-comm 的项目可以调用 `ConfigureGinWithSonic()` 来为 gin 的 JSON 处理启用 sonic
