package qlang

import (
	"fmt"
)

// RecoverAsError 在 defer 中使用，将 panic 转换为 error
// 返回的 error 需要作为命名返回值使用
// 用法: defer func() { err = comm.RecoverAsError(recover()) }()
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
// 用法: defer func() { comm.RecoverAndLog(recover(), logger, "operation name") }()
func RecoverAndLog(r any, logger Logger, operation string) {
	if r != nil {
		logger.Error(r).Str("operation", operation).Msg("goroutine panic recovered")
	}
}
