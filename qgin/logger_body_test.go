package qgin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ==================== isTextContentType ====================

func TestIsTextContentType_JSON(t *testing.T) {
	a := require.New(t)
	a.True(isTextContentType("application/json"))
	a.True(isTextContentType("application/json; charset=utf-8"))
}

func TestIsTextContentType_XML(t *testing.T) {
	a := require.New(t)
	a.True(isTextContentType("application/xml"))
	a.True(isTextContentType("text/xml"))
	a.True(isTextContentType("application/xml; charset=utf-8"))
}

func TestIsTextContentType_Text(t *testing.T) {
	a := require.New(t)
	a.True(isTextContentType("text/plain"))
	a.True(isTextContentType("text/html"))
	a.True(isTextContentType("text/event-stream"))
	a.True(isTextContentType("text/css"))
}

func TestIsTextContentType_Form(t *testing.T) {
	a := require.New(t)
	a.True(isTextContentType("application/x-www-form-urlencoded"))
}

func TestIsTextContentType_SSE(t *testing.T) {
	a := require.New(t)
	a.True(isTextContentType("text/event-stream"))
}

func TestIsTextContentType_Binary(t *testing.T) {
	a := require.New(t)
	a.False(isTextContentType("image/png"))
	a.False(isTextContentType("image/jpeg"))
	a.False(isTextContentType("video/mp4"))
	a.False(isTextContentType("audio/mpeg"))
	a.False(isTextContentType("application/octet-stream"))
	a.False(isTextContentType("multipart/form-data"))
}

func TestIsTextContentType_Empty(t *testing.T) {
	a := require.New(t)
	a.False(isTextContentType(""))
}

func TestIsTextContentType_Other(t *testing.T) {
	a := require.New(t)
	a.False(isTextContentType("application/pdf"))
	a.False(isTextContentType("application/zip"))
}

func TestIsTextContentType_InvalidMediaType(t *testing.T) {
	a := require.New(t)
	// 测试解析失败的 content-type，会走简单匹配分支
	// 带有无效字符的 media type
	a.True(isTextContentType("invalid json"))  // 包含 json
	a.True(isTextContentType("invalid xml"))   // 包含 xml
	a.True(isTextContentType("text/invalid"))  // 以 text/ 开头
	a.False(isTextContentType("completely invalid")) // 不匹配任何条件
}

// ==================== maskAll ====================

func TestMaskAll(t *testing.T) {
	a := require.New(t)
	a.Equal("****", maskAll("secret"))
	a.Equal("****", maskAll(""))
	a.Equal("****", maskAll("Bearer eyJhbGci..."))
}

// ==================== maskHead ====================

func TestMaskHead_Short(t *testing.T) {
	a := require.New(t)
	a.Equal("****", maskHead("abc", 4))
	a.Equal("****", maskHead("", 4))
}

func TestMaskHead_Long(t *testing.T) {
	a := require.New(t)
	// "Bearer eyJhbGci..." 有 18 个字符，mask 前 4 个字符
	a.Equal("****er eyJhbGci...", maskHead("Bearer eyJhbGci...", 4))
}

func TestMaskHead_Zero(t *testing.T) {
	a := require.New(t)
	a.Equal("****", maskHead("Bearer", 0))
}

// ==================== maskTail ====================

func TestMaskTail_Short(t *testing.T) {
	a := require.New(t)
	a.Equal("****", maskTail("abc", 4))
	a.Equal("****", maskTail("", 4))
}

func TestMaskTail_Long(t *testing.T) {
	a := require.New(t)
	// "Bearer eyJhbGci..." 有 18 个字符，mask 后 4 个字符
	a.Equal("Bearer eyJhbGc****", maskTail("Bearer eyJhbGci...", 4))
}

func TestMaskTail_Zero(t *testing.T) {
	a := require.New(t)
	a.Equal("****", maskTail("Bearer", 0))
}

// ==================== applySensitiveStrategy ====================

func TestApplySensitiveStrategy_Full(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderFull}
	result := applySensitiveStrategy(value, cfg)

	a.Equal("Bearer token123", result)
}

func TestApplySensitiveStrategy_Exclude(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderExclude}
	result := applySensitiveStrategy(value, cfg)

	a.Equal("", result)
}

func TestApplySensitiveStrategy_MaskAll(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderMaskAll}
	result := applySensitiveStrategy(value, cfg)

	a.Equal("****", result)
}

func TestApplySensitiveStrategy_MaskHead(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderMaskHead, MaskSize: 4}
	result := applySensitiveStrategy(value, cfg)

	a.Equal("****er token123", result)
}

func TestApplySensitiveStrategy_MaskTail(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderMaskTail, MaskSize: 4}
	result := applySensitiveStrategy(value, cfg)

	a.Equal("Bearer toke****", result)
}

func TestApplySensitiveStrategy_DefaultMaskSize(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderMaskHead, MaskSize: 0}
	result := applySensitiveStrategy(value, cfg)

	// MaskSize 为 0 时应使用默认值 4
	a.Equal("****er token123", result)
}

