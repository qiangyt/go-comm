package comm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoverAsError_Nil(t *testing.T) {
	// recover() 返回 nil 时，RecoverAsError 应该返回 nil
	err := RecoverAsError(nil)
	assert.Nil(t, err)
}

func TestRecoverAsError_ErrorType(t *testing.T) {
	// panic 的值是 error 类型时，应该返回该 error
	originalErr := fmt.Errorf("test error")
	err := RecoverAsError(originalErr)
	assert.Equal(t, originalErr, err)
}

func TestRecoverAsError_NonErrorType(t *testing.T) {
	// panic 的值不是 error 类型时，应该转换为 fmt.Errorf
	err := RecoverAsError("string panic")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "string panic")
}

func TestRecoverAsError_IntType(t *testing.T) {
	// panic 的值是 int 类型时
	err := RecoverAsError(123)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "123")
}

func TestRecoverAndLog_Nil(t *testing.T) {
	// recover() 返回 nil 时，不应该记录日志
	logger := NewDiscardLogger()
	RecoverAndLog(nil, logger, "test-operation")
	// 不应该 panic
}

func TestRecoverAndLog_WithError(t *testing.T) {
	// panic 的值是 error 时，应该记录日志
	logger := NewDiscardLogger()
	RecoverAndLog(fmt.Errorf("test error"), logger, "test-operation")
	// 不应该 panic
}

func TestRecoverAndLog_WithString(t *testing.T) {
	// panic 的值是 string 时，应该记录日志
	logger := NewDiscardLogger()
	RecoverAndLog("string panic", logger, "test-operation")
	// 不应该 panic
}

func TestRecoverAsError_WithPanic(t *testing.T) {
	// 在实际的 panic/recover 场景中使用
	var err error
	func() {
		defer func() {
			err = RecoverAsError(recover())
		}()
		panic(fmt.Errorf("panic error"))
	}()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "panic error")
}

func TestRecoverAsError_WithSystemError(t *testing.T) {
	// 测试 SystemError 类型
	originalErr := NewSystemError("system error", nil)
	var err error
	func() {
		defer func() {
			err = RecoverAsError(recover())
		}()
		panic(originalErr)
	}()

	assert.NotNil(t, err)
	// 应该是原始的 SystemError
	assert.Equal(t, originalErr, err)
}
