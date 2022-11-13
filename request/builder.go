package request

import (
	"context"
	"io"
	"strings"

	stdHTTP "net/http"
	"net/url"
)

// Builder is for convenience and isolates the building of HTTP request types from
// doing the requests. This Builder can be used by potential future Client implementations.
type Builder struct{}

// BuildGetRequest builds an HTTP GET request using the same strategy as the standard library.
// It enforces that the provided context has a deadline set.
func (b Builder) BuildGetRequest(ctx context.Context, url string) (*stdHTTP.Request, error) {

	if !ContextHasDeadline(ctx) {
		return nil, ErrorNoContextDeadline
	}

	return stdHTTP.NewRequestWithContext(ctx, stdHTTP.MethodGet, url, nil)
}

// BuildHeadRequest builds an HTTP HEAD request using the same strategy as the standard library.
// It enforces that the provided context has a deadline set.
func (b Builder) BuildHeadRequest(ctx context.Context, url string) (*stdHTTP.Request, error) {

	if !ContextHasDeadline(ctx) {
		return nil, ErrorNoContextDeadline
	}

	return stdHTTP.NewRequestWithContext(ctx, stdHTTP.MethodHead, url, nil)
}

// BuildPostRequest builds an HTTP POST request using the same strategy as the standard library.
// It enforces that the provided context has a deadline set.
func (b Builder) BuildPostRequest(
	ctx context.Context,
	url string,
	contentType string,
	body io.Reader,
) (*stdHTTP.Request, error) {

	if !ContextHasDeadline(ctx) {
		return nil, ErrorNoContextDeadline
	}

	req, err := stdHTTP.NewRequestWithContext(ctx, stdHTTP.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	return req, nil
}

// BuildPostFormRequest builds an HTTP POST request using the same strategy as the standard library.
// It enforces that the provided context has a deadline set.
func (b Builder) BuildPostFormRequest(
	ctx context.Context,
	url string,
	data url.Values,
) (*stdHTTP.Request, error) {
	return b.BuildPostRequest(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
