## 背景

go-comm 库是一个通用工具包，提供 JSON 解析、文件处理和其他工具。目前，它使用 Go 标准库 `encoding/json` 包进行所有 JSON 操作。此次变更将其替换为字节跳动的 sonic 库，在保持 API 兼容性的同时获得更好的性能。

**当前状态:**
- 5 个文件使用 `encoding/json`：`file.go`、`command.go`、`lock_file.go`、`lock_file_posix.go`、`lock_file_windows.go`
- 操作包括：`json.Marshal()` 和 `json.Unmarshal()` 用于配置解析、命令输出解析和锁文件数据

**约束条件:**
- 必须保持 API 兼容性（sonic 是即插即用的替换）
- Go 版本：1.25
- Gin v1.12.0 已包含内置 sonic 支持

## 目标 / 非目标

**目标:**
- 将 encoding/json 替换为 sonic v1.15.0 用于内部 JSON 操作
- 升级 gin 到 v1.12.0 及相关依赖
- 为 gin 用户提供辅助函数以启用 sonic 编解码器
- 保持 100% API 兼容性
- 记录性能对比数据

**非目标:**
- 更改 `github.com/iancoleman/orderedmap`（由该库自行处理）
- 对任何公共 API 进行破坏性更改

## 决策

### 决策 1: 使用 sonic v1.15.0（而非最新版本）

**理由:** 用户明确要求 v1.15.0 以确保稳定性。Sonic 在各版本间保持 API 兼容性，因此 v1.15.0 是一个安全、生产就绪的选择。

**备选方案:** 最新版本 - 被拒绝以确保可预测的行为。

### 决策 2: 利用 gin 的内置 sonic 支持

**理由:** Gin v1.12.0 通过 `codec/json/sonic.go` 包含原生 sonic 支持。无需实现自定义 `json.Core` 接口，我们只需提供一个调用 gin 内置 `json.EnableSonic()` 的辅助函数。

**备选方案:** 自定义 `json.Core` 实现 - 被拒绝，因为 gin 已内置此功能。

### 决策 3: 直接替换导入

**理由:** Sonic 与 encoding/json API 兼容。只需将导入从 `encoding/json` 改为 `github.com/bytedance/sonic` 即可 - 无需修改代码。

**备选方案:** 包装器抽象 - 被拒绝，因为是不必要的复杂性。

## 风险与权衡

**风险:** Sonic 可能在边缘情况下与 encoding/json 有行为差异
→ **缓解措施:** 运行现有测试套件验证兼容性；sonic 设计为即插即用替换

**风险:** Sonic 使用 CGO 进行性能优化（汇编优化）
→ **缓解措施:** 存在纯 Go 回退；即使没有汇编，性能仍优于 encoding/json

**风险:** 新依赖增加构建时间
→ **缓解措施:** Sonic 是轻量级的；构建时间影响可忽略不计
