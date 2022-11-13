package client

import (
	"github.com/deelawn/go-http/request"
	"github.com/deelawn/go-http/request/retry/backoff"
	"github.com/deelawn/go-http/response/body"
)

// Config is used to configure an HTTP client via its constructor.
type Config struct {
	// RequestDoer does an HTTP request. It is set to an instance of stdlib/net/http/Client by default
	// in the constructor if no value is provided in the Config.
	RequestDoer request.Doer
	// ResponseBodyDecoder is used to decode the response body to the data type provided
	// as a parameter to the Do method. It is set to an instance of response/body/json/Decoder by
	// default in the constructor if no value is provided.
	ResponseBodyDecoder body.Decoder
	// RequestRetryBackoffStrategy is used to obtain request retry interval durations based on the number of
	// retries that have already been issued. It is set to an instance of request/retry/StrategyConstant
	// with an interval of 0 by default in the constructor if no value is provided.
	RequestRetryBackoffStrategy backoff.Strategy
	// MaxRequestRetries is the maximum number of times a request will be retried before returning an error.
	MaxRequestRetries uint
}
