package qgin

import (
	"regexp"
	"strings"
)

// sseEventSplitRegex 匹配 SSE 事件分隔符（\n\n 或 \r\n\r\n）
var sseEventSplitRegex = regexp.MustCompile(`\r?\n\r?\n`)

// parseSSEEvents 解析 SSE 事件（以 \n\n 或 \r\n\r\n 分隔）
// 只返回完整的事件（末尾有分隔符的事件）
// 返回事件列表（不含空事件）
func parseSSEEvents(body string) []string {
	if body == "" {
		return nil
	}

	// 查找所有完整事件的结束位置
	var events []string
	lastEnd := 0

	// 使用正则表达式查找所有分隔符
	matches := sseEventSplitRegex.FindAllStringIndex(body, -1)

	for _, match := range matches {
		// 提取事件内容（不含分隔符）
		event := strings.TrimSpace(body[lastEnd:match[0]])
		if event != "" {
			events = append(events, event)
		}
		lastEnd = match[1]
	}

	return events
}

// truncateSSEEvents 根据 SSELogConfig 截取 SSE 事件
// 含 overlap 边界处理：当事件数 < n*2 时返回全部
func truncateSSEEvents(events []string, cfg SSELogConfig) string {
	if len(events) == 0 {
		return ""
	}

	// None 策略不记录
	if cfg.Strategy == SSETruncateNone {
		return ""
	}

	// 获取截取大小，默认 10
	n := cfg.TruncateSize
	if n <= 0 {
		n = 10
	}

	// Full 策略返回全部
	if cfg.Strategy == SSETruncateFull {
		return strings.Join(events, "\n\n")
	}

	// 事件数不足（< n*2），返回全部（避免 overlap）
	if len(events) < n*2 {
		return strings.Join(events, "\n\n")
	}

	switch cfg.Strategy {
	case SSETruncateHead:
		// 只取前 n 个事件
		selected := events[:n]
		return strings.Join(selected, "\n\n") + "\n\n...(truncated)"

	case SSETruncateTail:
		// 只取后 n 个事件
		selected := events[len(events)-n:]
		return "...(truncated)\n\n" + strings.Join(selected, "\n\n")

	case SSETruncateHeadAndTail:
		// 取前后各 n 个事件
		head := events[:n]
		tail := events[len(events)-n:]
		return strings.Join(head, "\n\n") + "\n\n...(truncated)...\n\n" + strings.Join(tail, "\n\n")

	default:
		// 默认使用 HeadAndTail
		head := events[:n]
		tail := events[len(events)-n:]
		return strings.Join(head, "\n\n") + "\n\n...(truncated)...\n\n" + strings.Join(tail, "\n\n")
	}
}
