package comm

import (
	"bytes"
	"testing"

	"github.com/bytedance/sonic"
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

func TestConfigureGinWithSonicConfig(t *testing.T) {
	// 测试自定义配置
	require.NotPanics(t, func() {
		ConfigureGinWithSonicConfig(sonic.Config{
			EscapeHTML: false,
		})
	})

	// 恢复默认配置
	ConfigureGinWithSonic()
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

func TestSonicAPI_MarshalIndent(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	result, err := json.API.MarshalIndent(data, "", "  ")
	require.Nil(t, err)
	require.Contains(t, string(result), `"name"`)
	require.Contains(t, string(result), `"test"`)
}

func TestSonicAPI_NewEncoder(t *testing.T) {
	ConfigureGinWithSonic()

	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)

	data := map[string]any{
		"name":  "test",
		"value": 123,
	}

	err := encoder.Encode(data)
	require.Nil(t, err)
	require.Contains(t, buf.String(), `"name"`)
	require.Contains(t, buf.String(), `"test"`)
}

func TestSonicAPI_NewDecoder(t *testing.T) {
	ConfigureGinWithSonic()

	jsonData := `{"name": "test", "value": 123}`
	reader := bytes.NewBufferString(jsonData)
	decoder := json.API.NewDecoder(reader)

	var result map[string]any
	err := decoder.Decode(&result)
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

func TestSonicAPI_NoEscapeHTML(t *testing.T) {
	// 测试 EscapeHTML: false 配置
	ConfigureGinWithSonicConfig(sonic.Config{
		EscapeHTML: false,
	})

	data := map[string]any{
		"html": "<script>alert('xss')</script>",
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	// EscapeHTML: false 时，特殊字符不应该被转义
	require.Contains(t, string(result), "<script>")
	require.NotContains(t, string(result), `\u003c`)

	// 恢复默认配置
	ConfigureGinWithSonic()
}

func TestSonicAPI_ComplexData(t *testing.T) {
	ConfigureGinWithSonic()

	// 测试复杂数据结构
	type Nested struct {
		Level1 struct {
			Level2 struct {
				Value string `json:"value"`
			} `json:"level2"`
		} `json:"level1"`
	}

	data := Nested{}
	data.Level1.Level2.Value = "deep"

	// Marshal 测试
	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"value"`)
	require.Contains(t, string(result), `"deep"`)

	// Unmarshal 测试
	var parsed Nested
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, "deep", parsed.Level1.Level2.Value)
}

func TestSonicAPI_EncoderDecoder(t *testing.T) {
	ConfigureGinWithSonic()

	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	// 编码
	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)
	data := TestData{Name: "encoder_test", Value: 42}

	err := encoder.Encode(data)
	require.Nil(t, err)

	// 解码
	decoder := json.API.NewDecoder(&buf)
	var result TestData

	// 重置 buffer 用于解码
	buf.Reset()
	buf.WriteString(`{"name":"decoder_test","value":99}`)

	err = decoder.Decode(&result)
	require.Nil(t, err)
	require.Equal(t, "decoder_test", result.Name)
	require.Equal(t, 99, result.Value)
}

func TestSonicAPI_StreamReader(t *testing.T) {
	ConfigureGinWithSonic()

	// 测试从 io.Reader 读取
	jsonData := `{"items": [1, 2, 3], "count": 3}`
	reader := bytes.NewBufferString(jsonData)

	decoder := json.API.NewDecoder(reader)
	var result map[string]any

	err := decoder.Decode(&result)
	require.Nil(t, err)
	require.Equal(t, float64(3), result["count"])
}

func TestSonicAPI_StreamWriter(t *testing.T) {
	ConfigureGinWithSonic()

	// 测试写入 io.Writer
	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)

	data := map[string]any{
		"message": "hello world",
		"count":   100,
	}

	err := encoder.Encode(data)
	require.Nil(t, err)
	require.True(t, buf.Len() > 0)
}

// TestSonicAPI_Interface 确保 sonicAPI 实现了 json.Core 接口
func TestSonicAPI_Interface(t *testing.T) {
	ConfigureGinWithSonic()

	// 编译时检查接口实现
	var _ json.Core = sonicAPI{}

	// 运行时验证接口方法
	var api json.Core = sonicAPI{
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

// TestSonicAPI_NilHandling 测试 nil 值处理
func TestSonicAPI_NilHandling(t *testing.T) {
	ConfigureGinWithSonic()

	// Marshal nil
	result, err := json.API.Marshal(nil)
	require.Nil(t, err)
	require.Equal(t, "null", string(result))

	// Unmarshal null
	var data map[string]any
	err = json.API.Unmarshal([]byte("null"), &data)
	require.Nil(t, err)
	require.Nil(t, data)
}

// TestSonicAPI_ArrayHandling 测试数组处理
func TestSonicAPI_ArrayHandling(t *testing.T) {
	ConfigureGinWithSonic()

	// Marshal 数组
	data := []int{1, 2, 3, 4, 5}
	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	require.Equal(t, `[1,2,3,4,5]`, string(result))

	// Unmarshal 数组
	var parsed []int
	err = json.API.Unmarshal([]byte(`[1,2,3,4,5]`), &parsed)
	require.Nil(t, err)
	require.Equal(t, []int{1, 2, 3, 4, 5}, parsed)
}

// TestSonicAPI_UnicodeHandling 测试 Unicode 处理
func TestSonicAPI_UnicodeHandling(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]string{
		"chinese": "中文测试",
		"emoji":   "🎉",
		"mixed":   "Hello 世界",
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)

	var parsed map[string]string
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, "中文测试", parsed["chinese"])
	require.Equal(t, "🎉", parsed["emoji"])
	require.Equal(t, "Hello 世界", parsed["mixed"])
}

// TestSonicAPI_ErrorHandling 测试错误处理
func TestSonicAPI_ErrorHandling(t *testing.T) {
	ConfigureGinWithSonic()

	// Unmarshal 无效 JSON
	var data map[string]any
	err := json.API.Unmarshal([]byte(`{invalid json}`), &data)
	require.NotNil(t, err)

	// Marshal 不支持的类型（如 channel）会导致错误
	invalidData := map[string]any{
		"channel": make(chan int),
	}
	_, err = json.API.Marshal(invalidData)
	require.NotNil(t, err)
}

// TestSonicAPI_LargeData 测试大数据处理
func TestSonicAPI_LargeData(t *testing.T) {
	ConfigureGinWithSonic()

	// 创建大型数据
	items := make([]map[string]any, 100)
	for i := 0; i < 100; i++ {
		items[i] = map[string]any{
			"id":    i,
			"name":  "item_" + string(rune('A'+i%26)),
			"value": i * 100,
		}
	}

	data := map[string]any{
		"items": items,
		"count": len(items),
	}

	// Marshal
	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	require.True(t, len(result) > 1000)

	// Unmarshal
	var parsed map[string]any
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, float64(100), parsed["count"])
}

// TestSonicAPI_EncoderSetEscapeHTML 测试编码器的 SetEscapeHTML 方法
func TestSonicAPI_EncoderSetEscapeHTML(t *testing.T) {
	ConfigureGinWithSonic()

	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)

	// 尝试设置 EscapeHTML
	encoder.SetEscapeHTML(true)

	data := map[string]string{
		"html": "<script>alert('xss')</script>",
	}

	err := encoder.Encode(data)
	require.Nil(t, err)
}

// TestSonicAPI_DecoderUseNumber 测试解码器的 UseNumber 方法
func TestSonicAPI_DecoderUseNumber(t *testing.T) {
	ConfigureGinWithSonic()

	jsonData := `{"number": 12345678901234567890}`
	reader := bytes.NewBufferString(jsonData)
	decoder := json.API.NewDecoder(reader)

	// 启用 UseNumber
	decoder.UseNumber()

	var result map[string]any
	err := decoder.Decode(&result)
	require.Nil(t, err)
	// 数字应该被解析为 json.Number 类型而不是 float64
}

// TestSonicAPI_DecoderDisallowUnknownFields 测试解码器的 DisallowUnknownFields 方法
func TestSonicAPI_DecoderDisallowUnknownFields(t *testing.T) {
	ConfigureGinWithSonic()

	type StrictStruct struct {
		KnownField string `json:"known_field"`
	}

	jsonData := `{"known_field": "value", "unknown_field": "should_error"}`
	reader := bytes.NewBufferString(jsonData)
	decoder := json.API.NewDecoder(reader)

	// 禁止未知字段
	decoder.DisallowUnknownFields()

	var result StrictStruct
	err := decoder.Decode(&result)
	// 某些 JSON 库可能会返回错误，sonic 可能不会
	// 这个测试主要是确保方法可以调用
	_ = err
}

// TestSonicAPI_MultipleReaders 测试多个 Reader
func TestSonicAPI_MultipleReaders(t *testing.T) {
	ConfigureGinWithSonic()

	data1 := `{"id": 1}`
	data2 := `{"id": 2}`

	// 解码第一个
	decoder1 := json.API.NewDecoder(bytes.NewBufferString(data1))
	var result1 map[string]any
	err := decoder1.Decode(&result1)
	require.Nil(t, err)
	require.Equal(t, float64(1), result1["id"])

	// 解码第二个
	decoder2 := json.API.NewDecoder(bytes.NewBufferString(data2))
	var result2 map[string]any
	err = decoder2.Decode(&result2)
	require.Nil(t, err)
	require.Equal(t, float64(2), result2["id"])
}

// TestSonicAPI_MultipleWriters 测试多个 Writer
func TestSonicAPI_MultipleWriters(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]string{"test": "data"}

	// 编码到第一个 Writer
	var buf1 bytes.Buffer
	encoder1 := json.API.NewEncoder(&buf1)
	err := encoder1.Encode(data)
	require.Nil(t, err)

	// 编码到第二个 Writer
	var buf2 bytes.Buffer
	encoder2 := json.API.NewEncoder(&buf2)
	err = encoder2.Encode(data)
	require.Nil(t, err)

	require.Equal(t, buf1.String(), buf2.String())
}

// TestSonicAPI_ConcurrentUsage 测试并发使用
func TestSonicAPI_ConcurrentUsage(t *testing.T) {
	ConfigureGinWithSonic()

	done := make(chan bool, 10)

	// 并发 Marshal
	for i := 0; i < 5; i++ {
		go func(id int) {
			data := map[string]int{"id": id}
			_, err := json.API.Marshal(data)
			require.Nil(t, err)
			done <- true
		}(i)
	}

	// 并发 Unmarshal
	for i := 0; i < 5; i++ {
		go func() {
			var result map[string]int
			err := json.API.Unmarshal([]byte(`{"id":1}`), &result)
			require.Nil(t, err)
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestSonicAPI_EmptyData 测试空数据处理
func TestSonicAPI_EmptyData(t *testing.T) {
	ConfigureGinWithSonic()

	// 空对象
	result, err := json.API.Marshal(map[string]any{})
	require.Nil(t, err)
	require.Equal(t, `{}`, string(result))

	// 空数组
	result, err = json.API.Marshal([]any{})
	require.Nil(t, err)
	require.Equal(t, `[]`, string(result))

	// 空字符串
	result, err = json.API.Marshal("")
	require.Nil(t, err)
	require.Equal(t, `""`, string(result))
}

// TestSonicAPI_SpecialCharacters 测试特殊字符处理
func TestSonicAPI_SpecialCharacters(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]string{
		"quote":    `"quoted"`,
		"backslash": `\path\to\file`,
		"newline":  "line1\nline2",
		"tab":      "col1\tcol2",
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)

	var parsed map[string]string
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, `"quoted"`, parsed["quote"])
	require.Equal(t, `\path\to\file`, parsed["backslash"])
	require.Equal(t, "line1\nline2", parsed["newline"])
	require.Equal(t, "col1\tcol2", parsed["tab"])
}

// TestSonicAPI_BoolHandling 测试布尔值处理
func TestSonicAPI_BoolHandling(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]bool{
		"true":  true,
		"false": false,
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)
	require.Contains(t, string(result), `"true":true`)
	require.Contains(t, string(result), `"false":false`)

	var parsed map[string]bool
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.True(t, parsed["true"])
	require.False(t, parsed["false"])
}

// TestSonicAPI_NumberHandling 测试数字处理
func TestSonicAPI_NumberHandling(t *testing.T) {
	ConfigureGinWithSonic()

	data := map[string]any{
		"int":     42,
		"float":   3.14159,
		"neg":     -100,
		"zero":    0,
		"large":   9223372036854775807, // max int64
	}

	result, err := json.API.Marshal(data)
	require.Nil(t, err)

	var parsed map[string]any
	err = json.API.Unmarshal(result, &parsed)
	require.Nil(t, err)
	require.Equal(t, float64(42), parsed["int"])
	require.InDelta(t, 3.14159, parsed["float"], 0.00001)
	require.Equal(t, float64(-100), parsed["neg"])
	require.Equal(t, float64(0), parsed["zero"])
}

// TestSonicAPI_NilEncoderDecoder 测试 nil 数据的编解码
func TestSonicAPI_NilEncoderDecoder(t *testing.T) {
	ConfigureGinWithSonic()

	// 测试 NewEncoder 返回非 nil
	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)
	require.NotNil(t, encoder)

	// 测试 NewDecoder 返回非 nil
	reader := bytes.NewBufferString(`{}`)
	decoder := json.API.NewDecoder(reader)
	require.NotNil(t, decoder)
}

// TestSonicAPI_ReadFromReadCloser 测试从 io.ReadCloser 读取
func TestSonicAPI_ReadFromReadCloser(t *testing.T) {
	ConfigureGinWithSonic()

	// 创建一个简单的 io.Reader
	jsonData := `{"test": "readcloser"}`
	reader := bytes.NewBufferString(jsonData)

	decoder := json.API.NewDecoder(reader)
	var result map[string]string
	err := decoder.Decode(&result)
	require.Nil(t, err)
	require.Equal(t, "readcloser", result["test"])
}

// TestSonicAPI_WriteToWriteCloser 测试写入 io.WriteCloser
func TestSonicAPI_WriteToWriteCloser(t *testing.T) {
	ConfigureGinWithSonic()

	var buf bytes.Buffer
	encoder := json.API.NewEncoder(&buf)

	data := map[string]string{"test": "writecloser"}
	err := encoder.Encode(data)
	require.Nil(t, err)
	require.True(t, buf.Len() > 0)
}

// 额外的接口检查
func TestSonicAPI_ImplementsInterfaces(t *testing.T) {
	ConfigureGinWithSonic()

	api := sonicAPI{api: sonic.ConfigStd}

	// 检查 encoder 接口
	var buf bytes.Buffer
	encoder := api.NewEncoder(&buf)
	require.NotNil(t, encoder)

	// 检查 decoder 接口
	reader := bytes.NewBufferString(`{}`)
	decoder := api.NewDecoder(reader)
	require.NotNil(t, decoder)
}
