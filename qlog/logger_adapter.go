package qlog

import (
	"fmt"
)

// QlogLoggerAdapter 将 qlog.Logger 适配为 qgin.Logger 接口
type QlogLoggerAdapter struct {
	Logger Logger
}

// NewQlogLoggerAdapter 创建适配器
func NewQlogLoggerAdapter(logger Logger) *QlogLoggerAdapter {
	return &QlogLoggerAdapter{Logger: logger}
}

// parseFields 将 fields 解析为 map[string]any
// fields 格式: key1, value1, key2, value2, ...
func parseFields(fields []any) map[string]any {
	result := make(map[string]any)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			result[key] = fields[i+1]
		}
	}
	return result
}

// Info 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Info(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Info().Msg(msg)
	} else {
		a.Logger.Info().Fields(parseFields(fields)).Msg(msg)
	}
}

// Warn 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Warn(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Warn().Msg(msg)
	} else {
		a.Logger.Warn().Fields(parseFields(fields)).Msg(msg)
	}
}

// Error 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Error(msg string, fields ...any) {
	err := fmt.Errorf("error: %s", msg)
	if len(fields) == 0 {
		a.Logger.Error(err).Msg(msg)
	} else {
		a.Logger.Error(err).Fields(parseFields(fields)).Msg(msg)
	}
}

// Debug 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Debug(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Debug().Msg(msg)
	} else {
		a.Logger.Debug().Fields(parseFields(fields)).Msg(msg)
	}
}

// Trace 实现 qgin.Logger 接口
// phuslu/log 没有 Trace 级别，使用 Debug 代替
func (a *QlogLoggerAdapter) Trace(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Debug().Msg(msg)
	} else {
		a.Logger.Debug().Fields(parseFields(fields)).Msg(msg)
	}
}

// Fatal 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Fatal(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Fatal().Msg(msg)
	} else {
		a.Logger.Fatal().Fields(parseFields(fields)).Msg(msg)
	}
}

// Panic 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) Panic(msg string, fields ...any) {
	if len(fields) == 0 {
		a.Logger.Panic().Msg(msg)
	} else {
		a.Logger.Panic().Fields(parseFields(fields)).Msg(msg)
	}
}

// WithField 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) WithField(key string, value any) any {
	return &fieldWrapper{adapter: a, fields: map[string]any{key: value}}
}

// WithFields 实现 qgin.Logger 接口
func (a *QlogLoggerAdapter) WithFields(fields map[string]any) any {
	return &fieldWrapper{adapter: a, fields: fields}
}

// fieldWrapper 用于 WithField/WithFields 返回值
type fieldWrapper struct {
	adapter *QlogLoggerAdapter
	fields  map[string]any
}
