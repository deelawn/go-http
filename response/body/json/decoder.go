package json

import (
	stdJSON "encoding/json"
	"errors"
	"io"
)

var ErrNilDecodeSource = errors.New("can't decode from nil source")

// Decoder implements respons/body.Decoder for decoding JSON.
type Decoder struct{}

// Decode attempts to decode data from a JSON encoded source to a target.
func (d Decoder) Decode(source io.Reader, target any) error {

	if source == nil {
		return ErrNilDecodeSource
	}

	jsonDecoder := stdJSON.NewDecoder(source)
	return jsonDecoder.Decode(target)
}
