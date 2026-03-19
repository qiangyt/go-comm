package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON_Sonic(t *testing.T) {
	// 设置为 sonic 后端
	JSONConfig.Backend = JSONBackendSonic

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	err := UnmarshalJSON([]byte(jsonData), &result)
	require.Nil(t, err)
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])
}

func TestUnmarshalJSON_Stdlib(t *testing.T) {
	// 设置为 stdlib 后端
	JSONConfig.Backend = JSONBackendStdlib

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	err := UnmarshalJSON([]byte(jsonData), &result)
	require.Nil(t, err)
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}

func TestMarshalJSON_Sonic(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := MarshalJSON(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestMarshalJSON_Stdlib(t *testing.T) {
	JSONConfig.Backend = JSONBackendStdlib

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := MarshalJSON(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}

func TestMarshalJSONIndent_Sonic(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := MarshalJSONIndent(data, "", "  ")
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestMarshalJSONIndent_Stdlib(t *testing.T) {
	JSONConfig.Backend = JSONBackendStdlib

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := MarshalJSONIndent(data, "", "  ")
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}

func TestUnmarshalJSON_InvalidJSON(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	var result map[string]any
	err := UnmarshalJSON([]byte(`{invalid}`), &result)
	require.NotNil(t, err)
}

func TestMarshalJSON_ComplexData(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	type Nested struct {
		Items []int `json:"items"`
	}

	data := Nested{Items: []int{1, 2, 3}}

	result, err := MarshalJSON(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"items"`)
	require.Contains(t, string(result), `[1,2,3]`)

	var parsed Nested
	err = UnmarshalJSON(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, []int{1, 2, 3}, parsed.Items)
}

// ===== Panic 版本测试 =====

func TestUnmarshalJSONP_Sonic(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	require.NotPanics(t, func() {
		UnmarshalJSONP([]byte(jsonData), &result)
	})
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])
}

func TestUnmarshalJSONP_Stdlib(t *testing.T) {
	JSONConfig.Backend = JSONBackendStdlib

	jsonData := `{"name": "test", "value": 123}`
	var result map[string]any

	require.NotPanics(t, func() {
		UnmarshalJSONP([]byte(jsonData), &result)
	})
	require.Equal(t, "test", result["name"])
	require.Equal(t, float64(123), result["value"])

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}

func TestUnmarshalJSONP_PanicsOnInvalidJSON(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	var result map[string]any
	require.Panics(t, func() {
		UnmarshalJSONP([]byte(`{invalid}`), &result)
	})
}

func TestMarshalJSONP_Sonic(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	var result []byte
	require.NotPanics(t, func() {
		result = MarshalJSONP(data)
	})
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestMarshalJSONP_Stdlib(t *testing.T) {
	JSONConfig.Backend = JSONBackendStdlib

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	var result []byte
	require.NotPanics(t, func() {
		result = MarshalJSONP(data)
	})
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}

func TestMarshalJSONIndentP_Sonic(t *testing.T) {
	JSONConfig.Backend = JSONBackendSonic

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	var result []byte
	require.NotPanics(t, func() {
		result = MarshalJSONIndentP(data, "", "  ")
	})
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestMarshalJSONIndentP_Stdlib(t *testing.T) {
	JSONConfig.Backend = JSONBackendStdlib

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	var result []byte
	require.NotPanics(t, func() {
		result = MarshalJSONIndentP(data, "", "  ")
	})
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)

	// 恢复默认
	JSONConfig.Backend = JSONBackendSonic
}
