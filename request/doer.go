package request

import stdHTTP "net/http"

// Doer does an HTTP request. The Do method has the same signature as the
// stdlib/net/http.Client's Do method.
type Doer interface {
	Do(req *stdHTTP.Request) (*stdHTTP.Response, error)
}