func TestApplySensitiveStrategy_UnknownStrategy(t *testing.T) {
	a := require.New(t)

	value := "Bearer token123"
	cfg := SensitiveHeaderConfig{Strategy: SensitiveHeaderStrategy(999)}
	result := applySensitiveStrategy(value, cfg)

	// 未知策略应该使用默认值 MaskAll
	a.Equal("****", result)
}

// ==================== filterHeaders ====================

func TestFilterHeaders_None(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{"Content-Type": {"application/json"}}
	cfg := HeaderLogConfig{Strategy: HeaderLogNone}
	result := filterHeaders(headers, cfg)

	a.Nil(result)
}

func TestFilterHeaders_All(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Custom":     {"value"},
	}
	cfg := HeaderLogConfig{
		Strategy: HeaderLogAll,
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderFull,
			SensitiveList: []string{},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Equal("application/json", result["Content-Type"])
	a.Equal("value", result["X-Custom"])
}

func TestFilterHeaders_Whitelist(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Custom":     {"value"},
		"X-Other":      {"ignored"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogWhitelist,
		HeaderList: []string{"Content-Type", "X-Custom"},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Equal("application/json", result["Content-Type"])
	a.Equal("value", result["X-Custom"])
	a.NotContains(result, "X-Other")
}

func TestFilterHeaders_Blacklist(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Content-Type": {"application/json"},
		"X-Custom":     {"value"},
		"X-Other":      {"ignored"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogBlacklist,
		HeaderList: []string{"X-Other"},
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderFull,
			SensitiveList: []string{},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Equal("application/json", result["Content-Type"])
	a.Equal("value", result["X-Custom"])
	a.NotContains(result, "X-Other")
}

func TestFilterHeaders_WithSensitiveMask(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Authorization": {"Bearer token123"},
		"X-Custom":      {"value"},
	}
	cfg := HeaderLogConfig{
		Strategy: HeaderLogAll,
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderMaskAll,
			SensitiveList: []string{"Authorization"},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Equal("****", result["Authorization"]) // 已 mask
	a.Equal("value", result["X-Custom"])
}

func TestFilterHeaders_MultipleValues(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Set-Cookie": {"session=abc", "user=john"},
	}
	cfg := HeaderLogConfig{
		Strategy: HeaderLogAll,
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderFull,
			SensitiveList: []string{},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	// 多个值应该合并
	a.Contains(result["Set-Cookie"], "session=abc")
}

func TestFilterHeaders_WhitelistEmpty(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogWhitelist,
		HeaderList: []string{}, // 空白名单
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Empty(result) // 空结果
}

func TestFilterHeaders_BlacklistEmpty(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Content-Type": {"application/json"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogBlacklist,
		HeaderList: []string{}, // 空黑名单，记录所有
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderFull,
			SensitiveList: []string{},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.Equal("application/json", result["Content-Type"])
}

func TestFilterHeaders_BlacklistEmptyWithSensitiveExclude(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Authorization": {"Bearer token"},
		"Content-Type":  {"application/json"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogBlacklist,
		HeaderList: []string{}, // 空黑名单
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderExclude, // Exclude 策略，跳过敏感 header
			SensitiveList: []string{"Authorization"},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.NotContains(result, "Authorization") // Exclude 策略跳过
	a.Equal("application/json", result["Content-Type"])
}

func TestFilterHeaders_BlacklistWithSensitiveExclude(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Authorization": {"Bearer token"},
		"Content-Type":  {"application/json"},
	}
	cfg := HeaderLogConfig{
		Strategy:   HeaderLogBlacklist,
		HeaderList: []string{"X-Other"}, // 黑名单中不包含这些
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderExclude, // Exclude 策略，跳过敏感 header
			SensitiveList: []string{"Authorization"},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.NotContains(result, "Authorization") // Exclude 策略跳过
	a.Equal("application/json", result["Content-Type"])
}

func TestFilterHeaders_AllWithSensitiveExclude(t *testing.T) {
	a := require.New(t)

	headers := map[string][]string{
		"Authorization": {"Bearer token"},
		"Content-Type":  {"application/json"},
	}
	cfg := HeaderLogConfig{
		Strategy: HeaderLogAll,
		SensitiveConfig: &SensitiveHeaderConfig{
			Strategy:      SensitiveHeaderExclude, // Exclude 策略，跳过敏感 header
			SensitiveList: []string{"Authorization"},
		},
	}
	result := filterHeaders(headers, cfg)

	a.NotNil(result)
	a.NotContains(result, "Authorization") // Exclude 策略跳过
	a.Equal("application/json", result["Content-Type"])
}

// ==================== isSensitiveHeader ====================

func TestIsSensitiveHeader_Authorization(t *testing.T) {
	a := require.New(t)
	cfg := DefaultSensitiveHeaderConfig()
	a.True(isSensitiveHeader("Authorization", cfg))
	a.True(isSensitiveHeader("authorization", cfg)) // 大小写不敏感
}

func TestIsSensitiveHeader_Custom(t *testing.T) {
	a := require.New(t)
	cfg := &SensitiveHeaderConfig{
		SensitiveList: []string{"X-Custom-Secret"},
	}
	a.True(isSensitiveHeader("X-Custom-Secret", cfg))
	a.False(isSensitiveHeader("X-Public", cfg))
}
