package runapi

import (
	"github.com/runapi-ai/core-sdk/go/core"
)

// Error is the base error type for all RunAPI SDK errors.
type Error = core.Error

// ErrorCode identifies the category of an API error.
type ErrorCode = core.ErrorCode

// Error code constants for matching errors returned by the SDK.
const (
	ErrAuthentication      = core.ErrAuthentication
	ErrInsufficientCredits = core.ErrInsufficientCredits
	ErrNotFound            = core.ErrNotFound
	ErrValidation          = core.ErrValidation
	ErrConflict            = core.ErrConflict
	ErrRateLimit           = core.ErrRateLimit
	ErrServiceUnavailable  = core.ErrServiceUnavailable
	ErrServer              = core.ErrServer
	ErrNetwork             = core.ErrNetwork
	ErrTimeout             = core.ErrTimeout
	ErrTaskTimeout         = core.ErrTaskTimeout
	ErrTaskFailed          = core.ErrTaskFailed
)
