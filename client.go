package http

import (
	"context"
	"io"
	stdHTTP "net/http"
	"net/url"
)

// Client provides the same base functionality as stdlib/net/http.Client.
// It additionally acceps Contexts and targets to which to decode the HTTP
// response body. This interface definition is for the convenience of
// code that imports this package.
type Client interface {
	Do(req *stdHTTP.Request, respBodyTarget any) (*stdHTTP.Response, error)
	Get(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error)
	Head(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error)
	Post(
		ctx context.Context,
		url string,
		contentType string,
		body io.Reader,
		respBodyTarget any,
	) (*stdHTTP.Response, error)
	PostForm(ctx context.Context, url string, data url.Values, respBodyTarget any) (*stdHTTP.Response, error)
}
