package comm

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// ==================== FromJson 覆盖率测试 ====================

func TestFromJson_withEnvsubst(t *testing.T) {
	a := require.New(t)

	// 设置环境变量
	t.Setenv("TEST_NAME", "env_value")

	jsonText := `{"name": "${TEST_NAME}", "age": 30}`
	type Config struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var result Config
	err := FromJson(jsonText, true, &result)
	a.NoError(err)
	a.Equal("env_value", result.Name)
	a.Equal(30, result.Age)
}

func TestFromJson_withEnvsubst_error(t *testing.T) {
	// 测试 envsubst 错误路径
	// 由于很难触发 EnvSubst 的错误（它依赖系统环境变量），
	// 我们跳过这个测试，但保留注释说明需要覆盖的分支
	t.Skip("难以触发 EnvSubst 错误分支")
}

func TestFromJson_invalidJson(t *testing.T) {
	a := require.New(t)

	invalidJson := `{invalid json}`
	type Config struct {
		Name string `json:"name"`
	}

	var result Config
	err := FromJson(invalidJson, false, &result)
	a.Error(err)
	a.Contains(err.Error(), "parse json")
}

func TestFromJson_emptyResult(t *testing.T) {
	a := require.New(t)

	jsonText := `{"name": "test"}`
	err := FromJson(jsonText, false, nil)
	a.Error(err) // unmarshal to nil should fail
}

func TestFromJsonP_panicsOnError(t *testing.T) {
	a := require.New(t)

	defer func() {
		r := recover()
		a.NotNil(r)
		a.Contains(r.(error).Error(), "parse json")
	}()

	invalidJson := `{invalid}`
	FromJsonP(invalidJson, false, &map[string]any{})
}

func TestFromJsonP_withEnvsubst(t *testing.T) {
	a := require.New(t)

	t.Setenv("MY_VAR", "substituted")

	jsonText := `{"value": "${MY_VAR}"}`
	var result map[string]any
	FromJsonP(jsonText, true, &result)
	a.NotNil(result)
	a.Equal("substituted", result["value"])
}

// ==================== MapFromJson 覆盖率测试 ====================

func TestMapFromJson_invalidJson(t *testing.T) {
	a := require.New(t)

	_, err := MapFromJson(`{invalid}`, false)
	a.Error(err)
}

func TestMapFromJsonP_panicsOnError(t *testing.T) {
	a := require.New(t)

	defer func() {
		r := recover()
		a.NotNil(r)
	}()

	MapFromJsonP(`{invalid}`, false)
}

func TestMapFromJson_withEnvsubst(t *testing.T) {
	a := require.New(t)

	t.Setenv("MAP_VAR", "map_value")

	result, err := MapFromJson(`{"key": "${MAP_VAR}"}`, true)
	a.NoError(err)
	a.Equal("map_value", result["key"])
}

// ==================== 并发性测试 ====================

func TestFromJson_concurrent(t *testing.T) {
	a := require.New(t)

	jsonText := `{"name": "test", "value": 123}`
	var wg sync.WaitGroup
	concurrency := 16

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var result map[string]any
			err := FromJson(jsonText, false, &result)
			a.NoError(err)
			a.Equal("test", result["name"])
		}()
	}
	wg.Wait()
}

func TestFromJsonP_concurrent(t *testing.T) {
	a := require.New(t)

	jsonText := `{"name": "test", "value": 123}`
	var wg sync.WaitGroup
	concurrency := 16

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var result map[string]any
			FromJsonP(jsonText, false, &result)
			a.NotNil(result)
		}()
	}
	wg.Wait()
}

// ==================== 扩展性测试 ====================

func TestFromJson_largeJson(t *testing.T) {
	a := require.New(t)

	// 创建一个大的 JSON 对象
	largeJson := `{"items": [`
	for i := 0; i < 1000; i++ {
		if i > 0 {
			largeJson += ","
		}
		largeJson += `{"id": ` + string(rune('0'+i%10)) + `, "name": "item` + string(rune('0'+i%10)) + `"}`
	}
	largeJson += `]}`

	var result map[string]any
	err := FromJson(largeJson, false, &result)
	a.NoError(err)
	a.NotNil(result["items"])
}

