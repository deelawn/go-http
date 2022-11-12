package client

import (
	"time"

	"github.com/deelawn/go-http/body/decoder"
)

type Config struct {
	DecoderType    decoder.Type
	MaxRetries     uint
	RetryFrequency time.Duration
}
