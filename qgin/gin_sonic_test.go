package qgin

import (
	"bytes"
	"testing"

	"github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/require"
)

func TestConfigureGinWithSonic(t *testing.T) {
	// 测试可以安全地多次调用
	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
	})

	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
	})
}

func TestSonicAPI_Marshal(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestSonicAPI_Unmarshal(t *testing.T) {
	ConfigureGinWithSonic()

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	err := json.API.Unmarshal([]byte(jsonData), &result)
	require.Nil(t, err)
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])
}

func TestSonicAPI_EscapeHTML(t *testing.T) {
	// 测试 EscapeHTML: true 配置
	ConfigureGinWithSonic()

	data := map[string]any{
		"html": "<script>alert('xss')</script>",
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	// EscapeHTML: true 时，特殊字符应该被转义
	require.Contains(t, string(result), `\u003c`)
	require.Contains(t, string(result), `\u003e`)
}

// TestSonicAPI_Interface 确保 sonicAPI 实现了 json.Core 接口
func TestSonicAPI_Interface(t *testing.T) {
	ConfigureGinWithSonic()

	// 编译时检查接口实现 - 使用指针类型检查
	var _ json.Core = (*SonicAPI)(nil)

	// 运行时验证接口方法 - 使用已配置的 json.API
	api := json.API

	data := map[string]string{"key": "value"}

	// Marshal
	_, err := api.Marshal(data)
	require.Nil(t, err)

	// Unmarshal
	var parsed map[string]string
	err = api.Unmarshal([]byte(`{"key":"value"}`), &parsed)
	require.Nil(t, err)

	// MarshalIndent
	_, err = api.MarshalIndent(data, "", "  ")
	require.Nil(t, err)

	// NewEncoder
	var buf bytes.Buffer
	encoder := api.NewEncoder(&buf)
	require.NotNil(t, encoder)

	// NewDecoder
	reader := bytes.NewBufferString(`{}`)
	decoder := api.NewDecoder(reader)
	require.NotNil(t, decoder)
}

// TestBackwardCompatibility 测试向后兼容的包级函数
func TestBackwardCompatibility(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	// 测试 Marshal 函数（向后兼容）
	result, err := Marshal(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)

	// 测试 Unmarshal 函数（向后兼容）
	var parsed map[string]any
	err = Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, "test", parsed["name"])

	// 测试 MarshalIndent 函数（向后兼容）
	_, err = MarshalIndent(data, "", "  ")
	require.Nil(t, err)

	// 测试 NewEncoder 函数（向后兼容）
	var buf bytes.Buffer
	encoder := NewEncoder(&buf)
	require.NotNil(t, encoder)

	// 测试 NewDecoder 函数（向后兼容）
	reader := bytes.NewBufferString(`{}`)
	decoder := NewDecoder(reader)
	require.NotNil(t, decoder)
}
