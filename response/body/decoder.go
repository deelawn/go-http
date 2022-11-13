package body

import "io"

// Decoder decodes data from source to target.
type Decoder interface {
	Decode(source io.Reader, target any) error
}
