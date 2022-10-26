package errors

import (
	"fmt"
)

type WrapperError struct {
	error
	Key string
}

func (w *WrapperError) Error() string {
	return fmt.Sprintf("%s: %v", w.Key, *w)
}

func Wrap(key string, err error) *WrapperError {
	return &WrapperError{Key: key, error: err}
}
