package qlang

import (
	"fmt"
)

// RecoverAsError 在 defer 中使用，将 panic 转换为 error
// 返回的 error 需要作为命名返回值使用
// 用法: defer func() { err = qlang.RecoverAsError(recover()) }()
func RecoverAsError(r any) error {
	if r == nil {
		return nil
	}

	switch v := r.(type) {
	case error:
		return v
	default:
		return fmt.Errorf("%v", v)
	}
}

// RecoverAndLog 在 goroutine 的 defer 中使用，记录 panic 日志
// 用法: defer func() { qlang.RecoverAndLog(recover(), logger, "operation name") }()
func RecoverAndLog(r any, logger Logger, operation string) error {
	err := RecoverAsError(r)
	if logger != nil {
		var log LogEntry
		if err != nil {
			log = logger.Error(err)
		} else {
			log = logger.Warn()
		}
		if operation != "" {
			log = log.Str("operation", operation)
		}
		log.Msg("panic recovered")
	}
	return err
}
