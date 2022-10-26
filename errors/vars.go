package errors

import "errors"

var (
	NetworkInvalid = errors.New("network is invalid")
	FatalInvalid   = errors.New("occur a Fatal error")
)
