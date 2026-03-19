package qgin

import (
	"fmt"
)

const truncatedMarker = "...(truncated)"
const truncatedMiddleMarker = "...(truncated)..."

// truncateHead 截取字符串前 n 个字符
func truncateHead(s string, n int) string {
	if n <= 0 {
		return truncatedMarker
	}
	if len(s) <= n {
		return s
	}
	return s[:n] + truncatedMarker
}

// truncateTail 截取字符串后 n 个字符
func truncateTail(s string, n int) string {
	if n <= 0 {
		return truncatedMarker
	}
	if len(s) <= n {
		return s
	}
	return truncatedMarker + s[len(s)-n:]
}

// truncateHeadAndTail 截取字符串前后各 n 个字符
// 当 body 长度 <= n*2 时，返回完整 body（不做截取，避免重叠）
func truncateHeadAndTail(s string, n int) string {
	if n <= 0 {
		return truncatedMiddleMarker
	}
	if len(s) <= n*2 {
		// 避免 overlap，返回完整内容
		return s
	}
	return s[:n] + truncatedMiddleMarker + s[len(s)-n:]
}

// applyTruncateStrategy 根据配置应用截取策略
func applyTruncateStrategy(s string, cfg BodyLogConfig) string {
	// 空字符串直接返回
	if s == "" {
		return s
	}

	// 获取截取大小，默认 1024
	n := cfg.TruncateSize
	if n <= 0 {
		n = 1024
	}

	switch cfg.Strategy {
	case BodyTruncateNone:
		return ""
	case BodyTruncateFull:
		return s
	case BodyTruncateHead:
		return truncateHead(s, n)
	case BodyTruncateTail:
		return truncateTail(s, n)
	case BodyTruncateHeadAndTail:
		return truncateHeadAndTail(s, n)
	default:
		return truncateHeadAndTail(s, n)
	}
}

// formatTruncatedSize 格式化截取大小信息
func formatTruncatedSize(originalLen int, truncatedLen int) string {
	return fmt.Sprintf("original: %d bytes, truncated: %d bytes", originalLen, truncatedLen)
}
