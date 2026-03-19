package qgin

import (
	"bytes"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/require"
)

// TestQginPackageImport 测试 qgin 包可以正常导入和使用
func TestQginPackageImport(t *testing.T) {
	// 测试可以安全地多次调用 ConfigureGinWithSonic
	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
	})

	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
	})
}

// TestConfigureGinWithSonicConfig 测试自定义配置
func TestConfigureGinWithSonicConfig(t *testing.T) {
	require.NotPanics(t, func() {
		ConfigureGinWithSonicConfig(sonic.Config{
			EscapeHTML: false,
		})
	})

	// 恢复默认配置
	ConfigureGinWithSonic()
}

// TestSonicAPI_Marshal 测试 Marshal 功能
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

// TestSonicAPI_Unmarshal 测试 Unmarshal 功能
func TestSonicAPI_Unmarshal(t *testing.T) {
	ConfigureGinWithSonic()

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	err := json.API.Unmarshal([]byte(jsonData), &result)
	require.Nil(t, err)
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])
}

// TestSonicAPI_EscapeHTML 测试 HTML 转义
func TestSonicAPI_EscapeHTML(t *testing.T) {
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

// TestSonicAPI_Interface 确保 SonicAPI 实现了 json.Core 接口
func TestSonicAPI_Interface(t *testing.T) {
	ConfigureGinWithSonic()

	// 编译时检查接口实现
	var _ json.Core = SonicAPI{}

	// 运行时验证接口方法
	var api json.Core = SonicAPI{
		api: sonic.ConfigStd,
	}

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
