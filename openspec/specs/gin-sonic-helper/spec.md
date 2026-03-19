## ADDED Requirements

### Requirement: ConfigureGinWithSonic 函数

系统应提供 `ConfigureGinWithSonic()` 函数，用于启用 gin 内置的 sonic 编解码器进行 JSON 序列化和反序列化。

#### Scenario: 为 gin 启用 sonic

- **WHEN** 在应用初始化期间调用 `ConfigureGinWithSonic()`
- **THEN** gin 的 JSON 编解码器被配置为使用 sonic 处理所有后续 JSON 操作

#### Scenario: 函数可以安全地多次调用

- **WHEN** `ConfigureGinWithSonic()` 被多次调用
- **THEN** 不会发生错误或 panic
- **AND** gin 仍保持使用 sonic 配置

### Requirement: Sonic 内部替换 encoding/json

系统应使用 sonic v1.15.0 进行所有内部 JSON marshal 和 unmarshal 操作，替换 encoding/json。

#### Scenario: FromJson 使用 sonic

- **WHEN** 调用 `FromJson()` 解析 JSON 文本
- **THEN** 使用 sonic 进行 JSON unmarshaling

#### Scenario: CreateLockFile 使用 sonic

- **WHEN** `CreateLockFile()` 序列化锁文件数据
- **THEN** 使用 sonic 进行 JSON marshaling

#### Scenario: ParseCommandOutput 使用 sonic

- **WHEN** `ParseCommandOutput()` 解析 JSON 命令输出
- **THEN** 使用 sonic 进行 JSON unmarshaling
