package constants

import (
	"errors"
)

// Error Codes
var (
	ErrBindingRequest   = errors.New("error binding json")
	ErrForbidden        = errors.New("forbidden")
	ErrPanicRecovered   = errors.New("panic recovered")
	ErrQueuePublish     = errors.New("error publishing to queue")
	ErrProcessAnalytics = errors.New("error processing analytivs")
	ErrInternalServer   = errors.New("Internal server error")
)
