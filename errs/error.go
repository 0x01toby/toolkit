package errs

import (
	"errors"
	"fmt"
	"runtime"
)

type IError interface {
	// Wrap 包裹error
	Wrap(err IError) IError
	// Error 格式化输出string
	Error() string
	// String 输出为string
	String() string
	// BlankError 空白error
	BlankError() error
	// Next 下一个Error
	Next() IError
	// Is 判断当前error是否包含目标error
	Is(error) bool
}

// WrapError 包装Error
type WrapError struct {
	// error type
	error
	// err information
	Info   string
	Caller string
	Stack  IError
}

func New(err error, info string) IError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return &WrapError{Info: info, error: err, Caller: fmt.Sprintf("%s: %d", file, line)}
	}
	return &WrapError{Info: info, error: err, Caller: "", Stack: nil}
}

// BlankError 当前的error
func (w *WrapError) BlankError() error {
	return w.error
}

// Error 格式化输出string
func (w *WrapError) Error() string {
	if w.Stack != nil {
		return fmt.Sprintf("%s\t[%v]: %s, caller: %s \n", w.Stack.Error(), w.error, w.Info, w.Caller)
	}
	return fmt.Sprintf("[%v]: %s, caller: %s", w.error, w.Info, w.Caller)
}

// String 格式化输出string
func (w *WrapError) String() string {
	return w.Error()
}

// Next error stack
func (w *WrapError) Next() IError {
	return w.Stack
}

// Wrap 包裹error
func (w *WrapError) Wrap(e IError) IError {
	var err IError = w
	for {
		if err.Next() == nil {
			err.(*WrapError).Stack = new(WrapError)
			err.(*WrapError).Stack = e
			break
		}
		err = err.Next()
	}
	return w
}

// Is 判断是否包含目标target error
func (w *WrapError) Is(target error) bool {
	var err IError
	err = w
	for {
		if errors.Is(err.BlankError(), target) {
			return true
		}
		if err = err.Next(); err == nil {
			break
		}
	}
	return false
}
