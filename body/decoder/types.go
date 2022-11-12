package decoder

import "errors"

var TypeUnknownError = errors.New("unknown decoder type")

type Type int

const (
	TypeUnknown Type = iota
	TypeJSON
)
