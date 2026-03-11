package comm

import (
	"fmt"
)

type ErrCode string

func (me ErrCode) Error() string {
	return string(me)
}

const (
	ErrCodeConfig   ErrCode = "CFG"
	ErrCodeBusiness ErrCode = "BIZ"
	ErrCodeSystem   ErrCode = "SYS"
	ErrCodeSecurity ErrCode = "SEC"
)

type AppError struct {
	Code    ErrCode
	Message string
	Err     error
}

func (me *AppError) Error() string {
	if me.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", me.Code, me.Message, me.Err.Error())
	}
	return fmt.Sprintf("%s: %s", me.Code, me.Message)
}

func (me *AppError) Unwrap() error {
	return me.Err
}

func (me *AppError) Is(target error) bool {
	if code, ok := target.(ErrCode); ok {
		return me.Code == code
	}
	return false
}

func (me *AppError) As(target any) bool {
	if t, ok := target.(**AppError); ok {
		*t = me
		return true
	}
	return false
}

func (me *AppError) WithMessage(format string, args ...any) *AppError {
	return &AppError{
		Code:    me.Code,
		Message: fmt.Sprintf(format, args...),
		Err:     me.Err,
	}
}

func (me *AppError) WithCode(code ErrCode) *AppError {
	return &AppError{
		Code:    code,
		Message: me.Message,
		Err:     me.Err,
	}
}

func NewConfigError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeConfig,
		Message: message,
		Err:     err,
	}
}

func NewBusinessError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeBusiness,
		Message: message,
		Err:     err,
	}
}

func NewSystemError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeSystem,
		Message: message,
		Err:     err,
	}
}

func NewSecurityError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeSecurity,
		Message: message,
		Err:     err,
	}
}

func NewConfigErrorf(format string, args ...any) *AppError {
	return &AppError{
		Code:    ErrCodeConfig,
		Message: fmt.Sprintf(format, args...),
		Err:     nil,
	}
}

func NewBusinessErrorf(format string, args ...any) *AppError {
	return &AppError{
		Code:    ErrCodeBusiness,
		Message: fmt.Sprintf(format, args...),
		Err:     nil,
	}
}

func NewSystemErrorf(format string, args ...any) *AppError {
	return &AppError{
		Code:    ErrCodeSystem,
		Message: fmt.Sprintf(format, args...),
		Err:     nil,
	}
}

func NewSecurityErrorf(format string, args ...any) *AppError {
	return &AppError{
		Code:    ErrCodeSecurity,
		Message: fmt.Sprintf(format, args...),
		Err:     nil,
	}
}

var _ error = (*AppError)(nil)
