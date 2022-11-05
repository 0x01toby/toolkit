package errs

import "errors"

var (
	InValidRpc               = errors.New("invalid rpc")
	InValidNetwork           = errors.New("invalid network")
	InvalidParams            = errors.New("invalid params")
	InvalidInitComponent     = errors.New("invalid init component")
	FatalInvalid             = errors.New("fatal error")
	MaxRetryPollingBatchCall = errors.New("max retry polling batch call")
)
