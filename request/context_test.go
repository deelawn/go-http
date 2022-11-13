package request_test

import (
	"context"
	"testing"
	"time"

	"github.com/deelawn/go-http/request"
)

func ctxWithCancel() context.Context {
	ctx, _ := context.WithCancel(context.Background())
	return ctx
}

func ctxWithDeadline() context.Context {
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	return ctx
}

func ctxWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	return ctx
}

func TestContextHasDeadline(t *testing.T) {

	tests := []struct {
		name      string
		ctx       context.Context
		expResult bool
	}{
		{
			name: "nil",
		},
		{
			name: "background",
			ctx:  context.Background(),
		},
		{
			name: "with cancel",
			ctx:  ctxWithCancel(),
		},
		{
			name: "with value",
			ctx:  context.WithValue(context.Background(), "a", "b"),
		},
		{
			name:      "with deadline",
			ctx:       ctxWithDeadline(),
			expResult: true,
		},
		{
			name:      "with timeout",
			ctx:       ctxWithTimeout(),
			expResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := request.ContextHasDeadline(tt.ctx)

			if result != tt.expResult {
				t.Errorf("expected %t, got %t", tt.expResult, result)
			}
		})
	}
}
