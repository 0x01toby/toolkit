package errs

import "errors"

var (
	InValidNetwork = errors.New("invalid network")
	InvalidParams  = errors.New("invalid params")
	FatalInvalid   = errors.New("fatal error")
)
