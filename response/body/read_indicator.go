package body

import stdHTTP "net/http"

// ReadIndicator is a function type that returns true if a response body should be
// read (decoded) and closed before returning the response.
type ReadIndicator func(*stdHTTP.Response) bool

// ReadIndicatorDefault will always return true.
func ReadIndicatorDefault(resp *stdHTTP.Response) bool {

	// The body can't be read from a nil response.
	if resp == nil {
		return false
	}

	return true
}

// ReadIndicatorErrorCode will return true if the HTTP status code is not
// in the error range.
func ReadIndicatorErrorCode(resp *stdHTTP.Response) bool {

	// The body can't be read from a nil response.
	if resp == nil {
		return false
	}

	return !(resp.StatusCode >= stdHTTP.StatusBadRequest)
}
