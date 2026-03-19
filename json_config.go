package comm

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

// JSONBackend 定义 JSON 后端类型
type JSONBackend int

const (
	// JSONBackendSonic 使用 ByteDance sonic 库（默认，性能更优）
	JSONBackendSonic JSONBackend = iota
	// JSONBackendStdlib 使用 Go 标准库 encoding/json
	JSONBackendStdlib
)

// JSONConfig 全局 JSON 配置
var JSONConfig = struct {
	// Backend 控制使用哪个 JSON 后端
	// 默认为 JSONBackendSonic
	Backend JSONBackend
}{
	Backend: JSONBackendSonic,
}

// jsonUnmarshal 根据 JSONConfig.Backend 选择对应的 unmarshal 实现
func jsonUnmarshal(data []byte, v any) error {
	if JSONConfig.Backend == JSONBackendStdlib {
		return json.Unmarshal(data, v)
	}
	return sonic.Unmarshal(data, v)
}

// jsonMarshal 根据 JSONConfig.Backend 选择对应的 marshal 实现
func jsonMarshal(v any) ([]byte, error) {
	if JSONConfig.Backend == JSONBackendStdlib {
		return json.Marshal(v)
	}
	return sonic.Marshal(v)
}