func TestFromJson_deeplyNested(t *testing.T) {
	a := require.New(t)

	// 创建深度嵌套的 JSON
	nestedJson := `{"level1": {"level2": {"level3": {"level4": {"level5": {"value": "deep"}}}}}}`

	var result map[string]any
	err := FromJson(nestedJson, false, &result)
	a.NoError(err)

	// 验证嵌套结构
	l1 := result["level1"].(map[string]any)
	l2 := l1["level2"].(map[string]any)
	l3 := l2["level3"].(map[string]any)
	l4 := l3["level4"].(map[string]any)
	l5 := l4["level5"].(map[string]any)
	a.Equal("deep", l5["value"])
}

func TestFromJson_emptyObject(t *testing.T) {
	a := require.New(t)

	var result map[string]any
	err := FromJson(`{}`, false, &result)
	a.NoError(err)
	a.Empty(result)
}

func TestFromJson_emptyArray(t *testing.T) {
	a := require.New(t)

	var result []any
	err := FromJson(`[]`, false, &result)
	a.NoError(err)
	a.Empty(result)
}

func TestFromJson_null(t *testing.T) {
	a := require.New(t)

	var result any
	err := FromJson(`null`, false, &result)
	a.NoError(err)
	a.Nil(result)
}

func TestFromJson_unicode(t *testing.T) {
	a := require.New(t)

	jsonText := `{"chinese": "中文测试", "emoji": "🎉", "mixed": "Hello世界"}`
	var result map[string]any
	err := FromJson(jsonText, false, &result)
	a.NoError(err)
	a.Equal("中文测试", result["chinese"])
	a.Equal("🎉", result["emoji"])
	a.Equal("Hello世界", result["mixed"])
}

// ==================== 安全性测试 ====================

func TestFromJson_injectionAttempt(t *testing.T) {
	// 尝试 JSON 注入
	maliciousJson := `{"name": "test\"}", "extra": "injected"}`
	var result map[string]any
	err := FromJson(maliciousJson, false, &result)
	// 应该要么正确解析，要么返回错误，不应该崩溃
	_ = err // 可能是错误也可能是正确解析
}

func TestFromJson_veryLongString(t *testing.T) {
	a := require.New(t)

	// 非常长的字符串
	longStr := `{"value": "` + makeLongString(10000) + `"}`
	var result map[string]any
	err := FromJson(longStr, false, &result)
	a.NoError(err)
}

func TestFromJson_specialCharacters(t *testing.T) {
	a := require.New(t)

	jsonText := `{"newline": "line1\nline2", "tab": "col1\tcol2", "quote": "say \"hello\""}`
	var result map[string]any
	err := FromJson(jsonText, false, &result)
	a.NoError(err)
	a.Contains(result["newline"].(string), "line1")
	a.Contains(result["tab"].(string), "col1")
}

func makeLongString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = 'a' + byte(i%26)
	}
	return string(result)
}

func TestFromJsonP_happy(t *testing.T) {
	a := require.New(t)

	type Config struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonText := `{"name": "test", "age": 30}`
	var result Config
	FromJsonP(jsonText, false, &result)
	a.Equal("test", result.Name)
	a.Equal(30, result.Age)
}

func TestMapFromJson_happy(t *testing.T) {
	a := require.New(t)

	jsonText := `{"key1": "value1", "key2": "value2"}`
	result, err := MapFromJson(jsonText, false)
	a.NoError(err)
	a.NotNil(result)
	a.Equal("value1", result["key1"])
	a.Equal("value2", result["key2"])
}

func TestMapFromJsonP_happy(t *testing.T) {
	a := require.New(t)

	jsonText := `{"key1": "value1", "key2": "value2"}`
	result := MapFromJsonP(jsonText, false)
	a.NotNil(result)
	a.Equal("value1", result["key1"])
	a.Equal("value2", result["key2"])
}
