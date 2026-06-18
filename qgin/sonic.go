package qgin

import (
	"io"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin/codec/json"
)

// SonicAPI 实现 gin 的 json.Core 接口，使用 sonic 作为后端
type SonicAPI struct {
	api sonic.API
}

// Marshal 使用 sonic 序列化对象
func (s SonicAPI) Marshal(v any) ([]byte, error) {
	return s.api.Marshal(v)
}

// Unmarshal 使用 sonic 反序列化数据
func (s SonicAPI) Unmarshal(data []byte, v any) error {
	return s.api.Unmarshal(data, v)
}

// MarshalIndent 使用 sonic 序列化对象并添加缩进
func (s SonicAPI) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return s.api.MarshalIndent(v, prefix, indent)
}

// NewEncoder 创建 sonic 编码器
func (s SonicAPI) NewEncoder(writer io.Writer) json.Encoder {
	return s.api.NewEncoder(writer)
}

// NewDecoder 创建 sonic 解码器
func (s SonicAPI) NewDecoder(reader io.Reader) json.Decoder {
	return s.api.NewDecoder(reader)
}

// ConfigureGinWithSonic 配置 gin 使用 sonic 作为 JSON 编解码器
// gin v1.12.0+ 支持通过替换 json.API 来切换 JSON 后端
func ConfigureGinWithSonic() {
	json.API = SonicAPI{
		api: sonic.Config{
			EscapeHTML: true,
		}.Froze(),
	}
}

// ConfigureGinWithSonicConfig 使用自定义 sonic 配置来配置 gin
// 允许更精细地控制 sonic 的行为
func ConfigureGinWithSonicConfig(config sonic.Config) {
	json.API = SonicAPI{
		api: config.Froze(),
	}
}
