package request

import (
	"context"
	"io"
	"strings"

	stdHTTP "net/http"
	"net/url"
)

type Builder struct{}

func (b Builder) BuildGetRequest(ctx context.Context, url string) (*stdHTTP.Request, error) {

	if !ContextHasDeadline(ctx) {
		return nil, ErrorNoContextDeadline
	}

	return stdHTTP.NewRequestWithContext(ctx, stdHTTP.MethodGet, url, nil)
}

func (b Builder) BuildHeadRequest(ctx context.Context, url string) (*stdHTTP.Request, error) {

	if !ContextHasDeadline(ctx) {
		return nil, ErrorNoContextDeadline
	}

	return stdHTTP.NewRequestWithContext(ctx, stdHTTP.MethodHead, url, nil)
}

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

func (b Builder) BuildPostFormRequest(
	ctx context.Context,
	url string,
	data url.Values,
) (*stdHTTP.Request, error) {
	return b.BuildPostRequest(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
