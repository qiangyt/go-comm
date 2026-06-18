package qgin

import (
	"io"

	"github.com/gin-gonic/gin/codec/json"
)

// Marshal 使用 sonic 序列化对象
//
// Deprecated: 此方法仅为向后兼容保留，请使用 ConfigureGinWithSonic() 后直接调用 json.API.Marshal()
func Marshal(v any) ([]byte, error) {
	return json.API.Marshal(v)
}

// Unmarshal 使用 sonic 反序列化数据
//
// Deprecated: 此方法仅为向后兼容保留，请使用 ConfigureGinWithSonic() 后直接调用 json.API.Unmarshal()
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
