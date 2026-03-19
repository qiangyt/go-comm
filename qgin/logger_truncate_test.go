package qgin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ==================== truncateHead ====================

func TestTruncateHead_ShortString(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateHead(s, 10)

	a.Equal("hello", result)
}

func TestTruncateHead_ExactLength(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateHead(s, 5)

	a.Equal("hello", result)
}

func TestTruncateHead_LongString(t *testing.T) {
	a := require.New(t)

	s := "hello world, this is a test"
	result := truncateHead(s, 5)

	a.Equal("hello...(truncated)", result)
}

func TestTruncateHead_EmptyString(t *testing.T) {
	a := require.New(t)

	s := ""
	result := truncateHead(s, 10)

	a.Equal("", result)
}

func TestTruncateHead_ZeroSize(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateHead(s, 0)

	a.Equal("...(truncated)", result)
}

// ==================== truncateTail ====================

func TestTruncateTail_ShortString(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateTail(s, 10)

	a.Equal("hello", result)
}

func TestTruncateTail_ExactLength(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateTail(s, 5)

	a.Equal("hello", result)
}

func TestTruncateTail_LongString(t *testing.T) {
	a := require.New(t)

	s := "hello world, this is a test"
	result := truncateTail(s, 4)

	a.Equal("...(truncated)test", result)
}

func TestTruncateTail_EmptyString(t *testing.T) {
	a := require.New(t)

	s := ""
	result := truncateTail(s, 10)

	a.Equal("", result)
}

func TestTruncateTail_ZeroSize(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateTail(s, 0)

	a.Equal("...(truncated)", result)
}

// ==================== truncateHeadAndTail ====================

func TestTruncateHeadAndTail_ShortString(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateHeadAndTail(s, 10)

	// len(s) <= n*2，返回完整内容
	a.Equal("hello", result)
}

func TestTruncateHeadAndTail_ExactBoundary(t *testing.T) {
	a := require.New(t)

	s := "helloworld" // len = 10
	result := truncateHeadAndTail(s, 5)

	// len(s) == n*2，返回完整内容（避免 overlap）
	a.Equal("helloworld", result)
}

func TestTruncateHeadAndTail_LongString(t *testing.T) {
	a := require.New(t)

	s := "hello world, this is a test string" // 34 chars
	result := truncateHeadAndTail(s, 5)

	// len(s) > n*2 (34 > 10)，截取前后各 5 个字符
	a.Equal("hello...(truncated)...tring", result)
}

func TestTruncateHeadAndTail_EmptyString(t *testing.T) {
	a := require.New(t)

	s := ""
	result := truncateHeadAndTail(s, 10)

	a.Equal("", result)
}

func TestTruncateHeadAndTail_ZeroSize(t *testing.T) {
	a := require.New(t)

	s := "hello"
	result := truncateHeadAndTail(s, 0)

	a.Equal("...(truncated)...", result)
}

func TestTruncateHeadAndTail_SingleCharOverlap(t *testing.T) {
	a := require.New(t)

	s := "hello world" // len = 11
	result := truncateHeadAndTail(s, 5)

	// len(s) > n*2 (11 > 10)，截取前后
	a.Equal("hello...(truncated)...world", result)
}

// ==================== applyTruncateStrategy ====================

func TestApplyTruncateStrategy_None(t *testing.T) {
	a := require.New(t)

	s := "hello world"
	cfg := BodyLogConfig{Strategy: BodyTruncateNone, TruncateSize: 5}
	result := applyTruncateStrategy(s, cfg)

	a.Equal("", result)
}

func TestApplyTruncateStrategy_Full(t *testing.T) {
	a := require.New(t)

	s := "hello world"
	cfg := BodyLogConfig{Strategy: BodyTruncateFull, TruncateSize: 5}
	result := applyTruncateStrategy(s, cfg)

	a.Equal("hello world", result)
}

func TestApplyTruncateStrategy_Head(t *testing.T) {
	a := require.New(t)

	s := "hello world"
	cfg := BodyLogConfig{Strategy: BodyTruncateHead, TruncateSize: 5}
	result := applyTruncateStrategy(s, cfg)

	a.Equal("hello...(truncated)", result)
}

func TestApplyTruncateStrategy_Tail(t *testing.T) {
	a := require.New(t)

	s := "hello world"
	cfg := BodyLogConfig{Strategy: BodyTruncateTail, TruncateSize: 5}
	result := applyTruncateStrategy(s, cfg)

	a.Equal("...(truncated)world", result)
}

func TestApplyTruncateStrategy_HeadAndTail(t *testing.T) {
	a := require.New(t)

	s := "hello world, this is a test" // 27 chars
	cfg := BodyLogConfig{Strategy: BodyTruncateHeadAndTail, TruncateSize: 5}
	result := applyTruncateStrategy(s, cfg)

	// 后 5 个字符是 " test"（包含空格）
	a.Equal("hello...(truncated)... test", result)
}

func TestApplyTruncateStrategy_EmptyBody(t *testing.T) {
	a := require.New(t)

	s := ""
	cfg := DefaultBodyLogConfig()
	result := applyTruncateStrategy(s, cfg)

	a.Equal("", result)
}

func TestApplyTruncateStrategy_DefaultTruncateSize(t *testing.T) {
	a := require.New(t)

	// 创建一个足够长的字符串
	s := make([]byte, 3000)
	for i := range s {
		s[i] = 'a'
	}

	cfg := BodyLogConfig{Strategy: BodyTruncateHeadAndTail, TruncateSize: 0}
	// TruncateSize 为 0 时应使用默认值 1024
	result := applyTruncateStrategy(string(s), cfg)

	// 应该被截取
	a.Contains(result, "...(truncated)...")
}

func TestApplyTruncateStrategy_DefaultStrategy(t *testing.T) {
	a := require.New(t)

	// 创建一个足够长的字符串
	s := make([]byte, 3000)
	for i := range s {
		s[i] = 'a'
	}

	cfg := BodyLogConfig{Strategy: BodyTruncateStrategy(999), TruncateSize: 5}
	// 未知策略应该使用默认值 HeadAndTail
	result := applyTruncateStrategy(string(s), cfg)

	// 应该被截取
	a.Contains(result, "...(truncated)...")
}

// ==================== formatTruncatedSize ====================

func TestFormatTruncatedSize(t *testing.T) {
	a := require.New(t)

	result := formatTruncatedSize(1000, 100)
	a.Equal("original: 1000 bytes, truncated: 100 bytes", result)
}
