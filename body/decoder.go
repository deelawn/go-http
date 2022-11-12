package body

import "io"

type Decoder interface {
	Decode(source io.Reader, target any) error
}
