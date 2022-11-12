package json

import (
	stdJSON "encoding/json"
	"io"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d Decoder) Decode(source io.Reader, target any) error {

	jsonDecoder := stdJSON.NewDecoder(source)
	return jsonDecoder.Decode(target)
}
