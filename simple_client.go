package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdHTTP "net/http"
	"net/url"
	"time"

	"github.com/deelawn/go-http/client"
	"github.com/deelawn/go-http/request"
	"github.com/deelawn/go-http/request/retry/backoff"
	"github.com/deelawn/go-http/response/body"
	"github.com/deelawn/go-http/response/body/json"
)

var ErrNilClientFields = errors.New(
	"one or more of the required client fields is nil: RequestDoer, RequestRetryBackoffStrategy; " +
		"consider using the client constructor",
)

// SimpleClient is an HTTP client with response decoding and request retries baked in. It is
// highly encouraged that instances of SimpleClient be initialized using the NewSimpleClient
// constructor, as this will ensure that any required values not included in the Config
// are assigned to their default values.
type SimpleClient struct {
	RequestDoer                 request.Doer
	ResponseBodyDecoder         body.Decoder
	RequestRetryBackoffStrategy backoff.Strategy
	MaxRequestRetries           uint
	ResponseBodyReadIndicator   body.ReadIndicator

	requestBuilder request.Builder
}

// NewSimpleClient returns an instance of SimpleClient initialized from the provided Config or
// from defaults if required Config values are not provided.
func NewSimpleClient(config client.Config) *SimpleClient {

	simpleClient := SimpleClient{
		RequestDoer:                 config.RequestDoer,
		ResponseBodyDecoder:         config.ResponseBodyDecoder,
		RequestRetryBackoffStrategy: config.RequestRetryBackoffStrategy,
		MaxRequestRetries:           config.MaxRequestRetries,
		ResponseBodyReadIndicator:   config.ResponseBodyReadIndicator,
	}

	if simpleClient.RequestDoer == nil {
		simpleClient.RequestDoer = new(stdHTTP.Client)
	}

	if simpleClient.ResponseBodyDecoder == nil {
		simpleClient.ResponseBodyDecoder = new(json.Decoder)
	}

	if simpleClient.RequestRetryBackoffStrategy == nil {
		simpleClient.RequestRetryBackoffStrategy = backoff.NewStrategyConstant(0)
	}

	if simpleClient.ResponseBodyReadIndicator == nil {
		simpleClient.ResponseBodyReadIndicator = body.ReadIndicatorDefault
	}

	return &simpleClient
}

// Do does an HTTP request using the defined retry strategy in the event of failures. It handles reading the
// response body and uses the defined decoder to decode it to the provided respBodyTarget.
func (c *SimpleClient) Do(req *stdHTTP.Request, respBodyTarget any) (resp *stdHTTP.Response, err error) {

	// Avoid nil pointer errors. There is no context set so it will end up returning a
	// no context deadline error.
	if req == nil {
		req = new(stdHTTP.Request)
	}

	// It's going to fail sooner or later if no deadline is set, so might as well fail right away
	// to enforce good coding practices.
	ctx := req.Context()
	if !request.ContextHasDeadline(ctx) {
		return nil, request.ErrorNoContextDeadline
	}

	// There will be no problem if the constructor is used, but there could be otherwise. Do this
	// to avoid panics due to nil pointer references.
	if c.RequestDoer == nil || c.RequestRetryBackoffStrategy == nil {
		return nil, ErrNilClientFields
	}

	// If the context is done from the start, return an error. This avoids the case where the context
	// is done from the start and the request is done anyway because the timer channel read is selected.
	// This should make behavior more predictable.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var (
		retries uint
		// Don't wait at all before Doing the first request.
		timer = time.NewTimer(0)
	)

	for {

		if retries > c.MaxRequestRetries {
			return resp, err
		}

		select {
		case <-timer.C:
		case <-ctx.Done():
			// It's possible for the context to expire while waiting on the retry timer.
			// If there is no error value (likely impossible) then use the context error value.
			if err == nil {
				err = ctx.Err()
			}

			return resp, err
		}

		resp, err = c.RequestDoer.Do(req)
		if err == nil {
			break
		}

		// Error occurred so we need to retry using the defined strategy.
		retries++
		timer.Reset(c.RequestRetryBackoffStrategy.IntervalForRetry(retries))

	}

	// Do not parse nor read the response body nor close the reader if the custom
	// function provided evaluates to false based on the response contents.
	if !c.ResponseBodyReadIndicator(resp) {
		return
	}

	defer resp.Body.Close()

	// Don't decode if no decoder is provided.
	if c.ResponseBodyDecoder == nil {
		return
	}

	if err = c.ResponseBodyDecoder.Decode(resp.Body, respBodyTarget); err != nil {
		return resp, fmt.Errorf("error decoding response: %w", err)
	}

	return
}

// Get wraps Do by building the HTTP GET request. This is to achieve stdlib/net/http.Client feature parity.
func (c *SimpleClient) Get(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildGetRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}

// Head wraps Do by building the HTTP HEAD request. This is to achieve stdlib/net/http.Client feature parity.
func (c *SimpleClient) Head(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildHeadRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}

// Post wraps Do by building the HTTP POST request. This is to achieve stdlib/net/http.Client feature parity.
func (c *SimpleClient) Post(
	ctx context.Context,
	url string,
	contentType string,
	body io.Reader,
	respBodyTarget any,
) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildPostRequest(ctx, url, contentType, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}

// PostForm wraps Do by building the HTTP POST request. This is to achieve stdlib/net/http.Client feature parity.
func (c *SimpleClient) PostForm(
	ctx context.Context,
	url string,
	data url.Values,
	respBodyTarget any,
) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildPostFormRequest(ctx, url, data)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}
