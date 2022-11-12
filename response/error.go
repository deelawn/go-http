package response

import (
	"fmt"
)

type Error struct {
	Status     string
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("HTTP error %s: %s", e.Status, e.Message)
}
