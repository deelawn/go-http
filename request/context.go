package request

import (
	"context"
	"errors"
)

var (
	ErrorNoContextDeadline = errors.New("HTTP request is missing context deadline")
)

func ContextHasDeadline(ctx context.Context) bool {

	if ctx == nil {
		return false
	}

	_, ok := ctx.Deadline()
	return ok
}
