package qjson

import (
	"encoding/json"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/qerr"
	"github.com/qiangyt/go-comm/v3/qsys"
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

// JsonUnmarshal 根据 JSONConfig.Backend 选择对应的 unmarshal 实现
func JsonUnmarshal(data []byte, v any) error {
	if JSONConfig.Backend == JSONBackendStdlib {
		return json.Unmarshal(data, v)
	}
	return sonic.Unmarshal(data, v)
}

// JsonMarshal 根据 JSONConfig.Backend 选择对应的 marshal 实现
func JsonMarshal(v any) ([]byte, error) {
	if JSONConfig.Backend == JSONBackendStdlib {
		return json.Marshal(v)
	}
	return sonic.Marshal(v)
}

// UnmarshalJSON 公开的 JSON 反序列化函数，根据 JSONConfig.Backend 选择后端
// 使用方法与 json.Unmarshal 相同，但支持动态切换 sonic/stdlib
func UnmarshalJSON(data []byte, v any) error {
	return JsonUnmarshal(data, v)
}

// UnmarshalJSONP panic 版本的 JSON 反序列化函数
// 失败时 panic，符合项目 Error 处理规范
func UnmarshalJSONP(data []byte, v any) {
	if err := JsonUnmarshal(data, v); err != nil {
		panic(qerr.NewSystemError("JSON 反序列化失败", err))
	}
}

// MarshalJSON 公开的 JSON 序列化函数，根据 JSONConfig.Backend 选择后端
// 使用方法与 json.Marshal 相同，但支持动态切换 sonic/stdlib
func MarshalJSON(v any) ([]byte, error) {
	return JsonMarshal(v)
}

// MarshalJSONP panic 版本的 JSON 序列化函数
// 失败时 panic，符合项目 Error 处理规范
func MarshalJSONP(v any) []byte {
	result, err := JsonMarshal(v)
	if err != nil {
		panic(qerr.NewSystemError("JSON 序列化失败", err))
	}
	return result
}

// MarshalJSONIndent 公开的 JSON 带缩进序列化函数
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	if JSONConfig.Backend == JSONBackendStdlib {
		return json.MarshalIndent(v, prefix, indent)
	}
	return sonic.MarshalIndent(v, prefix, indent)
}

// MarshalJSONIndentP panic 版本的 JSON 带缩进序列化函数
// 失败时 panic，符合项目 Error 处理规范
func MarshalJSONIndentP(v any, prefix, indent string) []byte {
	result, err := MarshalJSONIndent(v, prefix, indent)
	if err != nil {
		panic(qerr.NewSystemError("JSON 序列化失败", err))
	}
	return result
}

func FromJsonP(jsonText string, envsubt bool, result any) {
	if err := FromJson(jsonText, envsubt, result); err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
}

func FromJson(jsonText string, envsubt bool, result any) (err error) {
	if envsubt {
		jsonText, err = qsys.EnvSubst(jsonText, nil)
		if err != nil {
			return err
		}
	}

	if err = JsonUnmarshal([]byte(jsonText), result); err != nil {
		return errors.Wrapf(err, "parse json: \n\n%s", jsonText)
	}
	return nil
}

func MapFromJsonP(yamlText string, envsubt bool) map[string]any {
	r, err := MapFromJson(yamlText, envsubt)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

func MapFromJson(yamlText string, envsubt bool) (map[string]any, error) {
	r := map[string]any{}
	if err := FromJson(yamlText, envsubt, &r); err != nil {
		return nil, err
	}

	return r, nil
}
