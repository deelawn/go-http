package client

import (
	"github.com/deelawn/go-http/request/retry/backoff"
	"github.com/deelawn/go-http/response/body/decoder"
)

type Config struct {
	DecoderType     decoder.Type
	MaxRetries      uint
	BackoffStrategy backoff.Strategy
}
