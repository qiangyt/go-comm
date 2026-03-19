package comm

import (
	"io"

	"github.com/gin-gonic/gin/codec/json"
	"github.com/qiangyt/go-comm/v2/qgin"
)

// sonicAPI 实现 gin 的 json.Core 接口，使用 sonic 作为后端
//
// Deprecated: 请使用 github.com/qiangyt/go-comm/v2/qgin 包中的类型
type sonicAPI = qgin.SonicAPI

// Marshal 使用 sonic 序列化对象
//
// Deprecated: 此方法仅为向后兼容保留，请使用 qgin.ConfigureGinWithSonic() 后直接调用 json.API.Marshal()
func Marshal(v any) ([]byte, error) {
	return json.API.Marshal(v)
}

// Unmarshal 使用 sonic 反序列化数据
//
// Deprecated: 此方法仅为向后兼容保留，请使用 qgin.ConfigureGinWithSonic() 后直接调用 json.API.Unmarshal()
func Unmarshal(data []byte, v any) error {
	return json.API.Unmarshal(data, v)
}

// MarshalIndent 使用 sonic 序列化对象并添加缩进
//
// Deprecated: 此方法仅为向后兼容保留
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.API.MarshalIndent(v, prefix, indent)
}

// NewEncoder 创建 sonic 编码器
//
// Deprecated: 此方法仅为向后兼容保留
func NewEncoder(writer io.Writer) json.Encoder {
	return json.API.NewEncoder(writer)
}

// NewDecoder 创建 sonic 解码器
//
// Deprecated: 此方法仅为向后兼容保留
func NewDecoder(reader io.Reader) json.Decoder {
	return json.API.NewDecoder(reader)
}

// ConfigureGinWithSonic 配置 gin 使用 sonic 作为 JSON 编解码器
// gin v1.12.0+ 支持通过替换 json.API 来切换 JSON 后端
//
// Deprecated: 请使用 qgin.ConfigureGinWithSonic()
var ConfigureGinWithSonic = qgin.ConfigureGinWithSonic

// ConfigureGinWithSonicConfig 使用自定义 sonic 配置来配置 gin
// 允许更精细地控制 sonic 的行为
//
// Deprecated: 请使用 qgin.ConfigureGinWithSonicConfig()
var ConfigureGinWithSonicConfig = qgin.ConfigureGinWithSonicConfig

// 以下类型别名用于向后兼容

// SonicAPI 是 qgin.SonicAPI 的别名
//
// Deprecated: 请使用 qgin.SonicAPI
type SonicAPI = qgin.SonicAPI
