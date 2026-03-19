package qgin

import (
	"mime"
	"strings"
)

// ==================== Content-Type 判断 ====================

// isTextContentType 判断 Content-Type 是否为文本类型
// 文本类型会记录 body 内容，二进制类型只记录类型和大小
func isTextContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	// 解析 media type
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		// 解析失败，尝试简单匹配
		return strings.Contains(contentType, "json") ||
			strings.Contains(contentType, "xml") ||
			strings.HasPrefix(contentType, "text/")
	}

	// JSON 类型
	if mediaType == "application/json" {
		return true
	}

	// XML 类型
	if mediaType == "application/xml" || mediaType == "text/xml" {
		return true
	}

	// text/* 类型
	if strings.HasPrefix(mediaType, "text/") {
		return true
	}

	// form-urlencoded
	if mediaType == "application/x-www-form-urlencoded" {
		return true
	}

	return false
}

// ==================== Mask 函数 ====================

// maskAll 将整个值替换为 ****
func maskAll(value string) string {
	return "****"
}

// maskHead mask 前 n 个字符
func maskHead(value string, n int) string {
	if n <= 0 {
		return "****"
	}
	if len(value) <= n {
		return "****"
	}
	return "****" + value[n:]
}

// maskTail mask 后 n 个字符
func maskTail(value string, n int) string {
	if n <= 0 {
		return "****"
	}
	if len(value) <= n {
		return "****"
	}
	return value[:len(value)-n] + "****"
}

// applySensitiveStrategy 根据配置应用敏感信息处理策略
func applySensitiveStrategy(value string, cfg SensitiveHeaderConfig) string {
	maskSize := cfg.MaskSize
	if maskSize <= 0 {
		maskSize = 4
	}

	switch cfg.Strategy {
	case SensitiveHeaderFull:
		return value
	case SensitiveHeaderExclude:
		return ""
	case SensitiveHeaderMaskAll:
		return maskAll(value)
	case SensitiveHeaderMaskHead:
		return maskHead(value, maskSize)
	case SensitiveHeaderMaskTail:
		return maskTail(value, maskSize)
	default:
		return maskAll(value)
	}
}

// ==================== Header 过滤 ====================

// isSensitiveHeader 判断 header 是否为敏感 header（大小写不敏感）
func isSensitiveHeader(headerName string, cfg *SensitiveHeaderConfig) bool {
	if cfg == nil || len(cfg.SensitiveList) == 0 {
		return false
	}

	name := strings.ToLower(headerName)
	for _, h := range cfg.SensitiveList {
		if strings.ToLower(h) == name {
			return true
		}
	}
	return false
}

// filterHeaders 根据 HeaderLogConfig 过滤 headers
// 返回 map[string]string（key -> 合并后的值）
func filterHeaders(headers map[string][]string, cfg HeaderLogConfig) map[string]string {
	if cfg.Strategy == HeaderLogNone {
		return nil
	}

	result := make(map[string]string)

	switch cfg.Strategy {
	case HeaderLogAll:
		for k, v := range headers {
			value := strings.Join(v, "; ")
			// 检查是否为敏感 header
			if isSensitiveHeader(k, cfg.SensitiveConfig) {
				value = applySensitiveStrategy(value, *cfg.SensitiveConfig)
				if value == "" {
					continue // Exclude 策略时跳过
				}
			}
			result[k] = value
		}

	case HeaderLogWhitelist:
		if len(cfg.HeaderList) == 0 {
			return result
		}
		whitelist := make(map[string]bool)
		for _, h := range cfg.HeaderList {
			whitelist[strings.ToLower(h)] = true
		}
		for k, v := range headers {
			if whitelist[strings.ToLower(k)] {
				result[k] = strings.Join(v, "; ")
			}
		}

	case HeaderLogBlacklist:
		if len(cfg.HeaderList) == 0 {
			// 没有黑名单，记录所有
			for k, v := range headers {
				value := strings.Join(v, "; ")
				if isSensitiveHeader(k, cfg.SensitiveConfig) {
					value = applySensitiveStrategy(value, *cfg.SensitiveConfig)
					if value == "" {
						continue
					}
				}
				result[k] = value
			}
			return result
		}
		blacklist := make(map[string]bool)
		for _, h := range cfg.HeaderList {
			blacklist[strings.ToLower(h)] = true
		}
		for k, v := range headers {
			if blacklist[strings.ToLower(k)] {
				continue
			}
			value := strings.Join(v, "; ")
			if isSensitiveHeader(k, cfg.SensitiveConfig) {
				value = applySensitiveStrategy(value, *cfg.SensitiveConfig)
				if value == "" {
					continue
				}
			}
			result[k] = value
		}
	}

	return result
}
