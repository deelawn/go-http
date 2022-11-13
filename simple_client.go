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

type SimpleClient struct {
	RequestDoer                 request.Doer
	ResponseBodyDecoder         body.Decoder
	RequestRetryBackoffStrategy backoff.Strategy
	MaxRequestRetries           uint

	requestBuilder request.Builder
}

func NewSimpleClient(config client.Config) *SimpleClient {

	simpleClient := SimpleClient{
		RequestDoer:                 config.RequestDoer,
		ResponseBodyDecoder:         config.ResponseBodyDecoder,
		RequestRetryBackoffStrategy: config.RequestRetryBackoffStrategy,
		MaxRequestRetries:           config.MaxRequestRetries,
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

	return &simpleClient
}

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

	var (
		retries         uint
		timer           = time.NewTimer(0)
		backoffStrategy backoff.Strategy
	)

	backoffStrategy = c.RequestRetryBackoffStrategy
	if backoffStrategy == nil {
		backoffStrategy = backoff.NewStrategyConstant(0)
	}

	for {

		if retries > c.MaxRequestRetries {
			return resp, err
		}

		select {
		case <-timer.C:
		case <-ctx.Done():
			if err == nil {
				err = ctx.Err()
			}

			return resp, err
		}

		resp, err = c.RequestDoer.Do(req)
		if err == nil {
			break
		}

		retries++
		timer.Reset(backoffStrategy.IntervalForRetry(retries))

	}

	defer resp.Body.Close()

	if c.ResponseBodyDecoder == nil {
		return
	}

	if err = c.ResponseBodyDecoder.Decode(resp.Body, respBodyTarget); err != nil {
		return resp, fmt.Errorf("error decoding response: %w", err)
	}

	return
}

func (c *SimpleClient) Get(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildGetRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}

func (c *SimpleClient) Head(ctx context.Context, url string, respBodyTarget any) (*stdHTTP.Response, error) {

	req, err := c.requestBuilder.BuildHeadRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	return c.Do(req, respBodyTarget)
}

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
