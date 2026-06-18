package qgin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ==================== parseSSEEvents ====================

func TestParseSSEEvents_Single(t *testing.T) {
	a := require.New(t)

	body := "data: hello\n\n"
	events := parseSSEEvents(body)

	a.Len(events, 1)
	a.Equal("data: hello", events[0])
}

func TestParseSSEEvents_Multiple(t *testing.T) {
	a := require.New(t)

	body := "data: hello\n\ndata: world\n\n"
	events := parseSSEEvents(body)

	a.Len(events, 2)
	a.Equal("data: hello", events[0])
	a.Equal("data: world", events[1])
}

func TestParseSSEEvents_WithCRLF(t *testing.T) {
	a := require.New(t)

	body := "data: hello\r\n\r\ndata: world\r\n\r\n"
	events := parseSSEEvents(body)

	a.Len(events, 2)
	a.Equal("data: hello", events[0])
	a.Equal("data: world", events[1])
}

func TestParseSSEEvents_MixedLineEndings(t *testing.T) {
	a := require.New(t)

	body := "data: hello\n\ndata: world\r\n\r\n"
	events := parseSSEEvents(body)

	a.Len(events, 2)
	a.Equal("data: hello", events[0])
	a.Equal("data: world", events[1])
}

func TestParseSSEEvents_Empty(t *testing.T) {
	a := require.New(t)

	events := parseSSEEvents("")

	a.Len(events, 0)
}

func TestParseSSEEvents_Trailing(t *testing.T) {
	a := require.New(t)

	// 末尾没有完整的分隔符
	body := "data: hello\n\ndata: incomplete"
	events := parseSSEEvents(body)

	a.Len(events, 1)
	a.Equal("data: hello", events[0])
}

func TestParseSSEEvents_WithFields(t *testing.T) {
	a := require.New(t)

	body := "event: message\ndata: hello\nid: 123\n\n"
	events := parseSSEEvents(body)

	a.Len(events, 1)
	a.Contains(events[0], "event: message")
	a.Contains(events[0], "data: hello")
	a.Contains(events[0], "id: 123")
}

// ==================== truncateSSEEvents ====================

func TestTruncateSSEEvents_None(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3"}
	cfg := SSELogConfig{Strategy: SSETruncateNone, TruncateSize: 10}
	result := truncateSSEEvents(events, cfg)

	a.Equal("", result)
}

func TestTruncateSSEEvents_Full(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3"}
	cfg := SSELogConfig{Strategy: SSETruncateFull, TruncateSize: 10}
	result := truncateSSEEvents(events, cfg)

	a.Contains(result, "data: 1")
	a.Contains(result, "data: 2")
	a.Contains(result, "data: 3")
}

func TestTruncateSSEEvents_Head(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3", "data: 4", "data: 5"}
	cfg := SSELogConfig{Strategy: SSETruncateHead, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	a.Contains(result, "data: 1")
	a.Contains(result, "data: 2")
	a.NotContains(result, "data: 3")
	a.NotContains(result, "data: 4")
	a.NotContains(result, "data: 5")
	a.Contains(result, "...(truncated)")
}

func TestTruncateSSEEvents_Tail(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3", "data: 4", "data: 5"}
	cfg := SSELogConfig{Strategy: SSETruncateTail, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	a.NotContains(result, "data: 1")
	a.NotContains(result, "data: 2")
	a.NotContains(result, "data: 3")
	a.Contains(result, "data: 4")
	a.Contains(result, "data: 5")
	a.Contains(result, "...(truncated)")
}

func TestTruncateSSEEvents_HeadAndTail(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3", "data: 4", "data: 5"}
	cfg := SSELogConfig{Strategy: SSETruncateHeadAndTail, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	a.Contains(result, "data: 1")
	a.Contains(result, "data: 2")
	a.Contains(result, "data: 4")
	a.Contains(result, "data: 5")
	a.NotContains(result, "data: 3")
	a.Contains(result, "...(truncated)...")
}

func TestTruncateSSEEvents_HeadAndTailOverlap(t *testing.T) {
	a := require.New(t)

	// 事件数量 <= n*2 时，返回所有事件（避免 overlap）
	events := []string{"data: 1", "data: 2", "data: 3"}
	cfg := SSELogConfig{Strategy: SSETruncateHeadAndTail, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	// 3 个事件，n=2，n*2=4 > 3，所以返回所有
	a.Contains(result, "data: 1")
	a.Contains(result, "data: 2")
	a.Contains(result, "data: 3")
	a.NotContains(result, "...(truncated)...")
}

func TestTruncateSSEEvents_Empty(t *testing.T) {
	a := require.New(t)

	events := []string{}
	cfg := SSELogConfig{Strategy: SSETruncateHeadAndTail, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	a.Equal("", result)
}

func TestTruncateSSEEvents_DefaultSize(t *testing.T) {
	a := require.New(t)

	// 创建很多事件
	events := make([]string, 20)
	for i := 0; i < 20; i++ {
		events[i] = "data: event"
	}
	cfg := SSELogConfig{Strategy: SSETruncateHeadAndTail, TruncateSize: 0}
	result := truncateSSEEvents(events, cfg)

	// TruncateSize 为 0 时应使用默认值 10
	a.Contains(result, "...(truncated)...")
}

func TestTruncateSSEEvents_UnknownStrategy(t *testing.T) {
	a := require.New(t)

	events := []string{"data: 1", "data: 2", "data: 3", "data: 4", "data: 5"}
	cfg := SSELogConfig{Strategy: SSETruncateStrategy(999), TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	// 未知策略应该使用默认值 HeadAndTail
	a.Contains(result, "...(truncated)...")
}

func TestTruncateSSEEvents_HeadAndTailExactlyN2(t *testing.T) {
	a := require.New(t)

	// 事件数正好等于 n*2 + 1，应该截取
	events := []string{"data: 1", "data: 2", "data: 3", "data: 4", "data: 5"}
	cfg := SSELogConfig{Strategy: SSETruncateHeadAndTail, TruncateSize: 2}
	result := truncateSSEEvents(events, cfg)

	// 5 个事件，n=2，n*2=4 < 5，所以截取
	a.Contains(result, "data: 1")
	a.Contains(result, "data: 2")
	a.Contains(result, "data: 4")
	a.Contains(result, "data: 5")
	a.Contains(result, "...(truncated)...")
}
