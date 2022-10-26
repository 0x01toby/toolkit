package errs

import (
	"errors"
	"fmt"
	"runtime"
)

type WError interface {
	Wrap(err WError) WError
	Error() string
	String() string
	Err() error
	Next() WError
}

type WrapError struct {
	// error type
	error
	// err information
	Info   string
	Caller string
	Stack  WError
}

func New(err error, info string) WError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return &WrapError{Info: info, error: err, Caller: fmt.Sprintf("%s: %d", file, line)}
	}
	return &WrapError{Info: info, error: err, Caller: ""}
}

func (w *WrapError) Err() error {
	return w.error
}

func (w *WrapError) Error() string {
	if w.Stack != nil {
		return fmt.Sprintf("\t[%v]: %s, caller: %s \n%s", w.error, w.Info, w.Caller, w.Stack.Error())
	}
	return fmt.Sprintf("[%v]: %s, caller: %s", w.error, w.Info, w.Caller)
}

func (w *WrapError) String() string {
	return w.Error()
}

func (w *WrapError) Next() WError {
	return w.Stack
}

func (w *WrapError) Wrap(e WError) WError {
	w.Stack = new(WrapError)
	w.Stack = e
	return w
}

func Is(err WError, target error) bool {
	find := err
	for {
		if errors.Is(find.Err(), target) {
			return true
		}
		find = find.Next()
		if find == nil {
			break
		}
	}
	return false
}
