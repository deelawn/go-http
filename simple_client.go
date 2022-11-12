package http

import (
	"context"
	"fmt"
	"io"
	stdHTTP "net/http"
	"net/url"
	"time"

	"github.com/deelawn/go-http/body"
	"github.com/deelawn/go-http/body/decoder"
	"github.com/deelawn/go-http/body/json"
	"github.com/deelawn/go-http/client"
	"github.com/deelawn/go-http/request"
)

type SimpleClient struct {
	Doer        request.Doer
	BodyDecoder body.Decoder
	Config      client.Config

	requestBuilder request.Builder
}

func NewSimpleClient(config client.Config) (*SimpleClient, error) {

	var bodyDecoder body.Decoder
	switch config.DecoderType {
	case decoder.TypeJSON:
		bodyDecoder = json.NewDecoder()
	}

	return &SimpleClient{
		Doer:        new(stdHTTP.Client),
		BodyDecoder: bodyDecoder,
		Config:      config,
	}, nil
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

	var (
		retries uint
		timer   = time.NewTimer(0)
	)

	for {

		if retries > c.Config.MaxRetries {
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

		resp, err = c.Doer.Do(req)
		if err == nil {
			break
		}

		retries++
		timer.Reset(c.Config.RetryFrequency)

	}

	defer resp.Body.Close()

	if c.BodyDecoder == nil {
		return
	}

	if err = c.BodyDecoder.Decode(resp.Body, respBodyTarget); err != nil {
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
