package errors

import (
	"fmt"
	"runtime"
)

type ErrorType string

const (
	ConfigError     ErrorType = "ConfigError"
	UIError         ErrorType = "UIError"
	ProcessError    ErrorType = "ProcessError"
	MonitorError    ErrorType = "MonitorError"
	ValidationError ErrorType = "ValidationError"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
	Stack   string
}

func NewAppError(errType ErrorType, message string, err error) *AppError {
	stack := make([]byte, 4096)
	runtime.Stack(stack, false)

	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
		Stack:   string(stack),
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func IsType(err error, errType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errType
	}
	return false
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*AppError); ok {
		appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		return appErr
	}
	return NewAppError(MonitorError, message, err)
}
