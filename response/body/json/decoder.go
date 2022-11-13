package json

import (
	stdJSON "encoding/json"
	"io"
)

// Decoder implements respons/body.Decoder for decoding JSON.
type Decoder struct{}

// Decode attempts to decode data from a JSON encoded source to a target.
func (d Decoder) Decode(source io.Reader, target any) error {

	jsonDecoder := stdJSON.NewDecoder(source)
	return jsonDecoder.Decode(target)
}
