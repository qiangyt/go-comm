package qerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigError(t *testing.T) {
	baseErr := errors.New("invalid value")

	err := NewConfigError("配置项无效", baseErr)

	assert.NotNil(t, err)
	assert.Equal(t, ErrCodeConfig, err.Code)
	assert.Equal(t, "配置项无效", err.Message)
	assert.Equal(t, baseErr, err.Err)
	assert.Contains(t, err.Error(), "配置项无效")
	assert.Contains(t, err.Error(), "invalid value")
}

func TestNewBusinessError(t *testing.T) {
	err := NewBusinessError("会话不存在: session-123", nil)

	assert.NotNil(t, err)
	assert.Equal(t, ErrCodeBusiness, err.Code)
	assert.Equal(t, "会话不存在: session-123", err.Message)
	assert.Nil(t, err.Err)
	assert.Contains(t, err.Error(), "会话不存在: session-123")
}

func TestNewSystemError(t *testing.T) {
	baseErr := errors.New("connection refused")

	err := NewSystemError("连接 NATS 失败", baseErr)

	assert.NotNil(t, err)
	assert.Equal(t, ErrCodeSystem, err.Code)
	assert.Equal(t, "连接 NATS 失败", err.Message)
	assert.Equal(t, baseErr, err.Err)
	assert.Contains(t, err.Error(), "连接 NATS 失败")
	assert.Contains(t, err.Error(), "connection refused")
}

func TestNewSecurityError(t *testing.T) {
	err := NewSecurityError("命令被安全策略阻止: rm -rf /", nil)

	assert.NotNil(t, err)
	assert.Equal(t, ErrCodeSecurity, err.Code)
	assert.Equal(t, "命令被安全策略阻止: rm -rf /", err.Message)
	assert.Nil(t, err.Err)
	assert.Contains(t, err.Error(), "命令被安全策略阻止")
}

func TestAppError_Error(t *testing.T) {
	baseErr := errors.New("underlying error")
	err := &AppError{
		Code:    ErrCodeConfig,
		Message: "配置错误",
		Err:     baseErr,
	}

	errStr := err.Error()

	assert.Contains(t, errStr, "配置错误")
	assert.Contains(t, errStr, "underlying error")
	assert.Contains(t, errStr, string(ErrCodeConfig))
}

func TestAppError_Unwrap(t *testing.T) {
	baseErr := errors.New("underlying error")
	err := &AppError{
		Code:    ErrCodeBusiness,
		Message: "业务错误",
		Err:     baseErr,
	}

	unwrapped := errors.Unwrap(err)

	assert.Equal(t, baseErr, unwrapped)
}

func TestAppError_Is(t *testing.T) {
	err1 := NewConfigError("错误1", nil)
	err2 := NewConfigError("错误2", nil)
	err3 := NewBusinessError("业务错误", nil)

	assert.True(t, errors.Is(err1, ErrCodeConfig))
	assert.False(t, errors.Is(err1, ErrCodeBusiness))
	assert.False(t, errors.Is(err1, err2))
	assert.False(t, errors.Is(err1, err3))
}

func TestAppError_As(t *testing.T) {
	baseErr := errors.New("original error")
	appErr := NewConfigError("配置错误", baseErr)

	var target *AppError
	ok := errors.As(appErr, &target)

	assert.True(t, ok)
	assert.Equal(t, ErrCodeConfig, target.Code)
	assert.Equal(t, "配置错误", target.Message)
}

func TestAppError_WithMessage(t *testing.T) {
	original := NewConfigError("原始错误", nil)

	updated := original.WithMessage("新错误: %s", "详情")

	assert.Equal(t, ErrCodeConfig, updated.Code)
	assert.Equal(t, "新错误: 详情", updated.Message)
	assert.Equal(t, original.Err, updated.Err)
}

func TestAppError_WithCode(t *testing.T) {
	original := NewBusinessError("业务错误", nil)

	updated := original.WithCode(ErrCodeSecurity)

	assert.Equal(t, ErrCodeSecurity, updated.Code)
	assert.Equal(t, "业务错误", updated.Message)
}

func TestErrCodeValues(t *testing.T) {
	assert.Equal(t, ErrCode("CFG"), ErrCodeConfig)
	assert.Equal(t, ErrCode("BIZ"), ErrCodeBusiness)
	assert.Equal(t, ErrCode("SYS"), ErrCodeSystem)
	assert.Equal(t, ErrCode("SEC"), ErrCodeSecurity)
}

func TestNewConfigErrorf(t *testing.T) {
	err := NewConfigErrorf("配置项 %s 无效", "port")

	assert.Equal(t, ErrCodeConfig, err.Code)
	assert.Equal(t, "配置项 port 无效", err.Message)
}

func TestNewBusinessErrorf(t *testing.T) {
	err := NewBusinessErrorf("会话 %s 不存在", "abc123")

	assert.Equal(t, ErrCodeBusiness, err.Code)
	assert.Equal(t, "会话 abc123 不存在", err.Message)
}

func TestNewSystemErrorf(t *testing.T) {
	err := NewSystemErrorf("连接 %s 失败", "NATS")

	assert.Equal(t, ErrCodeSystem, err.Code)
	assert.Equal(t, "连接 NATS 失败", err.Message)
}

func TestNewSecurityErrorf(t *testing.T) {
	err := NewSecurityErrorf("路径 %s 不在允许的目录内", "/etc/passwd")

	assert.Equal(t, ErrCodeSecurity, err.Code)
	assert.Equal(t, "路径 /etc/passwd 不在允许的目录内", err.Message)
}
