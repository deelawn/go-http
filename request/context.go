package request

import (
	"context"
	"errors"
)

// ErrorNoContextDeadline is returned if a context is provided that doesn't have a deadline set.
var ErrorNoContextDeadline = errors.New("HTTP request is missing context deadline")

// ContextHasDeadline returns true if the provided context has a deadline set.
func ContextHasDeadline(ctx context.Context) bool {

	if ctx == nil {
		return false
	}

	_, ok := ctx.Deadline()
	return ok
}
