package body_test

import (
	stdHTTP "net/http"
	"testing"

	"github.com/deelawn/go-http/response/body"
)

func TestReadIndicator(t *testing.T) {

	tests := []struct {
		name          string
		readIndicator body.ReadIndicator
		resp          *stdHTTP.Response
		expResult     bool
	}{
		{
			name:          "default nil resp",
			readIndicator: body.ReadIndicatorDefault,
		},
		{
			name:          "default ok",
			readIndicator: body.ReadIndicatorDefault,
			resp:          &stdHTTP.Response{},
			expResult:     true,
		},
		{
			name:          "error code nil",
			readIndicator: body.ReadIndicatorErrorCode,
		},
		{
			name:          "error code with 400",
			readIndicator: body.ReadIndicatorErrorCode,
			resp:          &stdHTTP.Response{StatusCode: stdHTTP.StatusBadRequest},
		},
		{
			name:          "error code with 500",
			readIndicator: body.ReadIndicatorErrorCode,
			resp:          &stdHTTP.Response{StatusCode: stdHTTP.StatusInternalServerError},
		},
		{
			name:          "error code ok",
			readIndicator: body.ReadIndicatorErrorCode,
			resp:          &stdHTTP.Response{},
			expResult:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ok := tt.readIndicator(tt.resp)
			if ok != tt.expResult {
				t.Errorf("expected %t, got %t", tt.expResult, ok)
			}
		})
	}
}
